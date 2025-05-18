package domain

import (
	"context"
	"time"
)

type Rate struct {
	ID        uint `gorm:"primaryKey"`
	Currency  string
	Value     float64
	Date      time.Time
	CreatedAt time.Time
}

type RateRepository interface {
	Save(ctx context.Context, rate *Rate) error
	GetLatestRate(ctx context.Context, currency string) (Rate, error)
	GetByCurrencyAndPeriod(ctx context.Context, currency string, from, to time.Time) ([]Rate, error)
}

type RateService interface {
	FetchAndSaveRates(ctx context.Context) error
	GetLatestRate(ctx context.Context, currency string) (Rate, error)
	GetRatesHistory(ctx context.Context, currency string, from, to time.Time) ([]Rate, error)
	StartScheduler()
}
