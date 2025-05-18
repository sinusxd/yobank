package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"yobank/domain"
)

type rateService struct {
	rateRepository domain.RateRepository
	contextTimeout time.Duration
}

func NewRateService(rateRepository domain.RateRepository, timeout time.Duration) domain.RateService {
	return &rateService{
		rateRepository: rateRepository,
		contextTimeout: timeout,
	}
}

func (s *rateService) FetchAndSaveRates(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	resp, err := http.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data struct {
		Date   string `json:"Date"`
		Valute map[string]struct {
			Value float64 `json:"Value"`
		} `json:"Valute"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	t, _ := time.Parse(time.RFC3339, data.Date)
	for _, code := range []string{"USD", "EUR", "CNY"} {
		rate := &domain.Rate{
			Currency: code,
			Value:    data.Valute[code].Value,
			Date:     t,
		}
		_ = s.rateRepository.Save(ctx, rate)
	}
	return nil
}

func (s *rateService) GetRatesHistory(ctx context.Context, currency string, from, to time.Time) ([]domain.Rate, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	return s.rateRepository.GetByCurrencyAndPeriod(ctx, currency, from, to)
}

func (s *rateService) GetLatestRate(ctx context.Context, currency string) (domain.Rate, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	return s.rateRepository.GetLatestRate(ctx, currency)
}

func (s *rateService) StartScheduler() {
	go func() {
		for {
			ctx := context.Background()
			s.FetchAndSaveRates(ctx)
			time.Sleep(time.Hour)
		}
	}()
}
