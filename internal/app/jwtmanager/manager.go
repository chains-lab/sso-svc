package jwtmanager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"time"

	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-svc/internal/config"
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

func NewManager(cfg config.Config) Manager {
	return Manager{
		accessSK:  cfg.JWT.User.AccessToken.SecretKey,
		refreshSK: cfg.JWT.User.RefreshToken.SecretKey,

		accessTTL:  cfg.JWT.User.AccessToken.TokenLifetime,
		refreshTTL: cfg.JWT.User.RefreshToken.TokenLifetime,

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
	idn roles.Role,
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
	idn roles.Role,
) (string, error) {
	return auth.GenerateUserJWT(auth.GenerateUserJwtRequest{
		Issuer:  m.iss,
		User:    userID,
		Session: sessionID,
		Role:    idn,
		Ttl:     m.refreshTTL,
	}, m.refreshSK)
}
