package handlers

import (
	"log/slog"
	"net/http"
	"time"

	lpgrpc "github.com/DimTur/lp_api_gateway/internal/clients/lp/grpc"
	ssogrpc "github.com/DimTur/lp_api_gateway/internal/clients/sso/grpc"
	channelshandler "github.com/DimTur/lp_api_gateway/internal/handlers/learning_platform/channels"
	authmiddleware "github.com/DimTur/lp_api_gateway/internal/handlers/middleware/auth"
	authhandler "github.com/DimTur/lp_api_gateway/internal/handlers/users/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

type RouterConfigurator interface {
	ConfigureRouter() http.Handler
}

type ChiRouterConfigurator struct {
	AuthGRPCClient ssogrpc.Client
	LPGRPCClient   lpgrpc.Client
	Logger         *slog.Logger
}

func NewChiRouterConfigurator(authGRPCClient ssogrpc.Client, lpGRPCClient lpgrpc.Client, logger *slog.Logger) *ChiRouterConfigurator {
	return &ChiRouterConfigurator{
		AuthGRPCClient: authGRPCClient,
		LPGRPCClient:   lpGRPCClient,
		Logger:         logger,
	}
}

func (c *ChiRouterConfigurator) ConfigureRouter() http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.URLFormat)
	router.Use(httprate.LimitByIP(100, 1*time.Minute))
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
	router.Post("/sing_up", authhandler.SingUp(c.Logger, &c.AuthGRPCClient))
	router.Post("/sing_in", authhandler.SignIn(c.Logger, &c.AuthGRPCClient))

	// Learning Platform
	router.Group(func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware(&c.AuthGRPCClient))
		r.Post("/create_channel", channelshandler.CreateChannel(c.Logger, &c.LPGRPCClient))
		r.Get("/get_channel/{id}", channelshandler.GetChannel(c.Logger, &c.LPGRPCClient))
	})

	return router
}
