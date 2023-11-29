package router

import (
	"github.com/gofiber/fiber/v2"
	"order-rest-api/internal/api/order"
	"order-rest-api/internal/middleware"
)

type OrderRouter interface {
	Init(root *fiber.Router)
}

type orderRouter struct {
	handler    order.OrderApiHandler
	middleware *middleware.Middleware
}

func NewOrderRouter(handler order.OrderApiHandler, middleware *middleware.Middleware) OrderRouter {
	return orderRouter{
		handler:    handler,
		middleware: middleware,
	}
}

func (o orderRouter) Init(root *fiber.Router) {
	orderRouter := (*root).Group("/orders")
	{
		//admin
		orderRouter.Get("/admin", o.middleware.Authentication.RequiredAuthentication(), o.handler.ListOfOrder)
		orderRouter.Patch("/admin/:id", o.middleware.Authentication.RequiredAuthentication(), o.handler.UpdateOrderStatus)
		orderRouter.Patch("/admin/:id/complete", o.middleware.Authentication.RequiredAuthentication(), o.handler.UpdateOrderStatus)
		//user
		orderRouter.Post("/user", o.middleware.Authentication.RequiredAuthentication(), o.handler.CreateOrder)
		orderRouter.Patch("/user/cancel", o.middleware.Authentication.RequiredAuthentication(), o.handler.CancelOrder)
		orderRouter.Get("/user", o.middleware.Authentication.RequiredAuthentication(), o.handler.GetMyOrder)
		orderRouter.Get("/user/:id", o.middleware.Authentication.RequiredAuthentication(), o.handler.GetOrderByUUID)
		orderRouter.Patch("/:id/items", o.middleware.Authentication.RequiredAuthentication(), o.handler.UpdateOrderStatus)
		//store
		orderRouter.Get("/store", o.middleware.Authentication.RequiredStoreAuthentication(), o.handler.GetMyStoreOrder)
		orderRouter.Get("/store/:id", o.middleware.Authentication.RequiredStoreAuthentication(), o.handler.GetStoreOrderDetail)
		orderRouter.Patch("/store/:id/items", o.middleware.Authentication.RequiredStoreAuthentication(), o.handler.UpdateOrderItemStatus)
		//delivery
		orderRouter.Get("/delivery", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.handler.GetOrdersByDelivery)
		orderRouter.Patch("/delivery/:id", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.handler.UpdateStatusByDelivery)
	}

}
