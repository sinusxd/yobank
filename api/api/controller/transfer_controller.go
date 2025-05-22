package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"yobank/domain"
)

type TransferController struct {
	TransferService domain.TransferService
}

type transferRequest struct {
	SenderWalletID   uint  `json:"senderWalletId"`
	ReceiverWalletID uint  `json:"receiverWalletId"`
	Amount           int64 `json:"amount"`
}

func (tc *TransferController) CreateTransfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Некорректный запрос"})
		return
	}

	transfer, err := tc.TransferService.MakeTransfer(c.Request.Context(), req.SenderWalletID, req.ReceiverWalletID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transfer)
}

func (tc *TransferController) GetTransferHistory(c *gin.Context) {
	idStr := c.Param("walletId")
	walletID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Некорректный ID кошелька"})
		return
	}

	transfers, err := tc.TransferService.GetHistoryByWalletID(c.Request.Context(), uint(walletID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, transfers)
}

func (tc *TransferController) GetReceiverUsername(c *gin.Context) {
	idStr := c.Param("walletId")
	walletID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Некорректный ID кошелька"})
		return
	}

	user, err := tc.TransferService.GetUserInfoByWalletID(c.Request.Context(), uint(walletID))
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"username": user.Username, "avatarUrl": user.AvatarURL})
}
