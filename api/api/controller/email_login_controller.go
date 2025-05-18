package controller

import (
	"net/http"
	"time"
	"yobank/bootstrap"
	"yobank/domain"
	"yobank/internal/telegram"

	"github.com/gin-gonic/gin"
)

type EmailLoginController struct {
	LoginService   domain.LoginService
	CodeService    domain.EmailCodeService
	UserRepository domain.UserRepository
	WalletService  domain.WalletService
	Env            *bootstrap.Env
}

func (lc *EmailLoginController) RequestCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if err := lc.CodeService.RequestLoginCode(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Не удалось отправить код"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Код отправлен на почту"})
}

func (lc *EmailLoginController) VerifyCode(c *gin.Context) {
	time.Sleep(5 * time.Second)
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ok, err := lc.CodeService.VerifyLoginCode(c.Request.Context(), req.Email, req.Code)
	if err != nil || !ok {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Неверный или просроченный код"})
		return
	}

	user, err := lc.LoginService.GetUserByEmail(c, req.Email)
	if err != nil {
		// пользователь не найден — создаём
		user = domain.User{Email: &req.Email, Username: req.Email}
		if err := lc.UserRepository.Create(c.Request.Context(), &user); err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Не удалось создать пользователя"})
			return
		}
	}

	// Создаем кошелек для пользователя, если его нет
	if lc.WalletService != nil {
		wallet, err := lc.WalletService.InitWalletIfNotExists(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Ошибка при инициализации кошелька"})
			return
		}
		// Можно добавить информацию о созданном кошельке в ответ, если нужно
		_ = wallet
	}

	accessToken, err := lc.LoginService.CreateAccessToken(&user, lc.Env.AccessTokenSecret, lc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshToken, err := lc.LoginService.CreateRefreshToken(&user, lc.Env.RefreshTokenSecret, lc.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (lc *EmailLoginController) ShowInitDataHandler(c *gin.Context) {
	initData, ok := telegram.CtxInitData(c.Request.Context())
	if !ok {
		c.AbortWithStatusJSON(401, gin.H{"message": "Init data not found"})
		return
	}

	c.JSON(200, gin.H{
		"user_id":    initData.User.ID,
		"first_name": initData.User.FirstName,
		"username":   initData.User.Username,
		"auth_date":  initData.AuthDate(),
	})
}
