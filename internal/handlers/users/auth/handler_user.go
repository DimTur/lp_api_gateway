package authhandler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/DimTur/lp_api_gateway/internal/lib/api/response"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/validation"
	ssov1 "github.com/DimTur/lp_protos/gen/go/sso"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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
	resp.Response
	UserID int64 `json:"user_id,omitempty"`
}

type SingInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type SingInResponse struct {
	resp.Response
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

func SingUp(log *slog.Logger, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.auth.SingUp"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SingUpRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := Validate.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("err", err.Error()))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		respID, err := authService.RegisterUser(r.Context(), req.Email, req.Password)
		if errors.Is(err, ErrUserExists) {
			log.Info("user already exists", slog.String("user", req.Email))
			render.JSON(w, r, resp.Error("user already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add user", slog.String("err", err.Error()))
			render.JSON(w, r, resp.Error("failed to add user"))
			return
		}

		log.Info("user registered", slog.Int64("id", respID.UserId))

		render.JSON(w, r, SingUpResponse{
			Response: resp.OK(),
			UserID:   respID.UserId,
		})
	}
}

func SignIn(log *slog.Logger, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.auth.SignIn"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SingInRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := Validate.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("err", err.Error()))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		singInResponse, err := authService.LoginUser(r.Context(), req.Email, req.Password)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				switch st.Code() {
				case codes.Unauthenticated:
					log.Info("invalid email or password", slog.String("email", req.Email))
					render.JSON(w, r, resp.Error("invalid email or password"))
					return
				case codes.InvalidArgument:
					log.Error("invalid input", slog.String("err", st.Message()))
					render.JSON(w, r, resp.Error("invalid input"))
					return
				case codes.Internal:
					log.Error("internal server error", slog.String("err", st.Message()))
					render.JSON(w, r, resp.Error("internal server error"))
					return
				default:
					log.Error("unexpected error", slog.String("err", st.Message()))
					render.JSON(w, r, resp.Error("unexpected error"))
					return
				}
			}

			log.Error("failed to login user", slog.String("err", err.Error()))
			render.JSON(w, r, resp.Error("failed to login user"))
			return
		}

		log.Info("user logged in successfully")

		render.JSON(w, r, SingInResponse{
			Response:     resp.OK(),
			AccsessToken: singInResponse.AccessToken,
			RefreshToken: singInResponse.RefreshToken,
		})
	}
}
