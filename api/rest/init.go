package rest

import (
	"main/api/rest/op"
	"main/core/node"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func Init(root *node.Node) *chi.Mux {
	router := chi.NewMux()
	addMiddlewares(router)

	api := humachi.New(router, huma.DefaultConfig("Configuration API", "1.0.0"))
	op.RegisterCamRoutes(api, root)
	op.RegisterGroupRoutes(api, root)
	op.RegisterListRoutes(api, root)
	op.RegisterMediaRoutes(api, root)

	return router
}
