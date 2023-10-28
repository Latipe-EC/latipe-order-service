//go:build wireinject
// +build wireinject

package server

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/wire"
	"order-rest-api/config"
	"order-rest-api/internal/api"
	"order-rest-api/internal/app"
	"order-rest-api/internal/common/errors"
	"order-rest-api/internal/infrastructure/adapter/authserv"
	"order-rest-api/internal/infrastructure/adapter/deliveryserv"
	"order-rest-api/internal/infrastructure/adapter/productserv"
	"order-rest-api/internal/infrastructure/adapter/userserv"
	"order-rest-api/internal/infrastructure/persistence"
	"order-rest-api/internal/middleware"
	router2 "order-rest-api/internal/router"
	"order-rest-api/pkg/cache"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
}

func New() (*Server, error) {
	panic(wire.Build(wire.NewSet(
		NewServer,
		config.Set,
		api.Set,
		router2.Set,
		persistence.Set,
		userserv.Set,
		authserv.Set,
		deliveryserv.Set,
		productserv.Set,
		app.Set,
		middleware.Set,
		cache.Set,
	)))
}

func NewServer(
	cfg *config.Config,
	orderRouter router2.OrderRouter) *Server {

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
		ErrorHandler: errors.CustomErrorHandler,
	})
	// Initialize default config
	app.Use(logger.New())

	app.Get("", func(ctx *fiber.Ctx) error {
		return ctx.JSON("Orders service developed by Tien Dat")
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	orderRouter.Init(&v1)

	return &Server{
		cfg: cfg,
		app: app,
	}
}

func (serv Server) App() *fiber.App {
	return serv.app
}

func (serv Server) Config() *config.Config {
	return serv.cfg
}
