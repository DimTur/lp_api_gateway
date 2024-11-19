package serve

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DimTur/lp_api_gateway/internal/app"
	lpgrpc "github.com/DimTur/lp_api_gateway/internal/clients/lp/grpc"
	ssogrpc "github.com/DimTur/lp_api_gateway/internal/clients/sso/grpc"
	"github.com/DimTur/lp_api_gateway/internal/config"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/validation"
	lpservice "github.com/DimTur/lp_api_gateway/internal/services/lp"
	"github.com/DimTur/lp_api_gateway/internal/services/permissions"
	ssoservice "github.com/DimTur/lp_api_gateway/internal/services/sso"
	"github.com/DimTur/lp_api_gateway/internal/services/storage/redis"
	"github.com/DimTur/lp_api_gateway/pkg/meter"
	"github.com/DimTur/lp_api_gateway/pkg/tracer"
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

			lpClient, err := lpgrpc.New(
				ctx,
				log,
				cfg.Clients.LP.Address,
				cfg.Clients.LP.Timeout,
				cfg.Clients.LP.RetriesCount,
			)
			if err != nil {
				return err
			}

			traceService, err := tracer.InitTracer(cfg.Tracer.OpenTelemetry.Address, cfg.Tracer.OpenTelemetry.ServiceName)
			if err != nil {
				return err
			}

			meterService, err := meter.InitMeter(ctx, cfg.Meter.Prometheus.Address)
			if err != nil {
				return err
			}

			redisPermissions := &redis.RedisPermissions{
				Host:     cfg.Redis.Host,
				Port:     cfg.Redis.Port,
				DB:       cfg.Redis.PermissionsDB,
				Password: cfg.Redis.Password,
			}
			redisPerm, err := redis.NewRedisClient(*redisPermissions)
			if err != nil {
				log.Error("failed to close redis", slog.Any("err", err))
			}

			validate := validation.InitValidator()

			permService := permissions.New(log, validate, lpClient, lpClient, ssoClient, redisPerm)
			ssoService := ssoservice.New(log, validate, ssoClient, ssoClient)
			lpService := lpservice.New(log, validate, lpClient, lpClient, lpClient, ssoClient, *permService)

			application, err := app.NewApp(
				cfg.HTTPServer.Address,
				cfg.HTTPServer.Timeout,
				cfg.HTTPServer.Timeout,
				cfg.HTTPServer.IddleTimeout,
				*ssoService,
				*lpService,
				log,
				validate,
				traceService,
				meterService,
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
