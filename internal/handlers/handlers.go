package handlers

import (
	"log/slog"
	"net/http"

	"github.com/DimTur/lp_api_gateway/internal/clients/sso/grpc"
	authhandler "github.com/DimTur/lp_api_gateway/internal/handlers/users/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type RouterConfigurator interface {
	ConfigureRouter() http.Handler
}

type ChiRouterConfigurator struct {
	GRPCClient grpc.Client
	Logger     *slog.Logger
}

func NewChiRouterConfigurator(grpcClient grpc.Client, logger *slog.Logger) *ChiRouterConfigurator {
	return &ChiRouterConfigurator{
		GRPCClient: grpcClient,
		Logger:     logger,
	}
}

func (c *ChiRouterConfigurator) ConfigureRouter() http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.URLFormat)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Routes
	//
	// Server health cheker
	router.Get("/health", HealthCheckHandler)

	// Auth
	router.Post("/sing_up", authhandler.SingUp(c.Logger, &c.GRPCClient))

	return router
}
