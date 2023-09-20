//go:build wireinject
// +build wireinject

package server

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"order-service-rest-api/config"
	"order-service-rest-api/internal/message"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
}

func New() (*Server, error) {
	panic(wire.Build(wire.NewSet(
		NewServer,
		config.Set,
		message.Set,
	)))
}

func NewServer(
	cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})
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
