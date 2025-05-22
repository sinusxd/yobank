package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Transfer struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	SenderWalletID   uint      `gorm:"index;not null" json:"senderWalletId"`
	ReceiverWalletID uint      `gorm:"index;not null" json:"receiverWalletId"`
	Amount           int64     `gorm:"not null" json:"amount"`
	Currency         string    `gorm:"size:8;not null" json:"currency"`
	CreatedAt        time.Time `json:"createdAt"`
}

type TransferResponse struct {
	ID        uint      `json:"id"`
	From      string    `json:"from"` // номер отправителя
	To        string    `json:"to"`   // номер получателя
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"createdAt"`
}

func TransferToResponse(t Transfer, sender Wallet, receiver Wallet) TransferResponse {
	return TransferResponse{
		ID:        t.ID,
		From:      sender.Number,
		To:        receiver.Number,
		Amount:    t.Amount,
		Currency:  t.Currency,
		CreatedAt: t.CreatedAt,
	}
}

type TransferRepository interface {
	CreateWithTx(tx *gorm.DB, transfer *Transfer) error
	GetByWalletID(ctx context.Context, walletID uint) ([]Transfer, error)
}

type TransferService interface {
	MakeTransfer(ctx context.Context, senderWalletID, receiverWalletID uint, amount int64) (Transfer, error)
	GetHistoryByWalletID(ctx context.Context, walletID uint) ([]Transfer, error)
	GetUserInfoByWalletID(ctx context.Context, walletID uint) (*User, error)
}
