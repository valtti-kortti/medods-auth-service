package webhook

import (
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	Event     string    `json:"event"`
	UserGUID  uuid.UUID `json:"user_guid"`
	OldIP     string    `json:"old_ip"`
	NewIP     string    `json:"new_ip"`
	Timestamp time.Time `json:"timestamp"`
}
