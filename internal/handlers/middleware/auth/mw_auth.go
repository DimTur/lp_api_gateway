package authmiddleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	ssov1 "github.com/DimTur/lp_protos/gen/go/sso"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type AuthService interface {
	AuthCheck(ctx context.Context, accessToken string) (*ssov1.AuthCheckResponse, error)
}

func AuthMiddleware(log *slog.Logger, authService AuthService) func(next http.Handler) http.Handler {
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
			resp, err := authService.AuthCheck(r.Context(), accessToken)
			if err != nil {
				log.Error("error checking authorization", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusUnauthorized)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !resp.IsValid {
				log.Info("invalid authorization token", slog.String("user_id", strconv.FormatInt(resp.UserId, 10)))
				w.WriteHeader(http.StatusUnauthorized)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			span.AddEvent("completed_auth_cheking")
			span.SetAttributes(attribute.Int64("userID", resp.UserId))

			userIDStr := strconv.FormatInt(resp.UserId, 10)
			r.Header.Set("X-User-ID", userIDStr)

			log.Info("authorization successful", slog.String("user_id", userIDStr))

			next.ServeHTTP(w, r)
		})
	}
}
