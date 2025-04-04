package jwtmanager

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
	"github.com/hs-zavet/tokens/identity"
	"github.com/pkg/errors"
)

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	GenerateAccess(
		userID *uuid.UUID,
		sessionID *uuid.UUID,
		subTypeID *uuid.UUID,
		idn identity.IdnType,
	) (string, error)

	GenerateRefresh(
		userID *uuid.UUID,
		sessionID *uuid.UUID,
		subTypeID *uuid.UUID,
		idn identity.IdnType,
	) (string, error)
}

type jwtmanager struct {
	accessSK  string
	refreshSK string

	accessTTL  time.Duration
	refreshTTL time.Duration

	iss string
}

func NewJWTManager(cfg *config.Config) JWTManager {
	return &jwtmanager{
		accessSK:  cfg.JWT.AccessToken.SecretKey,
		refreshSK: cfg.JWT.RefreshToken.SecretKey,

		accessTTL:  cfg.JWT.AccessToken.TokenLifetime,
		refreshTTL: cfg.JWT.RefreshToken.TokenLifetime,

		iss: cfg.Server.Name,
	}
}

func (m *jwtmanager) EncryptAccess(token string) (string, error) {
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

func (m *jwtmanager) EncryptRefresh(token string) (string, error) {
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

func (m *jwtmanager) DecryptRefresh(encryptedToken string) (string, error) {
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

func (m *jwtmanager) GenerateAccess(
	userID *uuid.UUID,
	sessionID *uuid.UUID,
	subTypeID *uuid.UUID,
	idn identity.IdnType,
) (string, error) {
	return tokens.GenerateJWT(m.iss, userID.String(), nil, m.accessTTL, idn, sessionID, userID, subTypeID, m.accessSK)
}

func (m *jwtmanager) GenerateRefresh(
	userID *uuid.UUID,
	sessionID *uuid.UUID,
	subTypeID *uuid.UUID,
	idn identity.IdnType,
) (string, error) {
	return tokens.GenerateJWT(m.iss, userID.String(), nil, m.refreshTTL, idn, sessionID, userID, subTypeID, m.refreshSK)
}
