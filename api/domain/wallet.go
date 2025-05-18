package domain

import "time"

type Wallet struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"uniqueIndex;not null"`
	Number    string `gorm:"size:24;uniqueIndex;not null"`
	Balance   int64  `gorm:"not null;default:0"`
	Currency  string `gorm:"size:8;default:'RUB'"`
	Status    string `gorm:"size:16;default:'active'"`
	CreatedAt time.Time
}

type WalletRepository interface {
}

type WalletService interface {
}
