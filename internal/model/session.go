package model

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID           uuid.UUID
	UserID       uint
	Username     string
	RefreshToken string
	UserAgent    string
	ClientIp     string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
