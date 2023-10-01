package order

import (
	"order-rest-api/pkg/util/pagable"
)

type GetOrderListRequest struct {
	BaseHeader BaseHeader
	Query      *pagable.Query
}
type GetOrderListResponse struct {
	pagable.ListResponse
}
