package repository

import (
	"context"
	"gorm.io/gorm"
	"yobank/domain"
)

type transferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) domain.TransferRepository {
	return &transferRepository{db}
}

func (r *transferRepository) CreateWithTx(tx *gorm.DB, transfer *domain.Transfer) error {
	return tx.Create(transfer).Error
}

func (r *transferRepository) GetByWalletID(ctx context.Context, walletID uint) ([]domain.Transfer, error) {
	var transfers []domain.Transfer
	err := r.db.WithContext(ctx).
		Where("sender_wallet_id = ? OR receiver_wallet_id = ?", walletID, walletID).
		Order("created_at DESC").
		Find(&transfers).Error
	if err != nil {
		return nil, err
	}
	return transfers, nil
}
