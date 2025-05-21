package service

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
	"yobank/domain"
)

type transferService struct {
	db               *gorm.DB
	walletRepository domain.WalletRepository
	transferRepo     domain.TransferRepository
	contextTimeout   time.Duration
}

func NewTransferService(
	db *gorm.DB,
	walletRepo domain.WalletRepository,
	transferRepo domain.TransferRepository,
	timeout time.Duration,
) domain.TransferService {
	return &transferService{
		db:               db,
		walletRepository: walletRepo,
		transferRepo:     transferRepo,
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

	return transfer, err
}

func (s *transferService) GetHistoryByWalletID(ctx context.Context, walletID uint) ([]domain.Transfer, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	return s.transferRepo.GetByWalletID(ctx, walletID)
}
