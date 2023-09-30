package app

import (
	"order-worker/internal/app/orders"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	orders.NewOrderService,
)
