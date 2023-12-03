package order

type GetOrderByIDRequest struct {
	BaseHeader BaseHeader
	Role       string
	OwnerId    string
	OrderId    int `json:"order_id" params:"id"`
}

type GetOrderByUUIDRequest struct {
	BaseHeader BaseHeader
	Role       string
	OwnerId    string
	OrderId    string `json:"order_id" params:"id"`
}

type GetOrderResponse struct {
	Order OrderResponse `json:"order"`
}
