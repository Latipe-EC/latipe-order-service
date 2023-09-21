//go:build wireinject
// +build wireinject

package server

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/wire"
	"order-service-rest-api/config"
	"order-service-rest-api/internal/api"
	"order-service-rest-api/internal/common/errors"
	"order-service-rest-api/internal/infrastructure/persistence"
	"order-service-rest-api/internal/message"
	"order-service-rest-api/internal/router"
)

type Server struct {
	app      *fiber.App
	cfg      *config.Config
	consumer *message.ConsumerOrderMessage
}

func New() (*Server, error) {
	panic(wire.Build(wire.NewSet(
		NewServer,
		config.Set,
		message.Set,
		api.Set,
		router.Set,
		persistence.Set,
	)))
}

func NewServer(
	cfg *config.Config,
	orderRouter router.OrderRouter,
	consumer *message.ConsumerOrderMessage) *Server {

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
		cfg:      cfg,
		app:      app,
		consumer: consumer,
	}
}

func (serv Server) App() *fiber.App {
	return serv.app
}

func (serv Server) Config() *config.Config {
	return serv.cfg
}

func (serv Server) Consumer() *message.ConsumerOrderMessage {
	return serv.consumer
}
