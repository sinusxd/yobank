package route

import (
	"yobank/api/controller"
	"yobank/domain"

	"github.com/gin-gonic/gin"
)

func NewWalletRouter(walletService domain.WalletService, authGroup *gin.RouterGroup) {
	wc := &controller.WalletController{
		WalletService: walletService,
	}

	// Маршруты для работы с кошельком (требуют авторизации)
	authGroup.GET("/wallet", wc.GetUserWallet)
	authGroup.GET("/wallets/user/:id", wc.GetByUserID)
	authGroup.POST("/wallet", wc.CreateWallet)
	authGroup.POST("/wallet/init", wc.InitWallet)
	authGroup.POST("/wallet/topup", wc.TopUpWallet)
}
