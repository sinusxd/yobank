package service

import (
	"context"
	"fmt"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"gorm.io/gorm"
	"yobank/domain"
)

type userService struct {
	DB         *gorm.DB
	UserRepo   domain.UserRepository
	WalletRepo domain.WalletRepository
}

func NewUserService(
	db *gorm.DB,
	userRepo domain.UserRepository,
	walletRepo domain.WalletRepository,
) domain.UserService {
	return &userService{
		DB:         db,
		UserRepo:   userRepo,
		WalletRepo: walletRepo,
	}
}

func (s *userService) CreateUserWithWallet(ctx context.Context, tgUser initdata.User) (*domain.User, error) {
	var user *domain.User

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		u, err := s.UserRepo.GetByTelegramIDWithTx(tx, tgUser.ID)
		if err == nil && u != nil {
			user = u
			return nil
		}

		user = &domain.User{
			TelegramID:        &tgUser.ID,
			TelegramUsername:  &tgUser.Username,
			TelegramFirstName: &tgUser.FirstName,
			Username:          tgUser.Username,
			Email:             nil,
		}
		if err := s.UserRepo.CreateWithTx(tx, user); err != nil {
			return fmt.Errorf("cannot create user: %w", err)
		}

		number := s.WalletRepo.GenerateWalletNumber()

		wallet := &domain.Wallet{
			UserID:   user.ID,
			Number:   number,
			Balance:  0,
			Currency: "RUB",
			Status:   "active",
		}
		if err := s.WalletRepo.CreateWithTx(tx, wallet); err != nil {
			return fmt.Errorf("cannot create wallet: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
