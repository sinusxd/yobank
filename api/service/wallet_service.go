package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"yobank/domain"
	"yobank/internal/telegram"

	"gorm.io/gorm"
)

type walletService struct {
	walletRepository domain.WalletRepository
	userRepository   domain.UserRepository
	contextTimeout   time.Duration
}

func NewWalletService(walletRepository domain.WalletRepository, userRepository domain.UserRepository, timeout time.Duration) domain.WalletService {
	return &walletService{
		walletRepository: walletRepository,
		userRepository:   userRepository,
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
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Wallet{}, err
	}
	if len(wallets) == 0 {
		created, err := w.CreateWallet(ctx, userID, "RUB")
		if err != nil {
			return domain.Wallet{}, err
		}

		created, err = w.CreateWallet(ctx, userID, "USD")
		if err != nil {
			return domain.Wallet{}, err
		}

		created, err = w.CreateWallet(ctx, userID, "EUR")
		if err != nil {
			return domain.Wallet{}, err
		}

		created, err = w.CreateWallet(ctx, userID, "CNY")
		if err != nil {
			return domain.Wallet{}, err
		}

		return created, nil
	}
	return wallets[0], nil
}

func (w *walletService) TopUpWallet(ctx context.Context, userID uint, currency string, amount int64) (domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	if amount <= 0 {
		return domain.Wallet{}, errors.New("amount must be positive")
	}

	wallets, err := w.walletRepository.GetByUserID(ctx, userID)
	if err != nil {
		return domain.Wallet{}, err
	}

	var target *domain.Wallet
	for i := range wallets {
		if wallets[i].Currency == currency {
			target = &wallets[i]
			break
		}
	}
	if target == nil {
		return domain.Wallet{}, errors.New("кошелек с указанной валютой не найден")
	}
	if target.Status != "active" {
		return domain.Wallet{}, errors.New("кошелек неактивен")
	}

	target.Balance += amount

	if err := w.walletRepository.Update(ctx, target); err != nil {
		return domain.Wallet{}, err
	}

	// Уведомление пользователя
	user, err := w.userRepository.GetByID(ctx, fmt.Sprint(userID))
	if err == nil && user.TelegramID != nil {
		telegram.NotifyTopUp(*user.TelegramID, amount, currency)
	}

	return *target, nil
}
