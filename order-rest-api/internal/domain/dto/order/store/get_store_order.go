package store

import (
	"order-rest-api/internal/domain/dto/order"
	"order-rest-api/pkg/util/pagable"
)

type GetStoreOrderRequest struct {
	BaseHeader order.BaseHeader
	StoreID    string
	Query      *pagable.Query
}

type GetStoreOrderResponse struct {
	pagable.ListResponse
}
