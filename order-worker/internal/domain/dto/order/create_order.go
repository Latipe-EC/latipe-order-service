package order

type CreateOrderRequest struct {
	Header        BaseHeader
	UserRequest   UserRequest  `json:"user_request"`
	Amount        int          `json:"amount" validate:"required"`
	Discount      int          `json:"discount" validate:"required"`
	Total         int          `json:"total" validate:"required"`
	PaymentMethod int          `json:"payment_method" validate:"required"`
	Vouchers      []Voucher    `json:"vouchers"`
	Address       OrderAddress `json:"address" validate:"required"`
	Delivery      Delivery     `json:"delivery" validate:"required"`
	OrderItems    []OrderItems `json:"order_items" validate:"required"`
}

type UserRequest struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type Voucher struct {
	Code string `json:"code"`
	Type int    `json:"type"`
}

type OrderItems struct {
	CartItemId string `json:"cart_item_id"`
	ProductId  string `json:"product_id" validate:"required"`
	OptionId   string `json:"option_id"`
	Quantity   int    `json:"quantity" validate:"required"`
	Price      int    `json:"price" validate:"required"`
}

type OrderAddress struct {
	AddressId       string `json:"address_id"`
	ShippingName    string `json:"shipping_name" validate:"required"`
	ShippingPhone   string `json:"shipping_phone" validate:"required"`
	ShippingAddress string `json:"shipping_address" validate:"required"`
}

type Delivery struct {
	DeliveryId    string `json:"delivery_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Cost          int    `json:"cost" validate:"required"`
	ReceivingDate string `json:"receiving_date" validate:"required"`
}
