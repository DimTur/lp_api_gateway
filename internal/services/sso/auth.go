package ssoservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	ssogrpc "github.com/DimTur/lp_api_gateway/internal/clients/sso/grpc"
	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/DimTur/lp_api_gateway/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserExists          = errors.New("user already exists")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

func (sso *SsoService) RegisterUser(ctx context.Context, newUser *ssomodels.RegisterUser) (*ssomodels.RegisterResp, error) {
	const op = "internal.services.sso.auth.LoginUser"

	log := sso.Log.With(
		slog.String("op", op),
		slog.String("user_email", newUser.Email),
	)

	_, span := tracer.AuthTracer.Start(ctx, "RegisterUser")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := sso.Validator.Struct(newUser); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("email", newUser.Email))

	log.Info("registering user")

	// Start registration
	span.AddEvent("started_user_registering")
	reg, err := sso.AuthProvider.RegisterUser(ctx, newUser)
	if err != nil {
		switch {
		case errors.Is(err, ssogrpc.ErrUserExists):
			log.Error("user already exists", slog.Any("email", newUser.Email))
			return nil, fmt.Errorf("%s: %w", op, ErrUserExists)
		case errors.Is(err, ssogrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("email", newUser.Email))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("registratin failed", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_user_registering")
	span.SetAttributes(attribute.String("email", newUser.Email))

	log.Info("user registered")

	return &ssomodels.RegisterResp{
		Success: reg.Success,
	}, nil
}

func (sso *SsoService) LoginUser(ctx context.Context, logUser *ssomodels.LogIn) (*ssomodels.LogInResp, error) {
	const op = "internal.services.sso.auth.LoginUser"

	log := sso.Log.With(
		slog.String("op", op),
		slog.String("user_email", logUser.Email),
	)

	_, span := tracer.AuthTracer.Start(ctx, "LoginUser")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := sso.Validator.Struct(logUser); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("email", logUser.Email))

	log.Info("loging in started")

	// Start login
	span.AddEvent("started_user_login")
	logIn, err := sso.AuthProvider.LoginUser(ctx, logUser)
	if err != nil {
		switch {
		case errors.Is(err, ssogrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("email", logUser.Email))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to login user", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_user_login")
	span.SetAttributes(attribute.String("email", logUser.Email))

	log.Info("user logged in successfully")

	return &ssomodels.LogInResp{
		AccessToken:  logIn.AccessToken,
		RefreshToken: logIn.RefreshToken,
	}, nil
}

func (sso *SsoService) LogInViaTg(ctx context.Context, email *ssomodels.LogInViaTg) (*ssomodels.LogInViaTgResp, error) {
	const op = "internal.services.sso.auth.LogInViaTg"

	log := sso.Log.With(
		slog.String("op", op),
		slog.String("user_email", email.Email),
	)

	_, span := tracer.AuthTracer.Start(ctx, "LogInViaTg")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := sso.Validator.Struct(email); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("email", email.Email))

	log.Info("starting user login by telegram")

	// Start login
	span.AddEvent("started_user_login_by_telegram")
	resp, err := sso.AuthProvider.LogInViaTg(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, ssogrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("email", email.Email))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, ssogrpc.ErrUserNotFound):
			log.Error("user not found", slog.Any("email", email.Email))
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		default:
			log.Error("failed to login user by telegram", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_user_login_by_telegram")
	span.SetAttributes(attribute.String("email", email.Email))

	log.Info("otp code sent to telegram bot")

	return &ssomodels.LogInViaTgResp{
		Success: resp.Success,
		Info:    resp.Info,
	}, nil
}

func (sso *SsoService) CheckOTPAndLogIn(ctx context.Context, otp *ssomodels.CheckOTPAndLogIn) (*ssomodels.CheckOTPAndLogInResp, error) {
	const op = "internal.services.sso.auth.CheckOTPAndLogIn"

	log := sso.Log.With(
		slog.String("op", op),
		slog.String("user_email", otp.Email),
	)

	_, span := tracer.AuthTracer.Start(ctx, "CheckOTPAndLogIn")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := sso.Validator.Struct(otp); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("email", otp.Email))

	log.Info("starting checking otp code")

	// Start checking
	span.AddEvent("started_user_checking_otp_and_login")
	resp, err := sso.AuthProvider.CheckOTPAndLogIn(ctx, otp)
	if err != nil {
		switch {
		case errors.Is(err, ssogrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("email", otp.Email))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, ssogrpc.ErrUserNotFound):
			log.Error("user not found", slog.Any("email", otp.Email))
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		default:
			log.Error("failed to check otp and login", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_user_checking_otp_and_login")
	span.SetAttributes(attribute.String("email", otp.Email))

	log.Info("user logged in successfully")

	return &ssomodels.CheckOTPAndLogInResp{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (sso *SsoService) UpdateUserInfo(ctx context.Context, newInfo *ssomodels.UpdateUserInfo) (*ssomodels.UpdateUserInfoResp, error) {
	const op = "internal.services.sso.auth.UpdateUserInfo"

	log := sso.Log.With(
		slog.String("op", op),
		slog.String("user_id", newInfo.ID),
	)

	_, span := tracer.AuthTracer.Start(ctx, "UpdateUserInfo")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := sso.Validator.Struct(newInfo); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("user_id", newInfo.ID))

	log.Info("updating user info")

	// Start updating
	span.AddEvent("updating_user_info")
	resp, err := sso.AuthProvider.UpdateUserInfo(ctx, newInfo)
	if err != nil {
		switch {
		case errors.Is(err, ssogrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("user_id", newInfo.ID))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to update user info", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_update_user_info")
	span.SetAttributes(attribute.String("user_id", newInfo.ID))

	log.Info("user info updated successfully")

	return &ssomodels.UpdateUserInfoResp{
		Success: resp.Success,
	}, nil
}

func (sso *SsoService) AuthCheck(ctx context.Context, authCheck *ssomodels.AuthCheck) (*ssomodels.AuthCheckResp, error) {
	const op = "internal.services.sso.auth.AuthCheck"

	log := sso.Log.With(
		slog.String("op", op),
	)

	_, span := tracer.AuthTracer.Start(ctx, "AuthCheck")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := sso.Validator.Struct(authCheck); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	log.Info("auth checking")

	// Start auth checking
	span.AddEvent("started_auth_cheking")
	resp, err := sso.AuthProvider.AuthCheck(ctx, authCheck)
	if err != nil {
		switch {
		case errors.Is(err, ssogrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to auth check", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	if !resp.IsValid {
		log.Error("invalid credentinals", slog.Bool("is_valid", false))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("completed_auth_cheking")
	span.SetAttributes(attribute.String("userID", resp.UserID))

	return &ssomodels.AuthCheckResp{
		IsValid: resp.IsValid,
		UserID:  resp.UserID,
	}, nil
}
