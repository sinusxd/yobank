package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
	"yobank/domain"
)

type rateRepository struct {
	db *gorm.DB
}

func NewRateRepository(db *gorm.DB) domain.RateRepository {
	return &rateRepository{db}
}

func (r *rateRepository) Save(ctx context.Context, rate *domain.Rate) error {
	return r.db.WithContext(ctx).Create(rate).Error
}

func (r *rateRepository) GetByCurrencyAndPeriod(ctx context.Context, currency string, from, to time.Time) ([]domain.Rate, error) {
	var rates []domain.Rate
	result := r.db.WithContext(ctx).
		Where("currency = ? AND date >= ? AND date <= ?", currency, from, to).
		Order("date").
		Find(&rates)
	if result.Error != nil {
		return nil, result.Error
	}
	return rates, nil
}

func (r *rateRepository) GetLatestRate(ctx context.Context, currency string) (domain.Rate, error) {
	var rate domain.Rate
	err := r.db.WithContext(ctx).
		Where("currency = ?", currency).
		Order("date desc").
		Limit(1).
		First(&rate).Error
	return rate, err
}
