package bootstrap

import (
	"time"

	"yobank/domain"
	"yobank/internal/mailer"
	"yobank/repository"
	"yobank/service"
)

type Services struct {
	Login     domain.LoginService
	EmailCode domain.EmailCodeService
}

type Repositories struct {
	User      domain.UserRepository
	EmailCode domain.EmailCodeRepository
}

type Container struct {
	Services Services
	Repos    Repositories
}

func BuildContainer(app Application) Container {
	timeout := 5 * time.Second
	db := app.DB
	cfg := app.Env

	userRepo := repository.NewUserRepository(db)
	emailCodeRepo := repository.NewEmailCodeRepository(db)

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

	return Container{
		Services: Services{
			EmailCode: emailCodeService,
			Login:     loginService,
		},
		Repos: Repositories{
			User:      userRepo,
			EmailCode: emailCodeRepo,
		},
	}
}
