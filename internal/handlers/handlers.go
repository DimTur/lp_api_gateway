package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterConfigurator interface {
	ConfigureRouter() http.Handler
}

type ChiRouterConfigurator struct{}

func (c *ChiRouterConfigurator) ConfigureRouter() http.Handler {
	router := chi.NewRouter()

	// Подключаем middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	// Можно добавить маршруты
	// router.Get("/example", exampleHandler)

	return router
}
