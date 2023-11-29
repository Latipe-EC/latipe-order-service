package order

type BaseHeader struct {
	BearerToken string `reqHeader:"Authorization"`
}

type OrderMessage struct {
	Header        BaseHeader
	UserRequest   UserRequest       `json:"user_request,omitempty"`
	Status        int               `json:"status"`
	OrderUUID     string            `json:"order_uuid,omitempty"`
	Amount        int               `json:"amount,omitempty" validate:"required"`
	ShippingCost  int               `json:"shipping_cost,omitempty"`
	Discount      int               `json:"discount,omitempty" validate:"required"`
	SubTotal      int               `json:"sub_total,omitempty" validate:"required"`
	PaymentMethod int               `json:"payment_method,omitempty" validate:"required"`
	Vouchers      []string          `json:"vouchers,omitempty"`
	Address       OrderAddress      `json:"address,omitempty" validate:"required"`
	Delivery      Delivery          `json:"delivery,omitempty" validate:"required"`
	OrderItems    []OrderItemsCache `json:"order_items,omitempty" validate:"required"`
}

type UserRequest struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type OrderItemsCache struct {
	CartId      string      `json:"cart_id"`
	ProductItem ProductItem `json:"product_item"`
}

type ProductItem struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	NameOption  string `json:"name_option"`
	StoreID     string `json:"store_id"`
	OptionID    string `json:"option_id" `
	Image       string `json:"image"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	NetPrice    int    `json:"net_price"`
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
