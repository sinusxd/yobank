package service

import (
	"context"
	"time"

	"yobank/domain"
	"yobank/internal/tokenutil"
)

type loginService struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewLoginService(userRepository domain.UserRepository, timeout time.Duration) domain.LoginService {
	return &loginService{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (lu *loginService) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.GetByEmail(ctx, email)
}

func (lu *loginService) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}

func (lu *loginService) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}
