package token

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	SessionID int    `json:"session_id"`
	UserIP    string `json:"user_ip"`
	jwt.RegisteredClaims
}
