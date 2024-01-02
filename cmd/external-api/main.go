package main

import (
	"context"
	"github.com/fidesy-pay/external-api/internal/app/coingecko_api"
	"github.com/fidesy-pay/external-api/internal/config"
	coingeckoapi "github.com/fidesy-pay/external-api/internal/pkg/coingecko-api"
	"github.com/fidesyx/platform/pkg/scratch"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer cancel()

	err := config.Init()
	if err != nil {
		log.Fatalf("config.Init: %v", err)
	}

	coinGeckoAPI := coingeckoapi.New()

	impl := coingecko_api.New(coinGeckoAPI)

	app, err := scratch.New(ctx)
	if err != nil {
		log.Fatalf("scratch.New: %v", err)
	}

	err = app.Run(ctx, impl)
	if err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}
