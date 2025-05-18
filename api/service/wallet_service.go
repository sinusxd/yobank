package service

import (
	"context"
	"errors"
	"time"
	"yobank/domain"

	"gorm.io/gorm"
)

type walletService struct {
	walletRepository domain.WalletRepository
	contextTimeout   time.Duration
}

func NewWalletService(walletRepository domain.WalletRepository, timeout time.Duration) domain.WalletService {
	return &walletService{
		walletRepository: walletRepository,
		contextTimeout:   timeout,
	}
}

func (w *walletService) GetWalletByUserID(ctx context.Context, userID uint) ([]domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()
	return w.walletRepository.GetByUserID(ctx, userID)
}

func (w *walletService) CreateWallet(ctx context.Context, userID uint, currency string) (domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	// Список поддерживаемых валют
	supported := map[string]bool{
		"RUB": true,
		"USD": true,
		"EUR": true,
		"CNY": true,
	}
	if !supported[currency] {
		return domain.Wallet{}, errors.New("unsupported currency")
	}

	wallet := domain.Wallet{
		UserID:    userID,
		Number:    w.walletRepository.GenerateWalletNumber(),
		Balance:   0,
		Currency:  currency,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	if err := w.walletRepository.Create(ctx, &wallet); err != nil {
		return domain.Wallet{}, err
	}

	return wallet, nil
}

func (w *walletService) InitWalletIfNotExists(ctx context.Context, userID uint) (domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	wallets, err := w.walletRepository.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || len(wallets) == 0 {
			return w.CreateWallet(ctx, userID, "RUB")
		}
		return domain.Wallet{}, err
	}

	return wallets[0], nil
}
