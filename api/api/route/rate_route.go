package route

import (
	"yobank/api/controller"
	"yobank/domain"

	"github.com/gin-gonic/gin"
)

func NewRateRouter(rateService domain.RateService, group *gin.RouterGroup) {
	rc := &controller.RateController{
		RateService: rateService,
	}
	group.GET("/rates/:currency", rc.GetCurrencyRate)
}
