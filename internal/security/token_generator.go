package security

import "time"

type TokenGenerator interface {
	// Generate creates a new token for a specific username and duration
	Generate(userID uint, username string, duration time.Duration) (string, *TokenPayload, error)

	// Verify checks if the token is valid or not
	Verify(token string) (*TokenPayload, error)
}
