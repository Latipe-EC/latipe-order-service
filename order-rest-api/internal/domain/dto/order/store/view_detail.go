package store

import (
	"order-rest-api/internal/domain/dto/order"
	"time"
)

type GetOrderOfStoreByIDRequest struct {
	BaseHeader order.BaseHeader
	OrderUUID  string `json:"order_uuid" params:"id"`
	StoreID    string
}

type GetOrderOfStoreByIDResponse struct {
	StoreOrderResponse
}

type StoreOrderResponse struct {
	OrderUUID        string             `json:"order_uuid"`
	StoreOrderAmount int                `json:"store_order_amount,omitempty"`
	Status           int                `json:"status"`
	PaymentMethod    int                `json:"payment_method"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	Delivery         order.DeliveryResp `json:"delivery"`
	OrderItems       []OrderStoreItem   `json:"order_items,omitempty"`
}

type OrderStoreItem struct {
	Id          string `json:"item_id,omitempty"`
	ProductId   string `json:"product_id" `
	OptionId    string `json:"option_id"`
	Quantity    int    `json:"quantity" `
	Price       int    `json:"price" `
	Status      int    `json:"is_prepared"`
	SubTotal    int    `json:"sub_total"`
	ProductName string `json:"product_name"`
	ProdImg     string `json:"image"`
}
