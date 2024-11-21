package handlers

import (
	"log/slog"
	"net/http"
	"time"

	channelshandler "github.com/DimTur/lp_api_gateway/internal/handlers/learning_platform/channels"
	lessonshandler "github.com/DimTur/lp_api_gateway/internal/handlers/learning_platform/lessons"
	pageshandler "github.com/DimTur/lp_api_gateway/internal/handlers/learning_platform/pages"
	planshandler "github.com/DimTur/lp_api_gateway/internal/handlers/learning_platform/plans"
	authmiddleware "github.com/DimTur/lp_api_gateway/internal/handlers/middleware/auth"
	authhandler "github.com/DimTur/lp_api_gateway/internal/handlers/sso/auth"
	learninggrouphandler "github.com/DimTur/lp_api_gateway/internal/handlers/sso/learning_group"
	lpservice "github.com/DimTur/lp_api_gateway/internal/services/lp"
	ssoservice "github.com/DimTur/lp_api_gateway/internal/services/sso"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type RouterConfigurator interface {
	ConfigureRouter() http.Handler
}

type ChiRouterConfigurator struct {
	SsoService     ssoservice.SsoService
	LpService      lpservice.LpService
	Logger         *slog.Logger
	validator      *validator.Validate
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
}

func NewChiRouterConfigurator(
	ssoService ssoservice.SsoService,
	lpService lpservice.LpService,
	logger *slog.Logger,
	validator *validator.Validate,
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) *ChiRouterConfigurator {
	return &ChiRouterConfigurator{
		SsoService:     ssoService,
		LpService:      lpService,
		Logger:         logger,
		validator:      validator,
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
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
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-User-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Routes
	//
	// Server health cheker
	router.Get("/health", HealthCheckHandler)

	// Swagger
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// Trace and metrics
	router.Handle("/metrics", promhttp.Handler())

	// Auth
	router.Post("/sing_up", authhandler.SingUp(c.Logger, c.validator, &c.SsoService))
	router.Post("/sing_in", authhandler.SignIn(c.Logger, c.validator, &c.SsoService))
	router.Post("/sing_in_by_tg", authhandler.SignInByTelegram(c.Logger, c.validator, &c.SsoService))
	router.Post("/check_otp", authhandler.CheckOTPAndLogIn(c.Logger, c.validator, &c.SsoService))
	router.Group(func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware(c.Logger, c.validator, &c.SsoService))
		r.Patch("/profile/update_info", authhandler.UpdateUserInfo(c.Logger, c.validator, &c.SsoService))
	})

	// Lerning Groups
	router.Group(func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware(c.Logger, c.validator, &c.SsoService))
		r.Post("/learning_groups", learninggrouphandler.CreateLearningGroup(c.Logger, c.validator, &c.SsoService))
		r.Get("/learning_group/{id}", learninggrouphandler.GetLearningGroupByID(c.Logger, c.validator, &c.SsoService))
		r.Patch("/learning_group/{id}", learninggrouphandler.UpdateLearningGroup(c.Logger, c.validator, &c.SsoService))
		r.Delete("/learning_group/{id}", learninggrouphandler.DeleteLearningGroup(c.Logger, c.validator, &c.SsoService))
		r.Get("/learning_groups", learninggrouphandler.GetLearningGroups(c.Logger, c.validator, &c.SsoService))
	})

	// Learning Platform
	router.Group(func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware(c.Logger, c.validator, &c.SsoService))

		// Channels
		r.Post("/channels", channelshandler.CreateChannel(c.Logger, c.validator, &c.LpService))
		r.Get("/channels/{id}", channelshandler.GetChannel(c.Logger, c.validator, &c.LpService))
		r.Get("/channels", channelshandler.GetChannels(c.Logger, c.validator, &c.LpService))
		r.Patch("/channels/{id}", channelshandler.UpdateChannel(c.Logger, c.validator, &c.LpService))
		r.Delete("/channels/{id}", channelshandler.DeleteChannel(c.Logger, c.validator, &c.LpService))
		r.Post("/channels/{id}/share", channelshandler.ShareChannel(c.Logger, c.validator, &c.LpService))

		// Plans
		r.Post("/channels/{id}/plans", planshandler.CreatePlan(c.Logger, c.validator, &c.LpService))
		r.Get("/channels/{channel_id}/plans/{plan_id}", planshandler.GetPlan(c.Logger, c.validator, &c.LpService))
		r.Get("/channels/{id}/plans", planshandler.GetPlans(c.Logger, c.validator, &c.LpService))
		r.Patch("/channels/{channel_id}/plans/{plan_id}", planshandler.UpdatePlan(c.Logger, c.validator, &c.LpService))
		r.Delete("/channels/{channel_id}/plans/{plan_id}", planshandler.DeletePlan(c.Logger, c.validator, &c.LpService))
		r.Post("/channels/{channel_id}/plans/{plan_id}/share", planshandler.SharePlan(c.Logger, c.validator, &c.LpService))

		// Lessons
		r.Post("/channels/{channel_id}/plans/{plan_id}/lessons", lessonshandler.CreateLesson(c.Logger, c.validator, &c.LpService))
		r.Get("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}", lessonshandler.GetLesson(c.Logger, c.validator, &c.LpService))
		r.Get("/channels/{channel_id}/plans/{plan_id}/lessons", lessonshandler.GetLessons(c.Logger, c.validator, &c.LpService))
		r.Patch("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}", lessonshandler.UpdateLesson(c.Logger, c.validator, &c.LpService))
		r.Delete("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}", lessonshandler.DeleteLesson(c.Logger, c.validator, &c.LpService))

		// Pages
		r.Post("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/image_page", pageshandler.CreateImagePage(c.Logger, c.validator, &c.LpService))
		r.Post("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/video_page", pageshandler.CreateVideoPage(c.Logger, c.validator, &c.LpService))
		r.Post("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/pdf_page", pageshandler.CreateVideoPage(c.Logger, c.validator, &c.LpService))
		r.Get("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/image_page/{page_id}}", pageshandler.GetImagePage(c.Logger, c.validator, &c.LpService))
		r.Get("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/video_page/{page_id}}", pageshandler.GetVideoPage(c.Logger, c.validator, &c.LpService))
		r.Get("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/pdf_page/{page_id}}", pageshandler.GetPDFPage(c.Logger, c.validator, &c.LpService))
		r.Get("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/pages", pageshandler.GetPages(c.Logger, c.validator, &c.LpService))
		r.Patch("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/image_page/{page_id}", pageshandler.UpdateImagePage(c.Logger, c.validator, &c.LpService))
		r.Patch("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/video_page/{page_id}", pageshandler.UpdateVideoPage(c.Logger, c.validator, &c.LpService))
		r.Patch("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/pdf_page/{page_id}", pageshandler.UpdatePDFPage(c.Logger, c.validator, &c.LpService))
		r.Delete("/channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/pages/{page_id}", pageshandler.DeletePage(c.Logger, c.validator, &c.LpService))
	})

	return router
}
