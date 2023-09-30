package app

import (
	"github.com/google/wire"
	"order-service-rest-api/internal/app/orders"
)

var Set = wire.NewSet(
	orders.NewOrderService,
)
