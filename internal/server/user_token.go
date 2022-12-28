package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"short_url/internal/repository"
)

const password = "x35k9f"

func generateEncodedToken() (string, string) {
	userId := repository.GenUUID()
	userToken := encode([]byte(userId))

	return userId, userToken
}

func decodeToken(token string) (string, error) {
	userId, err := decode(token)
	return userId, err
}

func encode(msg []byte) string {
	key := sha256.Sum256([]byte(password))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		panic(err)
	}
	nonce := key[len(key)-aesgcm.NonceSize():]

	var token []byte

	token = aesgcm.Seal(nil, nonce, msg, nil)
	return hex.EncodeToString(token)
}

func decode(msg string) (string, error) {
	key := sha256.Sum256([]byte(password))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	// создаём вектор инициализации
	nonce := key[len(key)-aesgcm.NonceSize():]

	decoded, err := hex.DecodeString(msg)
	if err != nil {
		return "", err
	}

	decrypted, err := aesgcm.Open(nil, nonce, decoded, nil)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}
