package order

import "time"

type BaseHeader struct {
	BearerToken string `reqHeader:"Authorization"`
}

type OrderResponse struct {
	OrderUUID     string           `json:"order_uuid"`
	Amount        int              `json:"amount"`
	Discount      int              `json:"discount"`
	SubTotal      int              `json:"sub_total"`
	Status        int              `json:"status"`
	PaymentMethod int              `json:"payment_method"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated-at"`
	Delivery      DeliveryResp     `json:"delivery"`
	OrderItems    []OrderItemsResp `json:"order_items,omitempty"`
}

type DeliveryResp struct {
	DeliveryId      string    `json:"delivery_id"`
	DeliveryName    string    `json:"delivery_name"`
	Cost            int       `json:"payment_type"`
	ReceivingDate   time.Time `json:"receiving_date"`
	AddressId       string    `json:"address_id"`
	ShippingName    string    `json:"shipping_name" `
	ShippingPhone   string    `json:"shipping_phone" `
	ShippingAddress string    `json:"shipping_address" `
}

type OrderItemsResp struct {
	ProductId string `json:"product_id" `
	OptionId  string `json:"option_id"`
	Quantity  int    `json:"quantity" `
	Price     int    `json:"price" `
}
