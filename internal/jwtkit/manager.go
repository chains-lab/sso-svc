package jwtkit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/tokens"
	"github.com/hs-zavet/tokens/roles"
	"github.com/pkg/errors"
)

type Manager struct {
	accessSK  string
	refreshSK string

	accessTTL  time.Duration
	refreshTTL time.Duration

	iss string
}

func NewManager(cfg config.Config) Manager {
	return Manager{
		accessSK:  cfg.JWT.AccessToken.SecretKey,
		refreshSK: cfg.JWT.RefreshToken.SecretKey,

		accessTTL:  cfg.JWT.AccessToken.TokenLifetime,
		refreshTTL: cfg.JWT.RefreshToken.TokenLifetime,

		iss: cfg.Server.Name,
	}
}

func (m Manager) EncryptAccess(token string) (string, error) {
	block, err := aes.NewCipher([]byte(m.accessSK))
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12) // 96-битный nonce для AES-GCM
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(token), nil)
	return hex.EncodeToString(ciphertext), nil
}

func (m Manager) EncryptRefresh(token string) (string, error) {
	block, err := aes.NewCipher([]byte(m.refreshSK))
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12) // 96-битный nonce для AES-GCM
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(token), nil)
	return hex.EncodeToString(ciphertext), nil
}

func (m Manager) DecryptRefresh(encryptedToken string) (string, error) {
	ciphertext, err := hex.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(m.refreshSK))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aesGCM.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:aesGCM.NonceSize()], ciphertext[aesGCM.NonceSize():]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (m Manager) GenerateAccess(
	userID uuid.UUID,
	sessionID uuid.UUID,
	subTypeID uuid.UUID,
	idn roles.Role,
) (string, error) {
	return tokens.GenerateUserJWT(tokens.GenerateUserJwtRequest{
		Issuer:       m.iss,
		Account:      userID,
		Session:      sessionID,
		Subscription: subTypeID,
		Role:         idn,
		Ttl:          m.accessTTL,
	}, m.accessSK)
}

func (m Manager) GenerateRefresh(
	userID uuid.UUID,
	sessionID uuid.UUID,
	subTypeID uuid.UUID,
	idn roles.Role,
) (string, error) {
	return tokens.GenerateUserJWT(tokens.GenerateUserJwtRequest{
		Issuer:       m.iss,
		Account:      userID,
		Session:      sessionID,
		Subscription: subTypeID,
		Role:         idn,
		Ttl:          m.refreshTTL,
	}, m.refreshSK)
}
