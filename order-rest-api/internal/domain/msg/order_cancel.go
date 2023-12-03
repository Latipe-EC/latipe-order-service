package msg

import "time"

type OrderCancelMessage struct {
	OrderUUID     string    `json:"orderUUID"`
	OrderStatus   int       `json:"orderStatus"`
	PaymentMethod int       `json:"paymentMethod"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
