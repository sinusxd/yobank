package route

import (
	"github.com/gin-gonic/gin"
	"yobank/api/controller"
	"yobank/domain"
)

func NewTransferRouter(transferService domain.TransferService, group *gin.RouterGroup) {
	tc := &controller.TransferController{
		TransferService: transferService,
	}

	group.POST("/transfers", tc.CreateTransfer)
	group.GET("/transfers/wallet/:walletId", tc.GetTransferHistory)
	group.GET("/transfers/username/:walletId", tc.GetReceiverUsername)

}
