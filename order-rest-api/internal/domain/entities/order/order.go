package order

import (
	"time"
)

const (
	OI_PENDING  = 0
	OI_PREPARED = 1
)

type OrderItem struct {
	OrderType   string
	Id          int       `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	OrderID     int       `gorm:"not null;type:bigint" json:"order_id"`
	Order       *Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductID   string    `gorm:"not null;type:varchar(255)" json:"product_id"`
	ProductName string    `gorm:"not null;type:varchar(255)" json:"product_name"`
	ProdImg     string    `gorm:"not null;type:TEXT" json:"image"`
	RatingID    string    `gorm:"not null;type:varchar(255)" json:"rating_id"`
	Status      int       `gorm:"not null;type:int" json:"status"`
	StoreID     string    `gorm:"not null;type:varchar(255)" json:"store_id"`
	OptionID    string    `gorm:"not null;type:varchar(250)" json:"option_id" `
	Quantity    int       `gorm:"not null;type:int" json:"quantity"`
	Price       int       `gorm:"not null;type:bigint" json:"price"`
	NetPrice    int       `gorm:"not null;type:bigint" json:"net_price"`
	SubTotal    int       `gorm:"not null;type:bigint" json:"sub_total"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

type OrderStatusLog struct {
	OrderType    string
	Id           int       `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	OrderID      int       `gorm:"not null;type:bigint" json:"order_id"`
	Order        *Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Message      string    `gorm:"type:longtext" json:"message"`
	StatusChange int       `gorm:"type:int" json:"status_change"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
}

func (OrderStatusLog) TableName() string {
	return "order_status_logs"
}

type Order struct {
	Id               int                `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	OrderUUID        string             `gorm:"column:order_uuid;not null;type:varchar(250)" json:"order_uuid"`
	UserId           string             `gorm:"not null;type:varchar(250)" json:"user_id"`
	Username         string             `gorm:"not null;type:varchar(250)" json:"email"`
	Amount           int                `gorm:"not null;type:bigint" json:"amount"`
	Discount         int                `gorm:"not null;type:int" json:"discount"`
	ShippingCost     int                `gorm:"not null;type:int" json:"shipping_cost"`
	SubTotal         int                `gorm:"not null;type:int" json:"sub_total"`
	Status           int                `gorm:"not null;type:int" json:"status"`
	PaymentMethod    int                `json:"payment_method" gorm:"not null;type:int"`
	UpdatedAt        time.Time          `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt        time.Time          `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
	OrderItem        []*OrderItem       `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"order_items"`
	OrderCommissions []*OrderCommission `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"order_commissions"`
	OrderStatusLog   []*OrderStatusLog  `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"order_status"`
	Delivery         *DeliveryOrder     `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"delivery"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderCommission struct {
	OrderType      string
	Id             int       `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	OrderID        int       `gorm:"not null;type:bigint" json:"order_id"`
	Order          *Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StoreID        string    `gorm:"not null;type:varchar(250)" json:"store_id"`
	AmountReceived int       `gorm:"not null;type:bigint" json:"amount_received"`
	SystemFee      int       `gorm:"not null;type:bigint" json:"system_fee"`
	CreatedAt      time.Time `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
}

func (OrderCommission) TableName() string {
	return "orders_commission"
}
