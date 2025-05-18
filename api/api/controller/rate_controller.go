package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"yobank/domain"
)

type RateController struct {
	RateService domain.RateService
}

func (rc *RateController) GetCurrencyRate(c *gin.Context) {
	currency := c.Param("currency")
	if currency == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "currency is required"})
		return
	}

	dateStr := c.Query("date") // формат "2006-01-02"

	if dateStr == "" {
		// Нет даты — ищем самый свежий курс
		rate, err := rc.RateService.GetLatestRate(c.Request.Context(), currency)
		if err != nil {
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "rate not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"currency": rate.Currency,
			"value":    rate.Value,
			"date":     rate.Date.Format("2006-01-02 15:04:05"),
		})
		return
	}

	// Если дата указана — ищем все за этот день, возвращаем последний за день
	from, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid date format"})
		return
	}
	to := from.Add(24 * time.Hour)

	rates, err := rc.RateService.GetRatesHistory(c.Request.Context(), currency, from.UTC(), to.UTC())
	if err != nil || len(rates) == 0 {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "rate not found"})
		return
	}

	rate := rates[len(rates)-1]
	c.JSON(http.StatusOK, gin.H{
		"currency": rate.Currency,
		"value":    rate.Value,
		"date":     rate.Date.Format("2006-01-02 15:04:05"),
	})
}
