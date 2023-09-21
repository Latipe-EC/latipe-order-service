package api

import (
	"github.com/google/wire"
	"order-service-rest-api/internal/api/order"
)

var Set = wire.NewSet(
	order.NewOrderHandler,
)
