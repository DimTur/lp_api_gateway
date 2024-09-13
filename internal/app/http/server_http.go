package httpapp

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

const (
	// HTTPDefaultGracefulStopTimeout - время ожидания перед завершением сервера
	HTTPDefaultGracefulStopTimeout = 5 * time.Second
)

type APIServer struct {
	httpAddr            string
	httpSrv             *http.Server
	gracefulStopTimeout time.Duration
	logger              *slog.Logger
}

// https://github.com/gopherschool/http-rest-api/blob/master/internal/app/apiserver/server.go
func NewHTTPServer(
	httpAddr string,
	handler http.Handler,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	iddleTimeout time.Duration,
	logger *slog.Logger,
) (*APIServer, error) {
	const op = "http-server"

	logger = logger.With(
		slog.String("op", op),
		slog.String("addr", httpAddr),
	)

	httpSrv := &http.Server{
		Addr:         httpAddr,
		Handler:      handler,
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

// GracefulShutdown - function for monitoring signals and server termination
// func (a *APIServer) GracefulShutdown(closeFunc func() error) {
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, os.Interrupt)
// 	<-quit

// 	a.logger.Info("received shutdown signal")
// 	if err := closeFunc(); err != nil {
// 		a.logger.Error("error during shutdown", slog.Any("err", err))
// 	}
// }
