package security

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Different types of error returned by the Verify function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// TokenPayload contains the payload data of the token
type TokenPayload struct {
	ID        uuid.UUID `json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewTokenPayload creates a new token payload with a specific username and duration
func NewTokenPayload(userID uint, username string, duration time.Duration) (*TokenPayload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &TokenPayload{
		ID:        tokenID,
		UserID:    userID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *TokenPayload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
