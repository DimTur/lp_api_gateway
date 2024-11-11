package authhandler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"

	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
	ssoservice "github.com/DimTur/lp_api_gateway/internal/services/sso"
	"github.com/DimTur/lp_api_gateway/pkg/meter"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type AuthService interface {
	RegisterUser(ctx context.Context, newUser *ssomodels.RegisterUser) (*ssomodels.RegisterResp, error)
	LoginUser(ctx context.Context, logUser *ssomodels.LogIn) (*ssomodels.LogInResp, error)
	LogInViaTg(ctx context.Context, email *ssomodels.LogInViaTg) (*ssomodels.LogInViaTgResp, error)
	CheckOTPAndLogIn(ctx context.Context, otp *ssomodels.CheckOTPAndLogIn) (*ssomodels.CheckOTPAndLogInResp, error)
	UpdateUserInfo(ctx context.Context, newInfo *ssomodels.UpdateUserInfo) (*ssomodels.UpdateUserInfoResp, error)
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
		const op = "handlers.sso.auth.SingUp"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignUpReqCount.Add(r.Context(), 1)

		var req ssomodels.RegisterUser
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request from", req.Email))

		resp, err := authService.RegisterUser(r.Context(), &req)
		if err != nil {
			switch {
			case errors.Is(err, ssoservice.ErrUserExists):
				log.Error("user already exists", slog.Any("email", req.Email))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("user already exists"))
				return
			case errors.Is(err, ssoservice.ErrInvalidCredentials):
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
		const op = "handlers.sso.auth.SignIn"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignInReqCount.Add(r.Context(), 1)

		var req ssomodels.LogIn
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request from", req.Email))

		singInResponse, err := authService.LoginUser(r.Context(), &req)
		if err != nil {
			switch {
			case errors.Is(err, ssoservice.ErrInvalidCredentials):
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

		log.Info("user logged in successfully")

		render.JSON(w, r, SingInResponse{
			Response:     response.OK(),
			AccsessToken: singInResponse.AccessToken,
			RefreshToken: singInResponse.RefreshToken,
		})
	}
}

// SignInByTelegram godoc
// @Summary      User Login by telegram bot
// @Description  This endpoint allows users to sign in using their email and sends OTP code to chat.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        ssomodels.LogInViaTg body ssomodels.LogInViaTg true "Sign-in parameters"
// @Success      200 {object} authhandler.SingInByTgResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      404 {object} response.Response "User not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /sing_in_by_tg [post]
func SignInByTelegram(log *slog.Logger, val *validator.Validate, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.sso.auth.SignInByTelegram"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignInReqCount.Add(r.Context(), 1)

		var req ssomodels.LogInViaTg
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request from", req.Email))

		resp, err := authService.LogInViaTg(r.Context(), &req)
		if err != nil {
			switch {
			case errors.Is(err, ssoservice.ErrInvalidCredentials):
				log.Error("invalid input", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid input"))
				return
			case errors.Is(err, ssoservice.ErrUserNotFound):
				log.Error("user not found", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("user not found"))
				return
			default:
				log.Error("failed to login user by telegram", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to login user by telegram"))
				return
			}
		}

		log.Info("user logged in successfully")

		render.JSON(w, r, SingInByTgResponse{
			Response: response.OK(),
			Success:  resp.Success,
			Info:     resp.Info,
		})
	}
}

// CheckOTPAndLogIn godoc
// @Summary      User Login by telegram bot
// @Description  This endpoint allows users to sign in using their email and sends OTP code to chat.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        ssomodels.CheckOTPAndLogIn body ssomodels.CheckOTPAndLogIn true "Sign-in parameters"
// @Success      200 {object} authhandler.CheckOTPAndLogInResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      404 {object} response.Response "User not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /check_otp [post]
func CheckOTPAndLogIn(log *slog.Logger, val *validator.Validate, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.sso.auth.CheckOTPAndLogIn"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignInReqCount.Add(r.Context(), 1)

		var req ssomodels.CheckOTPAndLogIn
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request from", req.Email))

		resp, err := authService.CheckOTPAndLogIn(r.Context(), &req)
		if err != nil {
			switch {
			case errors.Is(err, ssoservice.ErrInvalidCredentials):
				log.Error("invalid input", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid input"))
				return
			case errors.Is(err, ssoservice.ErrUserNotFound):
				log.Error("user not found", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("user not found"))
				return
			default:
				log.Error("failed to check otp and login", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to check otp and login"))
				return
			}
		}

		log.Info("user logged in successfully")

		render.JSON(w, r, CheckOTPAndLogInResponse{
			Response:     response.OK(),
			AccsessToken: resp.AccessToken,
			RefreshToken: resp.RefreshToken,
		})
	}
}

// UpdateUserInfo godoc
// @Summary      Change self user info
// @Description  This endpoint allow users to change their profile.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        authhandler.UpdateUserInfoReq body authhandler.UpdateUserInfoReq true "Sign-in parameters"
// @Success      200 {object} authhandler.UpdateUserInfoResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      500 {object} response.Response "Server error"
// @Router       /profile/update_info [patch]
// @Security ApiKeyAuth
func UpdateUserInfo(log *slog.Logger, val *validator.Validate, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.sso.auth.UpdateUserInfo"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignInReqCount.Add(r.Context(), 1)

		var req UpdateUserInfoReq
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request from", req.Email))

		resp, err := authService.UpdateUserInfo(r.Context(), &ssomodels.UpdateUserInfo{
			ID:     r.Header.Get("X-User-ID"),
			Email:  req.Email,
			Name:   req.Name,
			TgLink: req.TgLink,
		})
		if err != nil {
			switch {
			case errors.Is(err, ssoservice.ErrInvalidCredentials):
				log.Error("invalid input", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid input"))
				return
			default:
				log.Error("failed to update user info", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to update user info"))
				return
			}
		}

		log.Info("user info updated in successfully")

		render.JSON(w, r, UpdateUserInfoResponse{
			Response: response.OK(),
			Success:  resp.Success,
		})
	}
}
