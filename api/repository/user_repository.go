package repository

import (
	"context"
	"strconv"

	"gorm.io/gorm"
	"yobank/domain"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Create(c context.Context, user *domain.User) error {
	result := ur.db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *userRepository) CreateWithTx(tx *gorm.DB, user *domain.User) error {
	return tx.Create(user).Error
}

func (ur *userRepository) Fetch(c context.Context) ([]domain.User, error) {
	var users []domain.User
	result := ur.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (ur *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	var user domain.User
	result := ur.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}

func (ur *userRepository) GetByID(c context.Context, id string) (domain.User, error) {
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return domain.User{}, err
	}

	var user domain.User
	result := ur.db.First(&user, uint(uid))
	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}

func (ur *userRepository) GetByTelegramID(c context.Context, tgID int64) (domain.User, error) {
	var user domain.User
	result := ur.db.WithContext(c).Where("telegram_id = ?", tgID).First(&user)
	return user, result.Error
}

func (r *userRepository) GetByTelegramIDWithTx(tx *gorm.DB, tgID int64) (*domain.User, error) {
	var user domain.User
	if err := tx.Where("telegram_id = ?", tgID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
