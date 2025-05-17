package domain

import (
	"context"
)

const (
	CollectionTask = "tasks"
)

type Task struct {
	ID     uint   `json:"id"`
	Title  string `json:"title" form:"title" binding:"required"`
	UserID uint   `json:"user_id"`
}

type TaskRepository interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}

type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}
