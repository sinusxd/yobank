package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"yobank/domain"
)

type emailCodeRepository struct {
	db *gorm.DB
}

func NewEmailCodeRepository(db *gorm.DB) domain.EmailCodeRepository {
	return &emailCodeRepository{
		db: db,
	}
}

func (r *emailCodeRepository) Create(ctx context.Context, code domain.EmailLoginCode) error {
	return r.db.WithContext(ctx).Create(&code).Error
}

func (r *emailCodeRepository) Verify(ctx context.Context, email string, code string) (domain.EmailLoginCode, error) {
	var result domain.EmailLoginCode
	err := r.db.WithContext(ctx).Where(
		"email = ? AND code = ? AND expires_at > ?", email, code, time.Now().UTC(),
	).First(&result).Error

	return result, err
}

func (r *emailCodeRepository) Delete(ctx context.Context, code domain.EmailLoginCode) error {
	return r.db.WithContext(ctx).Delete(&code).Error
}
