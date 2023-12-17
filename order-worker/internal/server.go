//go:build wireinject
// +build wireinject

package server

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/wire"
	"order-worker/config"
	"order-worker/internal/app"
	"order-worker/internal/infrastructure/adapter/storeserv"
	voucherserv "order-worker/internal/infrastructure/adapter/vouchersev"
	"order-worker/internal/order_cron"
	"order-worker/internal/publisher"
	"order-worker/internal/worker"

	"order-worker/internal/infrastructure/adapter/productserv"
	"order-worker/internal/infrastructure/adapter/userserv"
	"order-worker/internal/infrastructure/persistence"
)

type Server struct {
	app                  *fiber.App
	cfg                  *config.Config
	orderCreateConsumer  *worker.ConsumerOrderMessage
	ratingUpdateConsumer *worker.ConsumerRatingMessage
	orderCompleteCJ      *order_cron.OrderCompleteCronjob
}

func New() (*Server, error) {
	panic(wire.Build(wire.NewSet(
		NewServer,
		config.Set,
		publisher.Set,
		worker.Set,
		persistence.Set,
		userserv.Set,
		productserv.Set,
		storeserv.Set,
		voucherserv.Set,
		app.Set,
		order_cron.Set,
	)))
}

func NewServer(
	cfg *config.Config,
	orderSubscriber *worker.ConsumerOrderMessage,
	ratingSubscriber *worker.ConsumerRatingMessage,
	orderCompleteCron *order_cron.OrderCompleteCronjob) *Server {

	app := fiber.New(fiber.Config{})
	// Initialize default config
	app.Use(logger.New())
	//metric
	prometheus := fiberprometheus.New("order-worker")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Get("", func(ctx *fiber.Ctx) error {
		s := struct {
			Message string `json:"message"`
			Version string `json:"version"`
		}{
			Message: "Order worker was developed by TienDat",
			Version: "v0.0.1",
		}
		return ctx.JSON(s)
	})

	return &Server{
		cfg:                  cfg,
		app:                  app,
		ratingUpdateConsumer: ratingSubscriber,
		orderCreateConsumer:  orderSubscriber,
		orderCompleteCJ:      orderCompleteCron,
	}
}

func (serv Server) App() *fiber.App {
	return serv.app
}

func (serv Server) Config() *config.Config {
	return serv.cfg
}

func (serv Server) ConsumerOrderMessage() *worker.ConsumerOrderMessage {
	return serv.orderCreateConsumer
}

func (serv Server) ConsumerRatingMessage() *worker.ConsumerRatingMessage {
	return serv.ratingUpdateConsumer
}

func (serv Server) OrderCompleteCJ() *order_cron.OrderCompleteCronjob {
	return serv.orderCompleteCJ
}
