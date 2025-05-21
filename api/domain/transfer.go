package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Transfer struct {
	ID               uint   `gorm:"primaryKey"`
	SenderWalletID   uint   `gorm:"index;not null"`
	ReceiverWalletID uint   `gorm:"index;not null"`
	Amount           int64  `gorm:"not null"`
	Currency         string `gorm:"size:8;not null"`
	CreatedAt        time.Time
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
}
