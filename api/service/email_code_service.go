package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"yobank/domain"
	"yobank/internal/mailer"
)

type emailCodeService struct {
	codeRepository domain.EmailCodeRepository
	mailer         *mailer.GoMailer
	contextTimeout time.Duration
}

func NewEmailCodeService(repo domain.EmailCodeRepository, mailer *mailer.GoMailer, timeout time.Duration) domain.EmailCodeService {
	return &emailCodeService{
		codeRepository: repo,
		mailer:         mailer,
		contextTimeout: timeout,
	}
}

func (s *emailCodeService) RequestLoginCode(ctx context.Context, email string) error {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiresAt := time.Now().Add(10 * time.Minute)

	loginCode := domain.EmailLoginCode{
		Email:     email,
		Code:      code,
		ExpiresAt: expiresAt,
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if err := s.codeRepository.Create(ctx, loginCode); err != nil {
		return err
	}

	go func() {
		_ = s.mailer.SendLoginCode(email, code)
	}()

	return nil
}

func (s *emailCodeService) VerifyLoginCode(ctx context.Context, email, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	loginCode, err := s.codeRepository.Verify(ctx, email, code)
	if err != nil {
		return false, err
	}

	_ = s.codeRepository.Delete(ctx, loginCode)

	return true, nil
}
