package webhook

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type webhook struct {
	url string
}

type Webhook interface {
	SendWebhook(newIP, oldIP string, guid uuid.UUID)
}

func NewWebhook(url string) Webhook {
	return &webhook{url: url}
}

func (w *webhook) SendWebhook(newIP, oldIP string, guid uuid.UUID) {
	data := Payload{
		Event:     "Changed IP",
		UserGUID:  guid,
		OldIP:     oldIP,
		NewIP:     newIP,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Webhook marshal error: %v", err)
		return
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post(
		w.url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Printf("Webhook send failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("Webhook bad status: %d", resp.StatusCode)
	}
}
