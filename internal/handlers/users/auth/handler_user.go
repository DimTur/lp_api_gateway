package authhandler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	ssogrpc "github.com/DimTur/lp_api_gateway/internal/clients/sso/grpc"
	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
	"github.com/DimTur/lp_api_gateway/pkg/meter"
	"github.com/DimTur/lp_api_gateway/pkg/tracer"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidAppID        = errors.New("invalid app id")
	ErrUserExists          = errors.New("user already exists")
	ErrAppExists           = errors.New("app already exists")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type AuthService interface {
	RegisterUser(ctx context.Context, newUser *ssomodels.RegisterUser) (*ssomodels.RegisterResp, error)
	LoginUser(ctx context.Context, logUser *ssomodels.LogIn) (*ssomodels.LogInResp, error)
}

// SingUp godoc
// @Summary      Register a new user
// @Description  This endpoint allows users to register with an email and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        ssomodels.RegisterUser body ssomodels.RegisterUser true "Registration parameters"
// @Success      201 {object} authhandler.SingUpResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      500 {object} response.Response "Server error"
// @Router       /sing_up [post]
func SingUp(log *slog.Logger, val *validator.Validate, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.auth.SingUp"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignUpReqCount.Add(r.Context(), 1)

		_, span := tracer.AuthTracer.Start(r.Context(), "SingUp")
		defer span.End()

		var req ssomodels.RegisterUser
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request from", req.Email))

		span.AddEvent("validation_started")
		if err := val.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		span.AddEvent("validation_completed")
		span.SetAttributes(attribute.String("email", req.Email))

		span.AddEvent("started_user_registering")
		resp, err := authService.RegisterUser(r.Context(), &req)
		if err != nil {
			switch {
			case errors.Is(err, ssogrpc.ErrUserExists):
				log.Error("user already exists", slog.Any("email", req.Email))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("user already exists"))
				return
			case errors.Is(err, ssogrpc.ErrInvalidCredentials):
				log.Error("invalid credentinals", slog.Any("email", req.Email))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			default:
				log.Error("registratin failed", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("registratin failed"))
				return
			}
		}
		span.AddEvent("completed_user_registering")
		span.SetAttributes(attribute.String("email", req.Email))

		log.Info("user registered")

		render.JSON(w, r, SingUpResponse{
			Response: response.OK(),
			Success:  resp.Success,
		})
	}
}

// SignIn godoc
// @Summary      User Login
// @Description  This endpoint allows users to sign in using their email and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        ssomodels.LogIn body ssomodels.LogIn true "Sign-in parameters"
// @Success      200 {object} authhandler.SingInResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      404 {object} response.Response "User not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /sing_in [post]
func SignIn(log *slog.Logger, val *validator.Validate, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.auth.SignIn"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignInReqCount.Add(r.Context(), 1)

		_, span := tracer.AuthTracer.Start(r.Context(), "SignIn")
		defer span.End()

		var req ssomodels.LogIn
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request from", req.Email))

		span.AddEvent("validation_started")
		if err := val.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		span.AddEvent("validation_completed")
		span.SetAttributes(attribute.String("email", req.Email))

		span.AddEvent("started_user_login")
		singInResponse, err := authService.LoginUser(r.Context(), &req)
		if err != nil {
			switch {
			case errors.Is(err, ssogrpc.ErrInvalidCredentials):
				log.Error("invalid input", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid input"))
				return
			default:
				log.Error("failed to login user", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to login user"))
				return
			}
		}
		span.AddEvent("completed_user_login")
		span.SetAttributes(attribute.String("email", req.Email))

		log.Info("user logged in successfully")

		render.JSON(w, r, SingInResponse{
			Response:     response.OK(),
			AccsessToken: singInResponse.AccessToken,
			RefreshToken: singInResponse.RefreshToken,
		})
	}
}
