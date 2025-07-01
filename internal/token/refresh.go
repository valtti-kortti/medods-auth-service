package token

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

type RefreshManager interface {
	GenerateToken() (string, string, error)
	CompareHash(token, hash string) error
}
type refreshToken struct {
	lenToken int
}

func NewRefreshTokenManager(lenToken int) RefreshManager {
	return &refreshToken{lenToken: lenToken}
}

func (rm *refreshToken) GenerateToken() (string, string, error) {
	g := make([]byte, rm.lenToken)

	_, err := rand.Read(g)
	if err != nil {
		return "", "", err
	}

	refreshToken := base64.URLEncoding.EncodeToString(g)

	hashRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), 10)
	if err != nil {
		return "", "", err
	}
	return refreshToken, string(hashRefreshToken), nil

}

func (rm *refreshToken) CompareHash(token, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	if err != nil {
		return err
	}

	return nil
}
