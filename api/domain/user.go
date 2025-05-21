package domain

import (
	"context"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                uint     `gorm:"primaryKey" json:"id"`
	Email             *string  `gorm:"uniqueIndex" json:"email"`
	Username          string   `gorm:"uniqueIndex;not null" json:"username"` // для логики приложения
	TelegramID        *int64   `gorm:"uniqueIndex" json:"telegramId"`
	TelegramUsername  *string  `json:"telegramUsername"` // raw Telegram username
	TelegramFirstName *string  `json:"telegramFirstName"`
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
	GetByUsername(ctx context.Context, username string) (User, error)
}

type UserService interface {
	CreateUserWithWallet(ctx context.Context, tgUser initdata.User) (*User, error)
	GetUserInfoByID(ctx context.Context, userID uint) (*User, error)
	GetUserInfoByEmail(ctx context.Context, email string) (*User, error)
	GetUserInfoByTelegramID(ctx context.Context, tgID int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
}
