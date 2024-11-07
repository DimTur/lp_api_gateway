package httpapp

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	HTTPDefaultGracefulStopTimeout = 5 * time.Second
)

type APIServer struct {
	httpAddr            string
	httpSrv             *http.Server
	gracefulStopTimeout time.Duration
	logger              *slog.Logger
	validator           *validator.Validate
}

func NewHTTPServer(
	httpAddr string,
	router http.Handler,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	iddleTimeout time.Duration,
	logger *slog.Logger,
	validator *validator.Validate,
) (*APIServer, error) {
	const op = "http-server"

	logger = logger.With(
		slog.String("op", op),
		slog.String("addr", httpAddr),
	)

	httpSrv := &http.Server{
		Addr:         httpAddr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  iddleTimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return &APIServer{
		httpAddr:            httpAddr,
		httpSrv:             httpSrv,
		gracefulStopTimeout: HTTPDefaultGracefulStopTimeout,
		logger:              logger,
		validator:           validator,
	}, nil
}

// Run starts HTTP server
func (a *APIServer) Run() (func() error, error) {
	const op = "httpapp.Run"
	a.logger.With(slog.String("op", op)).Info("starting", slog.String("httpAddr", a.httpAddr))

	go func() {
		if err := a.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("http server", slog.Any("err", err))
		}
	}()

	a.logger.Info("http server is running", slog.String("addr", a.httpAddr))
	return a.close, nil
}

// close stops HTTP server with graceful shutdown
func (a *APIServer) close() error {
	const op = "httpapp.close"
	a.logger.With(slog.String("op", op)).Info("stopping HTTP server", slog.String("addr", a.httpAddr))

	ctx, cancel := context.WithTimeout(context.Background(), a.gracefulStopTimeout)
	defer cancel()

	if err := a.httpSrv.Shutdown(ctx); err != nil {
		a.logger.Error("graceful shutdown failed, forcing stop", slog.Any("err", err))
		if err := a.httpSrv.Close(); err != nil {
			return err
		}
	}

	a.logger.With(slog.String("op", op)).Info("stopped", slog.String("addr", a.httpAddr))
	return nil
}
