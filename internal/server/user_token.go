package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"short_url/internal/repository"
)

const password = "x35k9f"

func generateEncodedToken() (string, string) {
	userID := repository.GenUUID()
	userToken := encode([]byte(userID))

	return userID, userToken
}

func decodeToken(token string) (string, error) {
	userID, err := decode(token)
	return userID, err
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

func checkUserExist(r *http.Request, repo repository.URLRepo) (string, bool) {
	cookies := r.Cookies()
	var userExist bool
	var userId string
	var err error

	for _, v := range cookies {
		if v.Name == "Auth" {
			userId, err = decodeToken(v.Value)
			if err != nil {
				continue
			}
			userExist = repo.IsUserExist(userId)
			break
		}
	}
	return userId, userExist
}

func setNewUserToken(w http.ResponseWriter) string {
	userId, newUserToken := generateEncodedToken()
	newCookie := &http.Cookie{
		Name:  "Auth",
		Value: newUserToken,
	}
	http.SetCookie(w, newCookie)
	return userId
}
