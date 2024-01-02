package coingeckoapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fidesy-pay/external-api/internal/config"
	"io"
	"net/http"
)

const (
	coingeckoAPI = "https://api.coingecko.com/api/v3"
)

var (
	ErrExceedRateLimit = errors.New("you've exceeded the rate limit")
)

type (
	Service struct {
		httpClient *http.Client
		apiKey     string
	}
)

func New() *Service {
	s := &Service{
		httpClient: http.DefaultClient,
	}

	s.apiKey = config.Get(config.CoinGeckoAPIKey).(string)

	return s
}

func (s *Service) GetPrice(ctx context.Context, symbol string) (float64, error) {
	url := fmt.Sprintf(
		"%s/simple/price?ids=%s&vs_currencies=usd&x_cg_demo_api_key=%s",
		coingeckoAPI,
		symbol,
		s.apiKey,
	)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	response, err := s.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusTooManyRequests {
		return 0, ErrExceedRateLimit
	}

	body, _ := io.ReadAll(response.Body)
	var data map[string]CoinGeckoResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}

	if cryptoData, ok := data[symbol]; ok {
		return cryptoData.Usd, nil
	}

	return 0, fmt.Errorf("price data not available for symbol: %s", symbol)
}
