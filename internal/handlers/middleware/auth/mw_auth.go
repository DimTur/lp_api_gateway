package authmiddleware

import (
	"context"
	"net/http"
	"strconv"

	ssov1 "github.com/DimTur/lp_protos/gen/go/sso"
)

type AuthService interface {
	AuthCheck(ctx context.Context, accessToken string) (*ssov1.AuthCheckResponse, error)
}

func AuthMiddleware(authService AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.Header.Get("Authorization")
			if accessToken == "" {
				w.WriteHeader(http.StatusUnauthorized)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			resp, err := authService.AuthCheck(r.Context(), accessToken)
			if err != nil || !resp.IsValid {
				w.WriteHeader(http.StatusUnauthorized)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userIDStr := strconv.FormatInt(resp.UserId, 10)
			r.Header.Set("X-User-ID", userIDStr)

			next.ServeHTTP(w, r)
		})
	}
}
