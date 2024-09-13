package app

import (
	"log/slog"
	"time"

	httpapp "github.com/DimTur/lp_api_gateway/internal/app/http"
	"github.com/DimTur/lp_api_gateway/internal/handlers"
)

type App struct {
	HTTPServer *httpapp.APIServer
}

func NewApp(
	httpAddr string,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	iddleTimeout time.Duration,
	logger *slog.Logger,
) (*App, error) {
	routerConfigurator := &handlers.ChiRouterConfigurator{}
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
