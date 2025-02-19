package jwtmanager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/recovery-flow/roles"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/tokens"
)

type JWTManager interface {
	EncryptAccess(token string) (string, error)
	EncryptRefresh(token string) (string, error)
	DecryptRefresh(encryptedToken string) (string, error)

	GenerateAccess(
		sub uuid.UUID,
		sessionID uuid.UUID,
		role roles.UserRole,
	) (string, error)

	GenerateRefresh(
		sub uuid.UUID,
		sessionID uuid.UUID,
		role roles.UserRole,
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
	sub uuid.UUID,
	sessionID uuid.UUID,
	role roles.UserRole,
) (string, error) {
	roleU := string(role)
	ses := sessionID.String()
	return generateJWT(m.iss, sub.String(), nil, &roleU, &ses, m.accessTTL, m.accessSK)
}

func (m *jwtmanager) GenerateRefresh(
	sub uuid.UUID,
	sessionID uuid.UUID,
	role roles.UserRole,
) (string, error) {
	roleU := string(role)
	ses := sessionID.String()
	return generateJWT(m.iss, sub.String(), nil, &roleU, &ses, m.refreshTTL, m.refreshSK)
}

func generateJWT(
	iss string,
	sub string,
	aud []string,
	role *string,
	sessionID *string,
	ttl time.Duration,
	sk string,
) (string, error) {
	expirationTime := time.Now().Add(ttl * time.Second)
	claims := &tokens.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    iss,
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		Role:      role,
		SessionID: sessionID,
	}
	if role != nil {
		_, err := roles.ParseUserRole(*role)
		if err != nil {
			return "", fmt.Errorf("invalid role: %w", err)
		}
		claims.Role = role
	}
	if sessionID != nil {
		_, err := uuid.Parse(*sessionID)
		if err != nil {
			return "", fmt.Errorf("invalid device id: %w", err)
		}
		claims.SessionID = sessionID
	}
	if aud != nil {
		claims.Audience = aud
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sk))
}
