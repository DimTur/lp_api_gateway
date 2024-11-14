package app

import (
	"log/slog"
	"time"

	httpapp "github.com/DimTur/lp_api_gateway/internal/app/http"
	"github.com/DimTur/lp_api_gateway/internal/handlers"
	lpservice "github.com/DimTur/lp_api_gateway/internal/services/lp"
	ssoservice "github.com/DimTur/lp_api_gateway/internal/services/sso"
	"github.com/go-playground/validator/v10"
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
	ssoService ssoservice.SsoService,
	lpservice lpservice.LpService,
	logger *slog.Logger,
	validator *validator.Validate,
	traceProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) (*App, error) {
	routerConfigurator := handlers.NewChiRouterConfigurator(
		ssoService,
		lpservice,
		logger,
		validator,
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
		validator,
	)
	if err != nil {
		logger.Error("failed to create server", slog.Any("err", err))
		return nil, err
	}

	return &App{
		HTTPServer: httpServer,
	}, nil
}
