package coingecko_api

import (
	"context"
	desc "github.com/fidesy-pay/external-api/pkg/coingecko-api"
	"google.golang.org/grpc"
)

type (
	Implementation struct {
		desc.UnimplementedCoinGeckoAPIServer

		coinGeckoAPI CoinGeckoAPI
	}

	CoinGeckoAPI interface {
		GetPrice(ctx context.Context, symbol string) (float64, error)
	}
)

func New(coinGeckoAPI CoinGeckoAPI) *Implementation {
	return &Implementation{
		coinGeckoAPI: coinGeckoAPI,
	}
}

func (i *Implementation) GetDescription() *grpc.ServiceDesc {
	return &desc.CoinGeckoAPI_ServiceDesc
}
