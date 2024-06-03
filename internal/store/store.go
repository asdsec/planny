package db

import (
	"com.github/asdsec/planny/internal/model"
	"context"
	"gorm.io/gorm"
)

const (
	ErrRecordNotFound     = "record not found"
	ErrForeignKeyViolated = "foreign key constraint violated"
	ErrDuplicatedKey      = "duplicated key"
	ErrUnhandled          = "unhandled error"
)

type Store interface {
	CreateUser(ctx context.Context, arg CreateUserArg) (model.User, error)
	GetUserById(ctx context.Context, id uint) (model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	CreateSession(ctx context.Context, arg CreateSessionArg) (model.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (model.Session, error)
	CreatePlan(ctx context.Context, arg CreatePlanArg) (model.Plan, error)
	GetPlanByID(ctx context.Context, id uint) (model.Plan, error)
	ListPlansByUserID(ctx context.Context, userID uint) ([]model.Plan, error)
	UpdatePlanByID(ctx context.Context, arg UpdatePlanArg) (model.Plan, error)
	DeletePlanByID(ctx context.Context, id uint) error
}

// SQLStore represents the store
type SQLStore struct {
	db *gorm.DB
}

// NewStore creates a new store
func NewStore(db *gorm.DB) Store {
	return &SQLStore{db: db}
}
