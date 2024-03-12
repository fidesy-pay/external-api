package coingeckoapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	coingeckoAPI = "https://api.coingecko.com/api/v3"
)

var (
	ErrExceededRateLimit = errors.New("you've exceeded the rate limit")
)

type (
	Service struct {
		httpClient *http.Client
		apiKey     string
	}
)

type Option func(s *Service)

func WithAPIKey(apiKey string) Option {
	return func(s *Service) {
		s.apiKey = apiKey
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(s *Service) {
		s.httpClient = httpClient
	}
}

func New(opts ...Option) *Service {
	s := &Service{
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(s)
	}

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
		return 0, ErrExceededRateLimit
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
