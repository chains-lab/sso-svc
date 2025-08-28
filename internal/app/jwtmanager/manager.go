package jwtmanager

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/constant"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Manager struct {
	accessSK  string
	refreshSK string

	accessTTL  time.Duration
	refreshTTL time.Duration

	iss string
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func NewManager(cfg config.Config) Manager {
	return Manager{
		accessSK:  cfg.JWT.User.AccessToken.SecretKey,
		refreshSK: cfg.JWT.User.RefreshToken.SecretKey,

		accessTTL:  cfg.JWT.User.AccessToken.TokenLifetime,
		refreshTTL: cfg.JWT.User.RefreshToken.TokenLifetime,

		iss: cfg.Server.Name,
	}
}

func encryptAESGCM(plain string, key []byte) (string, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	// Пишем nonce в начало шифртекста (как ты и делал)
	ct := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return hex.EncodeToString(ct), nil
}

func decryptAESGCM(encHex string, key []byte) (string, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return "", err
	}
	ct, err := hex.DecodeString(encHex)
	if err != nil {
		return "", err
	}
	if len(ct) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := ct[:gcm.NonceSize()], ct[gcm.NonceSize():]
	pt, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

func (m Manager) EncryptAccess(token string) (string, error) {
	return encryptAESGCM(token, []byte(m.accessSK))
}

func (m Manager) EncryptRefresh(token string) (string, error) {
	return encryptAESGCM(token, []byte(m.refreshSK))
}

func (m Manager) DecryptRefresh(encryptedToken string) (string, error) {
	raw, err := decryptAESGCM(encryptedToken, []byte(m.refreshSK))
	if err != nil {
		return "", fmt.Errorf("decrypt refresh: %w", err)
	}

	return raw, nil
}

func (m Manager) ParseRefreshClaims(token string) (auth.UsersClaims, error) {
	return auth.VerifyUserJWT(context.Background(), token, m.refreshSK)
}

func (m Manager) GenerateAccess(
	userID uuid.UUID,
	sessionID uuid.UUID,
	idn string,
) (string, error) {
	return auth.GenerateUserJWT(auth.GenerateUserJwtRequest{
		Issuer:  m.iss,
		User:    userID,
		Session: sessionID,
		Role:    idn,
		Ttl:     m.accessTTL,
	}, m.accessSK)
}

func (m Manager) GenerateRefresh(
	userID uuid.UUID,
	sessionID uuid.UUID,
	role string,
	emailVerified bool,
) (string, error) {
	return auth.GenerateUserJWT(auth.GenerateUserJwtRequest{
		Issuer:   m.iss,
		Audience: []string{constant.ServiceApiGateway},
		User:     userID,
		Session:  sessionID,
		Verified: emailVerified,
		Role:     role,
		Ttl:      m.refreshTTL,
	}, m.refreshSK)
}
