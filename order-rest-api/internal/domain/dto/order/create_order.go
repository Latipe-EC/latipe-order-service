package order

type CreateOrderRequest struct {
	Header        BaseHeader
	UserRequest   UserRequest
	PaymentMethod int          `json:"payment_method" validate:"required"`
	VoucherCode   []string     `json:"vouchers"`
	Address       OrderAddress `json:"address" validate:"required"`
	Delivery      Delivery     `json:"delivery" validate:"required"`
	OrderItems    []OrderItems `json:"order_items" validate:"required"`
}

type CreateOrderResponse struct {
	UserOrder        UserRequest `json:"user_order"`
	OrderKey         string      `json:"order_key"`
	Amount           int         `json:"amount"`
	ShippingCost     int         `json:"shipping_cost"`
	ShippingDiscount int         ` json:"shipping_discount"`
	ItemDiscount     int         ` json:"item_discount"`
	SubTotal         int         `json:"sub_total" `
	PaymentMethod    int         `json:"payment_method"`
}

type UserRequest struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type OrderItems struct {
	CartId    string `json:"cart_id,omitempty"`
	ProductId string `json:"product_id" validate:"required"`
	OptionId  string `json:"option_id"`
	Quantity  int    `json:"quantity" validate:"required"`
	Price     int    `json:"price" validate:"required"`
}

type OrderAddress struct {
	AddressId string `json:"address_id" validate:"required"`
}

type Delivery struct {
	DeliveryId string `json:"delivery_id" validate:"required"`
}
