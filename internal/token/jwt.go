package token

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenManager interface {
	GenerateToken(sessionID int, userIP string) (string, error)
	ValidToken(accessToken string) (*Claims, error)
}

type jwtToken struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTTokenManager(secret []byte) JWTTokenManager {
	return &jwtToken{secret: secret, ttl: 15 * time.Minute}
}

func (t *jwtToken) GenerateToken(sessionID int, userIP string) (string, error) {
	sessionIDString := strconv.Itoa(sessionID)

	claims := Claims{
		SessionID: sessionID,
		UserIP:    userIP,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   sessionIDString,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessToken, err := token.SignedString(t.secret)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (t *jwtToken) ValidToken(accessToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.secret, nil
	})

	if err != nil {
		log.Println("JWT parse error:", err)
		return nil, err
	}

	if !token.Valid {
		log.Println("Token is not valid (but parse succeeded)")
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		log.Println("claims", err)
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}
