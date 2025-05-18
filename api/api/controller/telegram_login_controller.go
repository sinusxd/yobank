package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	initdata "github.com/telegram-mini-apps/init-data-golang"

	"yobank/bootstrap"
	"yobank/domain"
)

type TelegramLoginController struct {
	UserService  domain.UserService
	Env          *bootstrap.Env
	LoginService domain.LoginService
}

type TelegramLoginRequest struct {
	InitData string `json:"init_data" binding:"required"`
}

func (tc *TelegramLoginController) Login(c *gin.Context) {
	var req TelegramLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if err := initdata.Validate(req.InitData, tc.Env.TelegramBotToken, 10*time.Minute); err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid init data"})
		return
	}

	parsed, err := initdata.Parse(req.InitData)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Cannot parse init data"})
		return
	}

	tgUser := parsed.User

	user, err := tc.UserService.CreateUserWithWallet(c.Request.Context(), tgUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Ошибка при создании пользователя или кошелька: " + err.Error()})
		return
	}

	accessToken, err := tc.LoginService.CreateAccessToken(user, tc.Env.AccessTokenSecret, tc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshToken, err := tc.LoginService.CreateRefreshToken(user, tc.Env.RefreshTokenSecret, tc.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
