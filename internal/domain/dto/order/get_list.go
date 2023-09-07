package order

import "order-service-rest-api/pkg/util/pagable"

type GetOrderListRequest struct {
	BaseHeader BaseHeader
	Query      pagable.Query
}
type GetOrderListResponse struct {
	Data pagable.ListResponse
}
