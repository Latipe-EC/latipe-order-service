package delivery

import (
	"order-rest-api/internal/infrastructure/adapter/productserv/dto"
	"order-rest-api/pkg/util/pagable"
)

type GetOrderListRequest struct {
	BaseHeader dto.BaseHeader
	DeliveryID string
	Query      *pagable.Query
}

type GetOrderListResponse struct {
	pagable.ListResponse
}
