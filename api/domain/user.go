package domain

import (
	"context"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                uint    `gorm:"primaryKey"`
	Email             *string `gorm:"uniqueIndex"`
	Username          string  `gorm:"uniqueIndex;not null"` // для логики приложения
	TelegramID        *int64  `gorm:"uniqueIndex"`
	TelegramUsername  *string // raw Telegram username
	TelegramFirstName *string
	Wallets           []Wallet `gorm:"foreignKey:UserID"`

	Friends []Friend `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserRepository interface {
	CreateWithTx(tx *gorm.DB, user *User) error
	Create(c context.Context, user *User) error
	Fetch(c context.Context) ([]User, error)
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id string) (User, error)
	GetByTelegramID(c context.Context, tgID int64) (User, error)
	GetByTelegramIDWithTx(tx *gorm.DB, tgID int64) (*User, error)
}

type UserService interface {
	CreateUserWithWallet(ctx context.Context, tgUser initdata.User) (*User, error)
}
