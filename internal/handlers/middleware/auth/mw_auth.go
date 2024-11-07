package authmiddleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	ssogrpc "github.com/DimTur/lp_api_gateway/internal/clients/sso/grpc"
	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type AuthService interface {
	AuthCheck(ctx context.Context, authChek *ssomodels.AuthCheck) (*ssomodels.AuthCheckResp, error)
}

func AuthMiddleware(log *slog.Logger, val *validator.Validate, authService AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.Auth"

			log = log.With(
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(r.Context())),
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
			)

			tracer := otel.Tracer("AuthTracer")
			_, span := tracer.Start(r.Context(), "AuthMiddleware")
			defer span.End()

			accessToken := r.Header.Get("Authorization")
			if accessToken == "" {
				log.Info("authorization token not provided")
				w.WriteHeader(http.StatusUnauthorized)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			span.AddEvent("started_auth_cheking")
			authCheck := &ssomodels.AuthCheck{
				AccessToken: accessToken,
			}
			resp, err := authService.AuthCheck(r.Context(), authCheck)
			if err != nil {
				switch {
				case errors.Is(err, ssogrpc.ErrInvalidCredentials):
					log.Error("error checking authorization", slog.String("err", err.Error()))
					w.WriteHeader(http.StatusUnauthorized)
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				default:
					log.Error("internal error", slog.String("err", err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					http.Error(w, "internal error", http.StatusInternalServerError)
					return
				}
			}

			if !resp.IsValid {
				log.Info("invalid authorization token", slog.String("user_id", resp.UserID))
				w.WriteHeader(http.StatusUnauthorized)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			span.AddEvent("completed_auth_cheking")
			span.SetAttributes(attribute.String("userID", resp.UserID))

			r.Header.Set("X-User-ID", resp.UserID)

			log.Info("authorization successful", slog.String("user_id", resp.UserID))

			next.ServeHTTP(w, r)
		})
	}
}
