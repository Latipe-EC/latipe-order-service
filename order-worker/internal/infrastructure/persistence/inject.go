package persistence

import (
	"order-worker/internal/infrastructure/persistence/db"
	"order-worker/internal/infrastructure/persistence/order"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	db.NewMySQLConnection,
	order.NewGormRepository,
)
