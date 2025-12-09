package dto

import "time"

type CreateAPIKeyRequest struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Expiry      string   `json:"expiry"`
}

type CreateAPIKeyResponse struct {
	APIKey    string    `json:"api_key"`
	ExpiresAt time.Time `json:"expires_at"`
}
