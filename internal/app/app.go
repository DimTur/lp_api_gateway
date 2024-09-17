package app

import (
	"log/slog"
	"time"

	httpapp "github.com/DimTur/lp_api_gateway/internal/app/http"
	lpgrpc "github.com/DimTur/lp_api_gateway/internal/clients/lp/grpc"
	ssogrpc "github.com/DimTur/lp_api_gateway/internal/clients/sso/grpc"
	"github.com/DimTur/lp_api_gateway/internal/handlers"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	HTTPServer *httpapp.APIServer
}

func NewApp(
	httpAddr string,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	iddleTimeout time.Duration,
	authGRPCClient ssogrpc.Client,
	lpGRPCClient lpgrpc.Client,
	logger *slog.Logger,
	traceProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) (*App, error) {
	routerConfigurator := handlers.NewChiRouterConfigurator(
		authGRPCClient,
		lpGRPCClient,
		logger,
		traceProvider,
		meterProvider,
	)
	router := routerConfigurator.ConfigureRouter()

	httpServer, err := httpapp.NewHTTPServer(
		httpAddr,
		router,
		readTimeout,
		writeTimeout,
		iddleTimeout,
		logger,
	)
	if err != nil {
		logger.Error("failed to create server", slog.Any("err", err))
		return nil, err
	}

	return &App{
		HTTPServer: httpServer,
	}, nil
}
