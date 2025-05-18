package domain

import "time"

type Transaction struct {
	ID             uint    `gorm:"primaryKey"`
	FromUserID     *uint   `gorm:"index"`    // Отправитель (nullable для deposit)
	ToUserID       *uint   `gorm:"index"`    // Получатель (nullable для withdraw)
	Amount         int64   `gorm:"not null"` // В копейках!
	Currency       string  `gorm:"size:8;default:'RUB'"`
	Type           string  `gorm:"size:16;not null"` // transfer, deposit, withdrawal
	Status         string  `gorm:"size:16;not null"` // pending, completed, failed
	Description    *string `gorm:"size:255"`
	IdempotencyKey *string `gorm:"size:64;uniqueIndex"`
	RelatedTxnID   *uint   `gorm:"index"`
	CreatedAt      time.Time
}
