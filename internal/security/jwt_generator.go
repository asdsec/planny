package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

// JWTGenerator is a JSON Web Token generator
type JWTGenerator struct {
	secretKey string
}

// NewJWTGenerator creates a new JWTGenerator
func NewJWTGenerator(secretKey string) (TokenGenerator, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTGenerator{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (g *JWTGenerator) Generate(userID uint, username string, duration time.Duration) (string, *TokenPayload, error) {
	payload, err := NewTokenPayload(userID, username, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(g.secretKey))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (g *JWTGenerator) Verify(token string) (*TokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(g.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &TokenPayload{}, keyFunc)
	if err != nil {
		var verr *jwt.ValidationError
		ok := errors.As(err, &verr)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*TokenPayload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
