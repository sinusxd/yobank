package service

import (
	"context"
	"fmt"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"gorm.io/gorm"
	"strconv"
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
			AvatarURL:         &tgUser.PhotoURL,
			Email:             nil,
		}
		if err := s.UserRepo.CreateWithTx(tx, user); err != nil {
			return fmt.Errorf("cannot create user: %w", err)
		}

		wallet := &domain.Wallet{
			UserID:   user.ID,
			Number:   s.WalletRepo.GenerateWalletNumber(),
			Balance:  0,
			Currency: "RUB",
			Status:   "active",
		}
		if err := s.WalletRepo.CreateWithTx(tx, wallet); err != nil {
			return fmt.Errorf("cannot create wallet: %w", err)
		}

		wallet = &domain.Wallet{
			UserID:   user.ID,
			Number:   s.WalletRepo.GenerateWalletNumber(),
			Balance:  0,
			Currency: "USD",
			Status:   "active",
		}
		if err := s.WalletRepo.CreateWithTx(tx, wallet); err != nil {
			return fmt.Errorf("cannot create wallet: %w", err)
		}

		wallet = &domain.Wallet{
			UserID:   user.ID,
			Number:   s.WalletRepo.GenerateWalletNumber(),
			Balance:  0,
			Currency: "EUR",
			Status:   "active",
		}
		if err := s.WalletRepo.CreateWithTx(tx, wallet); err != nil {
			return fmt.Errorf("cannot create wallet: %w", err)
		}

		wallet = &domain.Wallet{
			UserID:   user.ID,
			Number:   s.WalletRepo.GenerateWalletNumber(),
			Balance:  0,
			Currency: "CNY",
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

func (s *userService) GetUserInfoByID(ctx context.Context, userID uint) (*domain.User, error) {
	user, err := s.UserRepo.GetByID(ctx, strconv.FormatUint(uint64(userID), 10))
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}
	return user, nil
}

func (s *userService) GetUserInfoByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}
	return &user, nil
}

func (s *userService) GetUserInfoByTelegramID(ctx context.Context, tgID int64) (*domain.User, error) {
	user, err := s.UserRepo.GetByTelegramID(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}
	return &user, nil
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.UserRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}
	return &user, nil
}

func (s *userService) GetUserInfoByWalletNumber(ctx context.Context, walletNumber string) (*domain.User, error) {
	wallet, err := s.WalletRepo.GetByNumber(ctx, walletNumber)
	if err != nil {
		return nil, fmt.Errorf("кошелек не найден: %w", err)
	}
	return s.UserRepo.GetByID(ctx, strconv.FormatUint(uint64(wallet.UserID), 10))
}

func (s *userService) Update(ctx context.Context, user *domain.User) error {
	return s.UserRepo.Update(ctx, user)
}
