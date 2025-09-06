package entity

import "time"

// Session represents a user session with refresh token and expiry
type Session struct {
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expires_at"`
}
