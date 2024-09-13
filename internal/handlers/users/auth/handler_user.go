package authhandler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/DimTur/lp_api_gateway/internal/lib/api/response"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/validation"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidAppID        = errors.New("invalid app id")
	ErrUserExists          = errors.New("user already exists")
	ErrAppExists           = errors.New("app already exists")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type Response struct {
	resp.Response
	UserID int64 `json:"user_id,omitempty"`
}

type AuthService interface {
	RegisterUser(ctx context.Context, email string, password string) (userID int64, err error)
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

		var req Request

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

		id, err := authService.RegisterUser(r.Context(), req.Email, req.Password)
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

		log.Info("user registered", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			UserID:   id,
		})
	}
}
