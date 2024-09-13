package serve

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimTur/lp_api_gateway/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	var configPath string

	c := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Start API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			defer cancel()

			router := chi.NewRouter()
			router.Use(middleware.RequestID)
			router.Use(middleware.Recoverer)
			router.Use(middleware.Logger)

			cfg, err := config.Parse(configPath)
			if err != nil {
				return err
			}

			application, err := app.NewApp(cfg.HTTPServer.Address, log)
			if err != nil {
				return err
			}

			log.Info("server listening:", slog.Any("port", cfg.HTTPServer.Address))
			<-ctx.Done()

			closeCtx, _ := context.WithTimeout(context.Background(), time.Second)
			if err := httpServer.Shutdown(closeCtx); err != nil {
				log.Error("httpServer.Shutdown", slog.Any("err", err))
			}

			return nil
		},
	}
	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}
