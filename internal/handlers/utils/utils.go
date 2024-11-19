package utils

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func GetHeaderID(r *http.Request, header string) (string, error) {
	id := r.Header.Get(header)
	if id == "" {
		return "", fmt.Errorf("missing %s in headers", header)
	}
	return id, nil
}

func GetURLParamInt64(r *http.Request, param string) (int64, error) {
	paramStr := chi.URLParam(r, param)
	if paramStr == "" {
		return 0, fmt.Errorf("missing %s in URL params", param)
	}
	p, err := strconv.ParseInt(paramStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s in URL params: %w", param, err)
	}
	return p, nil
}

func DecodeRequestBody[T any](r *http.Request, log *slog.Logger) (*T, error) {
	var req T
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}
	return &req, nil
}
