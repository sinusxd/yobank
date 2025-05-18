package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Wallet struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index;not null"` // FK на пользователя
	Number    string `gorm:"size:24;uniqueIndex;not null"`
	Balance   int64  `gorm:"not null;default:0"`
	Currency  string `gorm:"size:8;default:'RUB'"`
	Status    string `gorm:"size:16;default:'active'"`
	CreatedAt time.Time

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // для жадной загрузки
}

type WalletResponse struct {
	ID        uint   `json:"id"`
	Number    string `json:"number"`
	Balance   int64  `json:"balance"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
	CreatedAt time.Time
}

func WalletToResponse(wallet Wallet) WalletResponse {
	return WalletResponse{
		ID:        wallet.ID,
		Number:    wallet.Number,
		Currency:  wallet.Currency,
		Balance:   wallet.Balance,
		Status:    wallet.Status,
		CreatedAt: wallet.CreatedAt,
	}
}

func WalletsToResponse(wallets []Wallet) []WalletResponse {
	response := make([]WalletResponse, len(wallets))
	for i, wallet := range wallets {
		response[i] = WalletToResponse(wallet)
	}
	return response
}

type WalletRepository interface {
	GetByUserID(ctx context.Context, userID uint) ([]Wallet, error)
	Create(ctx context.Context, wallet *Wallet) error
	CreateWithTx(tx *gorm.DB, wallet *Wallet) error
	Update(ctx context.Context, wallet *Wallet) error
	GenerateWalletNumber() string
}

type WalletService interface {
	CreateWallet(ctx context.Context, userID uint, currency string) (Wallet, error)
	GetWalletByUserID(ctx context.Context, userID uint) ([]Wallet, error)
	InitWalletIfNotExists(ctx context.Context, userID uint) (Wallet, error)
}
