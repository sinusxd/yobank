package repository

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"yobank/domain"

	"gorm.io/gorm"
)

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) domain.WalletRepository {
	return &walletRepository{db}
}

func (r *walletRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Wallet, error) {
	var wallets []domain.Wallet
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&wallets)
	if result.Error != nil {
		return nil, result.Error
	}
	return wallets, nil
}

func (r *walletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *walletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}

func (r *walletRepository) GenerateWalletNumber() string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	walletNumber := fmt.Sprintf("%04d-%04d-%04d-%04d",
		random.Intn(10000),
		random.Intn(10000),
		random.Intn(10000),
		random.Intn(10000))
	return walletNumber
}

func (r *walletRepository) CreateWithTx(tx *gorm.DB, wallet *domain.Wallet) error {
	return tx.Create(wallet).Error
}

func (r *walletRepository) GetByIDTx(tx *gorm.DB, id uint) (*domain.Wallet, error) {
	var wallet domain.Wallet
	if err := tx.First(&wallet, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) UpdateWithTx(tx *gorm.DB, wallet *domain.Wallet) error {
	return tx.Save(wallet).Error
}

func (r *walletRepository) GetByNumber(ctx context.Context, number string) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.WithContext(ctx).
		Where("number = ?", number).
		First(&wallet).Error

	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) GetByID(ctx context.Context, id uint) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}
