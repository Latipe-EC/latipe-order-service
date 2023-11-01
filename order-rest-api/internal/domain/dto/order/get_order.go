package order

type GetOrderByIDRequest struct {
	BaseHeader BaseHeader
	OrderId    int `json:"order_id" params:"id"`
}

type GetOrderByUUIDRequest struct {
	BaseHeader BaseHeader
	OrderId    string `json:"order_id" params:"id"`
}

type GetOrderResponse struct {
	Order OrderResponse `json:"order"`
}
