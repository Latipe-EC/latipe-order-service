package order

type UpdateOrderStatusRequest struct {
	Role    string
	OrderId int `json:"order_id"`
	Status  int `json:"status"`
}
type UpdateOrderStatusResponse struct {
}

type UpdateOrderRequest struct {
	Header  BaseHeader
	OrderId int `json:"order_id"`
	Status  int `json:"status"`
}
