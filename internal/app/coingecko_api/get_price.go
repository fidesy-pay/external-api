package coingecko_api

import (
	"context"
	"errors"
	coingeckoapi "github.com/fidesy-pay/external-api/internal/pkg/coingecko-api"
	desc "github.com/fidesy-pay/external-api/pkg/coingecko-api"
	validation "github.com/go-ozzo/ozzo-validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetPrice(ctx context.Context, req *desc.GetPriceRequest) (*desc.GetPriceResponse, error) {
	err := validateGetPriceRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	price, err := i.coinGeckoAPI.GetPrice(ctx, req.GetSymbol())
	if err != nil {
		if errors.Is(err, coingeckoapi.ErrExceededRateLimit) {
			return nil, status.Errorf(codes.ResourceExhausted, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "coinGeckoAPI.GetPrice: %v", err)
	}

	return &desc.GetPriceResponse{
		PriceUsd: price,
	}, nil
}

func validateGetPriceRequest(req *desc.GetPriceRequest) error {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.Symbol, validation.Required),
	)
	return err
}
