package cacheapi

import (
	"context"
	"fmt"
	"github.com/fidesy-pay/external-api/internal/config"
	"time"
)

type (
	Service struct {
		cache        Cache
		coinGeckoAPI CoinGeckoAPI
	}

	Cache interface {
		Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error
		Get(ctx context.Context, key string, dst interface{}) (bool, error)
	}

	CoinGeckoAPI interface {
		GetPrice(ctx context.Context, symbol string) (float64, error)
	}
)

func New(
	cache Cache,
	coinGeckoAPI CoinGeckoAPI,
) *Service {
	return &Service{
		cache:        cache,
		coinGeckoAPI: coinGeckoAPI,
	}
}

func (s *Service) GetPrice(ctx context.Context, symbol string) (float64, error) {
	var price float64
	found, err := s.cache.Get(ctx, symbol, &price)
	if err != nil {
		return 0, fmt.Errorf("cache.Get: %w", err)
	}

	if found {
		return price, nil
	}

	price, err = s.coinGeckoAPI.GetPrice(ctx, symbol)
	if err != nil {
		return 0, fmt.Errorf("coinGeckoAPI.GetPrice: %w", err)
	}

	err = s.cache.Set(ctx, symbol, price, config.Get(config.CacheTTL).(time.Duration))
	if err != nil {
		return 0, fmt.Errorf("cache.Set: %w", err)
	}

	return price, nil
}
