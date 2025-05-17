package postgres

import (
	"time"

	"gorm.io/gorm"
	"yobank/domain"
)

type Task struct {
	ID        uint        `gorm:"primaryKey"`
	Title     string      `gorm:"not null"`
	UserID    uint        `gorm:"not null"`
	User      domain.User `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// ToDomain converts Task model to domain Task
func (t *Task) ToDomain() *domain.Task {
	return &domain.Task{
		ID:     t.ID,
		Title:  t.Title,
		UserID: t.UserID,
	}
}

// FromDomain converts domain Task to Task model
func (t *Task) FromDomain(task *domain.Task) {
	t.Title = task.Title
	t.UserID = task.UserID
}
