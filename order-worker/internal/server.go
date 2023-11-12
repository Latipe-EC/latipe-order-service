//go:build wireinject
// +build wireinject

package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/wire"
	"order-worker/config"
	"order-worker/internal/app"

	"order-worker/internal/infrastructure/adapter/productserv"
	"order-worker/internal/infrastructure/adapter/userserv"
	"order-worker/internal/infrastructure/persistence"
	"order-worker/internal/message"
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
		persistence.Set,
		userserv.Set,
		productserv.Set,
		app.Set,
	)))
}

func NewServer(
	cfg *config.Config,
	consumer *message.ConsumerOrderMessage) *Server {

	app := fiber.New(fiber.Config{})
	// Initialize default config
	app.Use(logger.New())

	app.Get("", func(ctx *fiber.Ctx) error {
		return ctx.JSON("Orders service developed by Tien Dat")
	})

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
