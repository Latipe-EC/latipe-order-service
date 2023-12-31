// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"encoding/json"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"order-rest-api/config"
	order2 "order-rest-api/internal/api/order"
	"order-rest-api/internal/app/orders"
	"order-rest-api/internal/common/errors"
	"order-rest-api/internal/infrastructure/adapter/authserv"
	"order-rest-api/internal/infrastructure/adapter/deliveryserv"
	"order-rest-api/internal/infrastructure/adapter/productserv"
	"order-rest-api/internal/infrastructure/adapter/storeserv"
	"order-rest-api/internal/infrastructure/adapter/userserv"
	"order-rest-api/internal/infrastructure/adapter/vouchersev"
	"order-rest-api/internal/infrastructure/persistence/db"
	"order-rest-api/internal/infrastructure/persistence/order"
	"order-rest-api/internal/middleware"
	"order-rest-api/internal/middleware/auth"
	"order-rest-api/internal/router"
	"order-rest-api/pkg/cache"
)

// Injectors from server.go:

func New() (*Server, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	gorm := db.NewMySQLConnection(configConfig)
	repository := order.NewGormRepository(gorm)
	service := productserv.NewProductServAdapter(configConfig)
	cacheEngine, err := cache.NewCacheEngine(configConfig)
	if err != nil {
		return nil, err
	}
	userservService := userserv.NewUserServHttpAdapter(configConfig)
	deliveryservService := deliveryserv.NewDeliServHttpAdapter(configConfig)
	voucherservService := voucherserv.NewUserServHttpAdapter(configConfig)
	usecase := orders.NewOrderService(configConfig, repository, service, cacheEngine, userservService, deliveryservService, voucherservService)
	orderApiHandler := order2.NewOrderHandler(usecase)
	orderStatisticApiHandler := order2.NewStatisticHandler(usecase)
	authservService := authserv.NewAuthServHttpAdapter(configConfig)
	storeservService := storeserv.NewStoreServiceAdapter(configConfig)
	authenticationMiddleware := auth.NewAuthMiddleware(authservService, storeservService, deliveryservService, configConfig)
	middlewareMiddleware := middleware.NewMiddleware(authenticationMiddleware)
	orderRouter := router.NewOrderRouter(orderApiHandler, orderStatisticApiHandler, middlewareMiddleware)
	server := NewServer(configConfig, orderRouter)
	return server, nil
}

// server.go:

type Server struct {
	app *fiber.App
	cfg *config.Config
}

func NewServer(
	cfg *config.Config,
	orderRouter router.OrderRouter) *Server {

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
		ErrorHandler: errors.CustomErrorHandler,
	})

	prometheus := fiberprometheus.New("order-rest-api")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Use(logger.New())

	app.Get("", func(ctx *fiber.Ctx) error {
		s := struct {
			Message string `json:"message"`
			Version string `json:"version"`
		}{
			Message: "Order rest-api was developed by TienDat",
			Version: "v0.0.1",
		}
		return ctx.JSON(s)
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
