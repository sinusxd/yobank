package bootstrap

import (
	"gorm.io/gorm"
	"time"

	"yobank/domain"
	"yobank/internal/mailer"
	"yobank/repository"
	"yobank/service"
)

type Services struct {
	Login     domain.LoginService
	EmailCode domain.EmailCodeService
	Wallet    domain.WalletService
	User      domain.UserService
	Rate      domain.RateService
	Transfer  domain.TransferService
}

type Repositories struct {
	User      domain.UserRepository
	EmailCode domain.EmailCodeRepository
	Wallet    domain.WalletRepository
	Rate      domain.RateRepository
	Transfer  domain.TransferRepository
}

type Container struct {
	Services Services
	Repos    Repositories
}

func BuildContainer(db *gorm.DB, cfg *Env) Container {
	timeout := 5 * time.Second

	userRepo := repository.NewUserRepository(db)
	emailCodeRepo := repository.NewEmailCodeRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	rateRepo := repository.NewRateRepository(db)
	transferRepo := repository.NewTransferRepository(db)

	// Mailer
	mail := mailer.NewGoMailer(
		cfg.SMTPUsername,
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
	)

	// UseCases / Services
	emailCodeService := service.NewEmailCodeService(emailCodeRepo, mail, timeout)
	loginService := service.NewLoginService(userRepo, timeout)
	walletService := service.NewWalletService(walletRepo, userRepo, mail, timeout)
	userService := service.NewUserService(db, userRepo, walletRepo)
	rateService := service.NewRateService(rateRepo, timeout)
	transferService := service.NewTransferService(db, walletRepo, transferRepo, userRepo, mail, timeout)

	return Container{
		Services: Services{
			EmailCode: emailCodeService,
			Login:     loginService,
			Wallet:    walletService,
			User:      userService,
			Rate:      rateService,
			Transfer:  transferService,
		},
		Repos: Repositories{
			User:      userRepo,
			EmailCode: emailCodeRepo,
			Wallet:    walletRepo,
			Rate:      rateRepo,
			Transfer:  transferRepo,
		},
	}
}
