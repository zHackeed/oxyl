package models

import "time"

type AuthenticationRequest struct {
	Identifier string `json:"agent_id"`
}

type AuthenticationResponse struct {
	AccessToken struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	} `json:"access_token"`
	RefreshToken struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	} `json:"refresh_token"`
}

type AuthenticationShutdownRequest struct {
	AgentId string `json:"agent_id"`
}
