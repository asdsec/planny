package db

import (
	"com.github/asdsec/planny/internal/model"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Status string

const (
	Done       Status = "done"
	InProgress Status = "in_progress"
	Cancelled  Status = "cancelled"
)

type PlanEntity struct {
	ID          uint `gorm:"primarykey"`
	Title       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Status      Status
	UserID      uint
	User        UserEntity `gorm:"foreignKey:UserID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreatePlanArg struct {
	Title       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Status      Status
	UserID      uint
}

type UpdatePlanArg struct {
	ID          uint
	Title       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Status      Status
}

func (arg *UpdatePlanArg) toEntity() *PlanEntity {
	return &PlanEntity{
		ID:          arg.ID,
		Title:       arg.Title,
		Description: arg.Description,
		StartDate:   arg.StartDate,
		EndDate:     arg.EndDate,
		Status:      arg.Status,
	}
}

func (store *SQLStore) CreatePlan(ctx context.Context, arg CreatePlanArg) (model.Plan, error) {
	planEntity := PlanEntity{
		Title:       arg.Title,
		Description: arg.Description,
		StartDate:   arg.StartDate,
		EndDate:     arg.EndDate,
		Status:      arg.Status,
		UserID:      arg.UserID,
	}
	err := store.db.Create(&planEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return planEntity.toEmpty(), errors.New(ErrForeignKeyViolated)
		}
		return planEntity.toEmpty(), errors.New(ErrUnhandled)
	}
	return planEntity.toPlan(), nil
}

func (store *SQLStore) GetPlanByID(ctx context.Context, id uint) (model.Plan, error) {
	var planEntity PlanEntity
	err := store.db.First(&planEntity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return planEntity.toEmpty(), errors.New(ErrRecordNotFound)
		}
		return planEntity.toEmpty(), errors.New(ErrUnhandled)
	}
	return planEntity.toPlan(), nil
}

func (store *SQLStore) ListPlansByUserID(ctx context.Context, userID uint) ([]model.Plan, error) {
	var planEntities []PlanEntity
	err := store.db.Where("user_id = ?", userID).Find(&planEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(ErrRecordNotFound)
		}
		return nil, errors.New(ErrUnhandled)
	}
	var plans []model.Plan
	for _, planEntity := range planEntities {
		plans = append(plans, planEntity.toPlan())
	}
	return plans, nil
}

func (store *SQLStore) UpdatePlanByID(ctx context.Context, arg UpdatePlanArg) (model.Plan, error) {
	planEntity := arg.toEntity()
	err := store.db.Model(planEntity).Updates(planEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return planEntity.toPlan(), errors.New(ErrRecordNotFound)
		}
		return planEntity.toPlan(), errors.New(ErrUnhandled)
	}
	err = store.db.Model(planEntity).Find(planEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return planEntity.toPlan(), errors.New(ErrRecordNotFound)
		}
		return planEntity.toPlan(), errors.New(ErrUnhandled)
	}
	return planEntity.toPlan(), nil
}

func (store *SQLStore) DeletePlanByID(ctx context.Context, id uint) error {
	err := store.db.Delete(&PlanEntity{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(ErrRecordNotFound)
		}
		return errors.New(ErrUnhandled)
	}
	return nil
}

func (e *PlanEntity) toEmpty() model.Plan {
	return model.Plan{}
}

func (e *PlanEntity) toPlan() model.Plan {
	return model.Plan{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		StartDate:   e.StartDate,
		EndDate:     e.EndDate,
		Status:      e.Status.toModelStatus(),
		UserID:      e.UserID,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func (s *Status) toModelStatus() model.Status {
	switch *s {
	case Done:
		return model.Done
	case Cancelled:
		return model.Cancelled
	default:
		return model.InProgress
	}
}
