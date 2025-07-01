package repository

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionID int       `json:"session_id"`
	UserGUID  uuid.UUID `json:"user_guid"`
	UserAgent string    `json:"user_agent"`
	TokenHash string    `json:"token_hash"`
	IpAddress string    `json:"ip_address"`
	ExpiresAt time.Time `json:"expires_at"`
}
