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
		adminRouter := orderRouter.Group("/admin")
		{
			adminRouter.Get("/:id", o.middleware.Authentication.RequiredAuthentication(), o.handler.GetOrderByUUID)
			adminRouter.Get("", o.middleware.Authentication.RequiredAuthentication(), o.handler.ListOfOrder)
			adminRouter.Get("/count", o.middleware.Authentication.RequiredAuthentication(), o.handler.AdminCountingOrder)
			adminRouter.Patch("/:id", o.middleware.Authentication.RequiredAuthentication(), o.handler.UpdateOrderStatus)
			adminRouter.Patch("/:id/complete", o.middleware.Authentication.RequiredAuthentication(), o.handler.UpdateOrderStatus)
		}

		//user
		userRouter := orderRouter.Group("/user")
		{
			userRouter.Get("", o.middleware.Authentication.RequiredAuthentication(), o.handler.GetMyOrder)
			userRouter.Get("/count", o.middleware.Authentication.RequiredAuthentication(), o.handler.UserCountingOrder)
			userRouter.Get("/:id", o.middleware.Authentication.RequiredAuthentication(), o.handler.UserGetOrderByUUID)
			userRouter.Post("", o.middleware.Authentication.RequiredAuthentication(), o.handler.CreateOrder)
			userRouter.Patch("/cancel", o.middleware.Authentication.RequiredAuthentication(), o.handler.CancelOrder)
		}

		//store
		storeRouter := orderRouter.Group("/store")
		{
			storeRouter.Get("", o.middleware.Authentication.RequiredStoreAuthentication(), o.handler.GetMyStoreOrder)
			storeRouter.Get("/count", o.middleware.Authentication.RequiredStoreAuthentication(), o.handler.StoreCountingOrder)
			storeRouter.Get("/:id", o.middleware.Authentication.RequiredStoreAuthentication(), o.handler.GetStoreOrderDetail)
			storeRouter.Patch("/:id/items", o.middleware.Authentication.RequiredStoreAuthentication(), o.handler.UpdateOrderItemStatus)
		}

		//delivery
		deliveryRouter := orderRouter.Group("/delivery")
		{
			deliveryRouter.Get("", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.handler.GetOrdersByDelivery)
			deliveryRouter.Get("/:id", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.handler.DeliveryGetOrderByUUID)
			deliveryRouter.Get("/count", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.handler.DeliveryCountingOrder)
			deliveryRouter.Patch("/:id", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.handler.UpdateStatusByDelivery)
		}

		//internal
		internalRouter := orderRouter.Group("/internal")
		{
			internalRouter.Get("/rating/:id", o.middleware.Authentication.RequiredInternalService(), o.handler.InternalGetOrderByUUID)
		}

	}

}
