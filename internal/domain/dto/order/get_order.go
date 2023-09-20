package order

type GetOrderRequest struct {
	BaseHeader BaseHeader
	OrderId    int `json:"order_id"`
}
type GetOrderResponse struct {
	Order OrderResponse `json:"order"`
}
