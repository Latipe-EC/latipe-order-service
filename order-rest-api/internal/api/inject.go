package api

import (
	"github.com/google/wire"
	"order-rest-api/internal/api/order"
)

var Set = wire.NewSet(
	order.NewOrderHandler,
	order.NewStatisticHandler,
)
