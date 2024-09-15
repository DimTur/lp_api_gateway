package handlers

import (
	"log/slog"
	"net/http"
	"time"

	ssogrpc "github.com/DimTur/lp_api_gateway/internal/clients/sso/grpc"
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
	Logger         *slog.Logger
}

func NewChiRouterConfigurator(grpcClient ssogrpc.Client, logger *slog.Logger) *ChiRouterConfigurator {
	return &ChiRouterConfigurator{
		AuthGRPCClient: grpcClient,
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
		r.Get("/protected", ProtectedHandler)
	})

	return router
}

// Test handler
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected route"))
}
