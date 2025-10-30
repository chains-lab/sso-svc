package jwtmanager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

type Service struct {
	accessSK  string
	refreshSK string

	accessTTL  time.Duration
	refreshTTL time.Duration

	iss string
}

type Config struct {
	AccessSK  string
	RefreshSK string

	AccessTTL  time.Duration
	RefreshTTL time.Duration

	Iss string
}

func NewManager(cfg Config) Service {
	return Service{
		accessSK:  cfg.AccessSK,
		refreshSK: cfg.RefreshSK,

		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,

		iss: cfg.Iss,
	}
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
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
