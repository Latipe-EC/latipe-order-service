package order

type UpdateOrderRequest struct {
	Header  BaseHeader
	OrderId int `json:"order_id"`
	Status  int `json:"status"`
}
type UpdateOrderResponse struct {
}
