package authhandler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/validation"
	"github.com/DimTur/lp_api_gateway/pkg/meter"
	"github.com/DimTur/lp_api_gateway/pkg/tracer"
	ssov1 "github.com/DimTur/lp_protos/gen/go/sso"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidAppID        = errors.New("invalid app id")
	ErrUserExists          = errors.New("user already exists")
	ErrAppExists           = errors.New("app already exists")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type SingUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type SingUpResponse struct {
	response.Response
	UserID int64 `json:"user_id,omitempty"`
}

type SingInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type SingInResponse struct {
	response.Response
	AccsessToken string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type AuthService interface {
	RegisterUser(ctx context.Context, email string, password string) (*ssov1.RegisterUserResponse, error)
	LoginUser(ctx context.Context, email string, password string) (*ssov1.LoginUserResponse, error)
}

var Validate = validator.New()

func init() {
	Validate.RegisterValidation("password", validation.PasswordValidator)
}

// SingUp godoc
// @Summary      Register a new user
// @Description  This endpoint allows users to register with an email and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        SingUpRequest body authhandler.SingUpRequest true "Registration parameters"
// @Success      201 {object} authhandler.SingUpResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      500 {object} response.Response "Server error"
// @Router       /sing_up [post]
func SingUp(log *slog.Logger, authService AuthService) http.HandlerFunc {
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

		var req SingUpRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request from", req.Email))

		span.AddEvent("validation_started")
		if err := Validate.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		span.AddEvent("validation_completed")
		span.SetAttributes(attribute.String("email", req.Email))

		span.AddEvent("started_user_registering")
		respID, err := authService.RegisterUser(r.Context(), req.Email, req.Password)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == codes.AlreadyExists {
					log.Error("user already exists", slog.Any("email", req.Email))
					w.WriteHeader(http.StatusBadRequest)
					render.JSON(w, r, response.Error("user already exists"))
					return
				}
			}

			log.Error("failed to add user", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to add user"))
			return
		}
		span.AddEvent("completed_user_registering")
		span.SetAttributes(attribute.String("email", req.Email))

		log.Info("user registered", slog.Int64("id", respID.UserId))

		render.JSON(w, r, SingUpResponse{
			Response: response.OK(),
			UserID:   respID.UserId,
		})
	}
}

// SignIn godoc
// @Summary      User Login
// @Description  This endpoint allows users to sign in using their email and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        SingInRequest body authhandler.SingInRequest true "Sign-in parameters"
// @Success      200 {object} authhandler.SingInResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      404 {object} response.Response "User not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /sing_in [post]
func SignIn(log *slog.Logger, authService AuthService) http.HandlerFunc {
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

		var req SingInRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request from", req.Email))

		span.AddEvent("validation_started")
		if err := Validate.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		span.AddEvent("validation_completed")
		span.SetAttributes(attribute.String("email", req.Email))

		span.AddEvent("started_user_login")
		singInResponse, err := authService.LoginUser(r.Context(), req.Email, req.Password)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				switch st.Code() {
				case codes.Unauthenticated:
					log.Info("invalid email or password", slog.String("email", req.Email))
					w.WriteHeader(http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid email or password"))
					return
				case codes.InvalidArgument:
					log.Error("invalid input", slog.String("err", st.Message()))
					w.WriteHeader(http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid input"))
					return
				case codes.NotFound:
					log.Error("user not found", slog.String("err", st.Message()))
					w.WriteHeader(http.StatusNotFound)
					render.JSON(w, r, response.Error("user not found"))
					return
				case codes.Internal:
					log.Error("internal server error", slog.String("err", st.Message()))
					w.WriteHeader(http.StatusInternalServerError)
					render.JSON(w, r, response.Error("internal server error"))
					return
				default:
					log.Error("unexpected error", slog.String("err", st.Message()))
					w.WriteHeader(http.StatusBadRequest)
					render.JSON(w, r, response.Error("unexpected error"))
					return
				}
			}

			log.Error("failed to login user", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to login user"))
			return
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
