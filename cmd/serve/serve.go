package serve

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DimTur/lp_api_gateway/internal/app"
	ssogrpc "github.com/DimTur/lp_api_gateway/internal/clients/sso/grpc"
	"github.com/DimTur/lp_api_gateway/internal/config"
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

			cfg, err := config.Parse(configPath)
			if err != nil {
				return err
			}

			ssoClient, err := ssogrpc.New(
				ctx,
				log,
				cfg.Clients.SSO.Address,
				cfg.Clients.SSO.Timeout,
				cfg.Clients.SSO.RetriesCount,
			)
			if err != nil {
				return err
			}

			application, err := app.NewApp(
				cfg.HTTPServer.Address,
				cfg.HTTPServer.Timeout,
				cfg.HTTPServer.Timeout,
				cfg.HTTPServer.IddleTimeout,
				*ssoClient,
				log,
			)
			if err != nil {
				return err
			}

			httCloser, err := application.HTTPServer.Run()
			if err != nil {
				return err
			}

			log.Info("server listening:", slog.Any("port", cfg.HTTPServer.Address))
			<-ctx.Done()

			httCloser()

			return nil
		},
	}
	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}
