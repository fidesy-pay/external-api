package main

import (
	"context"
	"github.com/fidesy-pay/external-api/internal/app/coingecko_api"
	"github.com/fidesy-pay/external-api/internal/config"
	coingeckoapi "github.com/fidesy-pay/external-api/internal/pkg/coingecko-api"
	"github.com/fidesy/sdk/common/grpc"
	"github.com/fidesy/sdk/common/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	coinGeckoAPI := coingeckoapi.New(
		coingeckoapi.WithAPIKey(
			config.Get(config.CoinGeckoAPIKey).(string),
		),
		coingeckoapi.WithHTTPClient(httpClient),
	)

	impl := coingecko_api.New(coinGeckoAPI)

	server, err := grpc.NewServer(
		grpc.WithPort(os.Getenv("GRPC_PORT")),
		grpc.WithMetricsPort(os.Getenv("METRICS_PORT")),
		grpc.WithDomainNameService(ctx, "domain-name-service:10000"),
		grpc.WithGraylog("graylog:5555"),
		grpc.WithTracer("http://jaeger:14268/api/traces"),
	)
	if err != nil {
		log.Fatalf("grpc.NewServer: %v", err)
	}

	err = server.Run(ctx, impl)
	if err != nil {
		logger.Fatalf("app.Run: %v", err)
	}
}
