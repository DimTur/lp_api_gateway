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
	"google.golang.org/grpc/status"
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
	const op = "grpc.New"

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
	const op = "grpc.RegisterUser"

	resp, err := c.api.RegisterUser(ctx, &ssov1.RegisterUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return &ssov1.RegisterUserResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (c *Client) LoginUser(ctx context.Context, email string, password string) (*ssov1.LoginUserResponse, error) {
	const op = "grpc.LoginUser"

	resp, err := c.api.LoginUser(ctx, &ssov1.LoginUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("%s: authentication failed: %s", op, st.Message())
			case codes.InvalidArgument:
				return nil, fmt.Errorf("%s: invalid input: %s", op, st.Message())
			case codes.Internal:
				return nil, fmt.Errorf("%s: internal server error: %s", op, st.Message())
			default:
				return nil, fmt.Errorf("%s: unexpected error: %s", op, st.Message())
			}
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (c *Client) AuthCheck(ctx context.Context, accessToken string) (*ssov1.AuthCheckResponse, error) {
	const op = "grpc.AuthCheck"

	resp, err := c.api.AuthCheck(ctx, &ssov1.AuthCheckRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return nil, fmt.Errorf("%s: unauth: %s", op, st.Message())
		}
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
