package ssogrpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ssov1 "github.com/DimTur/lp_protos/gen/go/sso"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api ssov1.AuthClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "sso.grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	// TODO: secure conn
	cc, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: ssov1.NewAuthClient(cc),
		log: log,
	}, nil
}

func (c *Client) RegisterUser(ctx context.Context, email string, password string) (*ssov1.RegisterUserResponse, error) {
	const op = "sso.grpc.RegisterUser"

	resp, err := c.api.RegisterUser(ctx, &ssov1.RegisterUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		c.log.Error("received error from auth grpc service", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (c *Client) LoginUser(ctx context.Context, email string, password string) (*ssov1.LoginUserResponse, error) {
	const op = "sso.grpc.LoginUser"

	resp, err := c.api.LoginUser(ctx, &ssov1.LoginUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		c.log.Error("received error from auth grpc service", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (c *Client) AuthCheck(ctx context.Context, accessToken string) (*ssov1.AuthCheckResponse, error) {
	const op = "sso.grpc.AuthCheck"

	resp, err := c.api.AuthCheck(ctx, &ssov1.AuthCheckRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		c.log.Error("received error from auth grpc service", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

// InterceptorLogger adapts slog logger to iterceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
