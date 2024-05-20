package rest

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func addMiddlewares(router chi.Router) {
	router.Use(middleware.Logger)
	router.Use(cors.AllowAll().Handler)
}
