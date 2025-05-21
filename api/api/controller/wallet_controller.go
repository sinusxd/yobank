package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"yobank/domain"
)

type WalletController struct {
	WalletService domain.WalletService
}

func (wc *WalletController) GetUserWallet(c *gin.Context) {
	id, ok := c.Get("x-user-id")
	if ok != true {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
		return
	}

	userID, err := strconv.ParseUint(id.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
		return
	}

	wallets, err := wc.WalletService.GetWalletByUserID(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Кошелек не найден"})
		return
	}

	walletsResponse := domain.WalletsToResponse(wallets)

	c.JSON(http.StatusOK, walletsResponse)
}

// InitWallet проверяет наличие кошелька у пользователя и создает его, если отсутствует
func (wc *WalletController) InitWallet(c *gin.Context) {
	// Получаем ID пользователя из JWT токена
	id, ok := c.Get("x-user-id")
	if ok != true {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
		return
	}

	userID, err := strconv.ParseUint(id.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
		return
	}

	wallet, err := wc.WalletService.InitWalletIfNotExists(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Ошибка при инициализации кошелька"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (wc *WalletController) CreateWallet(c *gin.Context) {
	id, ok := c.Get("x-user-id")
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
		return
	}

	userID, err := strconv.ParseUint(id.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
		return
	}

	var req struct {
		Currency string `json:"currency" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Currency == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Поле currency обязательно"})
		return
	}

	ctx := c.Request.Context()
	wallet, err := wc.WalletService.CreateWallet(ctx, uint(userID), req.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Ошибка при создании кошелька"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (wc *WalletController) TopUpWallet(c *gin.Context) {
	id, ok := c.Get("x-user-id")
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
		return
	}

	userID, err := strconv.ParseUint(id.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
		return
	}

	// Структура запроса — валюта и сумма в копейках
	var req struct {
		Currency string `json:"currency" binding:"required"`
		Amount   int64  `json:"amount" binding:"required"` // в копейках
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 || req.Currency == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Некорректные данные для пополнения"})
		return
	}

	updatedWallet, err := wc.WalletService.TopUpWallet(c.Request.Context(), uint(userID), req.Currency, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedWallet)
}

func (wc *WalletController) GetByUserID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Некорректный ID"})
		return
	}

	wallets, err := wc.WalletService.GetWalletByUserID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Кошельки не найдены"})
		return
	}

	c.JSON(http.StatusOK, domain.WalletsToResponse(wallets))
}
