package order

import "order-service-rest-api/pkg/util/pagable"

type GetOrderRequest struct {
	BaseHeader BaseHeader
	Query      pagable.Query
}
type GetOrderResponse struct {
	Order OrderResponse `json:"order"`
}
