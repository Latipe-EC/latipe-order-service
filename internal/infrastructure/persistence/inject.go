package persistence

import (
	"github.com/google/wire"
	"order-service-rest-api/internal/infrastructure/persistence/db"
	"order-service-rest-api/internal/infrastructure/persistence/order"
)

var Set = wire.NewSet(
	db.NewMySQLConnection,
	order.NewGormRepository,
)
