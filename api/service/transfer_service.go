package service

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
	"yobank/domain"
	"yobank/internal/telegram"
)

type transferService struct {
	db               *gorm.DB
	walletRepository domain.WalletRepository
	transferRepo     domain.TransferRepository
	userRepo         domain.UserRepository
	contextTimeout   time.Duration
}

func NewTransferService(
	db *gorm.DB,
	walletRepo domain.WalletRepository,
	transferRepo domain.TransferRepository,
	userRepo domain.UserRepository,
	timeout time.Duration,
) domain.TransferService {
	return &transferService{
		db:               db,
		walletRepository: walletRepo,
		transferRepo:     transferRepo,
		userRepo:         userRepo,
		contextTimeout:   timeout,
	}
}

func (s *transferService) MakeTransfer(ctx context.Context, senderWalletID, receiverWalletID uint, amount int64) (domain.Transfer, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	var transfer domain.Transfer

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sender, err := s.walletRepository.GetByIDTx(tx, senderWalletID)
		if err != nil {
			return fmt.Errorf("отправитель не найден: %w", err)
		}

		receiver, err := s.walletRepository.GetByIDTx(tx, receiverWalletID)
		if err != nil {
			return fmt.Errorf("получатель не найден: %w", err)
		}

		if sender.Currency != receiver.Currency {
			return fmt.Errorf("кошельки в разных валютах")
		}

		if sender.Balance < amount {
			return fmt.Errorf("недостаточно средств")
		}

		sender.Balance -= amount
		receiver.Balance += amount

		if err := s.walletRepository.UpdateWithTx(tx, sender); err != nil {
			return err
		}
		if err := s.walletRepository.UpdateWithTx(tx, receiver); err != nil {
			return err
		}

		transfer = domain.Transfer{
			SenderWalletID:   senderWalletID,
			ReceiverWalletID: receiverWalletID,
			Amount:           amount,
			Currency:         sender.Currency,
			CreatedAt:        time.Now(),
		}

		if err := s.transferRepo.CreateWithTx(tx, &transfer); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return transfer, err
	}

	// Оповещение получателя
	receiverWallet, err := s.walletRepository.GetByID(ctx, transfer.ReceiverWalletID)
	if err == nil {
		receiverUser, err := s.userRepo.GetByID(ctx, strconv.Itoa(int(receiverWallet.UserID)))
		if err == nil && receiverUser.TelegramID != nil {
			senderWallet, _ := s.walletRepository.GetByID(ctx, transfer.SenderWalletID)
			senderUser, _ := s.userRepo.GetByID(ctx, strconv.Itoa(int(senderWallet.UserID)))

			senderUsername := "неизвестно"
			if senderUser.Username != "" {
				senderUsername = senderUser.Username
			}

			telegram.NotifyTransfer(*receiverUser.TelegramID, senderUsername, transfer.Amount, transfer.Currency, senderUser.TelegramID != nil)
		}
	}

	return transfer, nil
}

func (s *transferService) GetHistoryByWalletID(ctx context.Context, walletID uint) ([]domain.Transfer, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	return s.transferRepo.GetByWalletID(ctx, walletID)
}

func (s *transferService) GetUserInfoByWalletID(ctx context.Context, walletID uint) (*domain.User, error) {
	wallet, err := s.walletRepository.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, strconv.Itoa(int(wallet.UserID)))
	if err != nil {
		return nil, err
	}

	return user, nil
}
