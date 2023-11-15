package order

type CancelOrderRequest struct {
	Header    BaseHeader
	OrderUUID string `json:"order_uuid" validate:"required"`
	UserId    string
}

type CancelOrderResponse struct {
}
