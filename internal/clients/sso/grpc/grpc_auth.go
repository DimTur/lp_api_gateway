package ssogrpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	ssov1 "github.com/DimTur/lp_protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidAppID        = errors.New("invalid app id")
	ErrUserExists          = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrAppExists           = errors.New("app already exists")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrOtpNotFound         = errors.New("otp not found")

	ErrInternal = errors.New("internal error")
)

func (c *Client) RegisterUser(ctx context.Context, newUser *ssomodels.RegisterUser) (*ssomodels.RegisterResp, error) {
	const op = "sso.grpc_auth.RegisterUser"

	resp, err := c.api.RegisterUser(ctx, &ssov1.RegisterUserRequest{
		Email:    newUser.Email,
		Password: newUser.Password,
		Name:     newUser.Name,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.AlreadyExists:
			c.log.Error("user alredy exists", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrUserExists)
		case codes.InvalidArgument:
			c.log.Error("invalid arguments", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &ssomodels.RegisterResp{
		Success: resp.Success,
	}, nil
}

func (c *Client) LoginUser(ctx context.Context, logUser *ssomodels.LogIn) (*ssomodels.LogInResp, error) {
	const op = "sso.grpc_auth.LoginUser"

	resp, err := c.api.LoginUser(ctx, &ssov1.LoginUserRequest{
		Email:    logUser.Email,
		Password: logUser.Password,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			c.log.Error("invalid credentials", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case codes.NotFound:
			c.log.Error("user not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case codes.InvalidArgument:
			c.log.Error("invalid credentials", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &ssomodels.LogInResp{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (c *Client) LogInViaTg(ctx context.Context, email *ssomodels.LogInViaTg) (*ssomodels.LogInViaTgResp, error) {
	const op = "sso.grpc_auth.LogInViaTg"

	resp, err := c.api.LoginViaTg(ctx, &ssov1.LoginViaTgRequest{
		Email: email.Email,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("invalid email", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case codes.NotFound:
			c.log.Error("user not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &ssomodels.LogInViaTgResp{
		Success: resp.Success,
		Info:    resp.Info,
	}, nil
}

func (c *Client) CheckOTPAndLogIn(ctx context.Context, otp *ssomodels.CheckOTPAndLogIn) (*ssomodels.CheckOTPAndLogInResp, error) {
	const op = "sso.grpc_auth.CheckOTPAndLogIn"

	resp, err := c.api.CheckOTPAndLogIn(ctx, &ssov1.CheckOTPAndLogInRequest{
		Email: otp.Email,
		Code:  otp.Code,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("invalid credentials", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case codes.NotFound:
			c.log.Error("user not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &ssomodels.CheckOTPAndLogInResp{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (c *Client) UpdateUserInfo(ctx context.Context, newInfo *ssomodels.UpdateUserInfo) (*ssomodels.UpdateUserInfoResp, error) {
	const op = "sso.grpc_auth.UpdateUserInfo"

	resp, err := c.api.UpdateUserInfo(ctx, &ssov1.UpdateUserInfoRequest{
		Id:      newInfo.ID,
		Email:   newInfo.Email,
		Name:    newInfo.Name,
		TgLink:  newInfo.TgLink,
		IsAdmin: newInfo.IsAdmin,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("invalid credentials", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &ssomodels.UpdateUserInfoResp{
		Success: resp.Success,
	}, nil
}

func (c *Client) RefreshToken(ctx context.Context, refToken *ssomodels.RefreshToken) (*ssomodels.RefreshTokenResp, error) {
	const op = "sso.grpc_auth.RefreshToken"

	resp, err := c.api.RefreshToken(ctx, &ssov1.RefreshTokenRequest{
		RefreshToken: refToken.RefreshToken,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("invalid credentials", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &ssomodels.RefreshTokenResp{
		AccessToken: resp.AccessToken,
	}, nil
}

func (c *Client) IsAdmin(ctx context.Context, userID *ssomodels.IsAdmin) (*ssomodels.IsAdminResp, error) {
	const op = "sso.grpc_auth.IsAdmin"

	resp, err := c.api.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: userID.UserID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("user not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &ssomodels.IsAdminResp{
		IsAdmin: resp.IsAdmin,
	}, nil
}

func (c *Client) AuthCheck(ctx context.Context, authCheck *ssomodels.AuthCheck) (*ssomodels.AuthCheckResp, error) {
	const op = "sso.grpc_auth.AuthCheck"

	resp, err := c.api.AuthCheck(ctx, &ssov1.AuthCheckRequest{
		AccessToken: authCheck.AccessToken,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			c.log.Error("invalid credentials", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &ssomodels.AuthCheckResp{
		IsValid: resp.IsValid,
		UserID:  resp.UserId,
	}, nil
}
