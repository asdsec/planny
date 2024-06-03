package db

import (
	"com.github/asdsec/planny/internal/model"
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SessionEntity struct {
	ID           uuid.UUID `gorm:"primarykey"`
	Username     string
	RefreshToken string
	UserAgent    string
	ClientIp     string
	ExpiresAt    time.Time
	UserID       uint
	User         UserEntity `gorm:"foreignKey:UserID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateSessionArg struct {
	ID           uuid.UUID
	UserID       uint
	Username     string
	RefreshToken string
	UserAgent    string
	ClientIp     string
	ExpiresAt    time.Time
}

func (store *SQLStore) CreateSession(ctx context.Context, arg CreateSessionArg) (model.Session, error) {
	sessionEntity := SessionEntity{
		ID:           arg.ID,
		UserID:       arg.UserID,
		Username:     arg.Username,
		RefreshToken: arg.RefreshToken,
		UserAgent:    arg.UserAgent,
		ClientIp:     arg.ClientIp,
		ExpiresAt:    arg.ExpiresAt,
	}
	err := store.db.Create(&sessionEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return sessionEntity.toEmpty(), errors.New(ErrForeignKeyViolated)
		}
		return sessionEntity.toEmpty(), errors.New(ErrUnhandled)
	}
	return sessionEntity.toSession(), nil
}

func (store *SQLStore) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (model.Session, error) {
	var sessionEntity SessionEntity
	err := store.db.Where("refresh_token = ?", refreshToken).First(&sessionEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sessionEntity.toEmpty(), errors.New(ErrRecordNotFound)
		}
		return sessionEntity.toEmpty(), errors.New(ErrUnhandled)
	}
	return sessionEntity.toSession(), nil
}

func (s *SessionEntity) toSession() model.Session {
	return model.Session{
		ID:           s.ID,
		UserID:       s.UserID,
		Username:     s.Username,
		RefreshToken: s.RefreshToken,
		UserAgent:    s.UserAgent,
		ClientIp:     s.ClientIp,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

func (s *SessionEntity) toEmpty() model.Session {
	return model.Session{}
}
