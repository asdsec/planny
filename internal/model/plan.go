package model

import (
	"time"
)

type Status string

const (
	Done       Status = "done"
	InProgress Status = "in_progress"
	Cancelled  Status = "cancelled"
)

type Plan struct {
	ID          uint
	Title       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Status      Status
	UserID      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
