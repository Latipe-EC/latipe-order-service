package order

type UpdateOrderStatusRequest struct {
	Role      string
	OrderUUID string `json:"order_uuid"`
	UserId    string
	Status    int `json:"status"`
}
type UpdateOrderStatusResponse struct {
}

type UpdateOrderRequest struct {
	Header    BaseHeader
	OrderUUID string `json:"order_uuid"`
	UserId    string
	Status    int `json:"status"`
}
