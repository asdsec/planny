package db

import (
	"com.github/asdsec/planny/internal/model"
	"context"
	"errors"
	"gorm.io/gorm"
)

type UserEntity struct {
	gorm.Model
	Username       string `gorm:"unique"`
	Email          string
	FirstName      string
	LastName       string
	HashedPassword string
}

type CreateUserArg struct {
	Username       string
	Email          string
	FirstName      string
	LastName       string
	HashedPassword string
}

func (store *SQLStore) CreateUser(ctx context.Context, arg CreateUserArg) (model.User, error) {
	userEntity := UserEntity{
		Username:       arg.Username,
		Email:          arg.Email,
		FirstName:      arg.FirstName,
		LastName:       arg.LastName,
		HashedPassword: arg.HashedPassword,
	}
	err := store.db.Create(&userEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return userEntity.toEmpty(), errors.New(ErrForeignKeyViolated)
		}
		return userEntity.toEmpty(), errors.New(ErrUnhandled)
	}
	return userEntity.toUser(), nil
}

func (store *SQLStore) GetUserById(ctx context.Context, id uint) (model.User, error) {
	userEntity := UserEntity{}
	err := store.db.First(&userEntity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userEntity.toEmpty(), errors.New(ErrRecordNotFound)
		}
		return userEntity.toEmpty(), errors.New(ErrUnhandled)
	}
	return userEntity.toUser(), nil
}

func (store *SQLStore) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	userEntity := UserEntity{}
	err := store.db.First(&userEntity, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userEntity.toEmpty(), errors.New(ErrRecordNotFound)
		}
		return userEntity.toEmpty(), errors.New(ErrUnhandled)
	}
	return userEntity.toUser(), nil
}

func (store *SQLStore) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	userEntity := UserEntity{}
	err := store.db.First(&userEntity, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userEntity.toEmpty(), errors.New(ErrRecordNotFound)
		}
		return userEntity.toEmpty(), errors.New(ErrUnhandled)
	}
	return userEntity.toUser(), nil
}

func (u *UserEntity) toUser() model.User {
	return model.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Password:  u.HashedPassword,
		UpdatedAt: u.UpdatedAt,
		CreatedAt: u.CreatedAt,
	}
}

func (u *UserEntity) toEmpty() model.User {
	return model.User{}
}
