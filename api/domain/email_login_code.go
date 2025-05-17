package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type EmailLoginCode struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"index;not null"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type EmailCodeRepository interface {
	Create(ctx context.Context, code EmailLoginCode) error
	Verify(ctx context.Context, email string, code string) (EmailLoginCode, error)
	Delete(ctx context.Context, code EmailLoginCode) error
}

type EmailCodeService interface {
	RequestLoginCode(ctx context.Context, email string) error
	VerifyLoginCode(ctx context.Context, email, code string) (bool, error)
}
