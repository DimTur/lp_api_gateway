package app

import (
	"log/slog"

	httpapp "github.com/DimTur/lp_api_gateway/internal/app/http"
)

type App struct {
	HTTPServer *httpapp.APIServer
}

func NewApp(
	httpAddr string,
	logger *slog.Logger,
)
