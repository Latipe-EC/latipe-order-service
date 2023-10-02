package order

import (
	"time"
)

type OrderItem struct {
	OrderType string
	Id        int       `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	OrderID   int       `gorm:"not null;type:bigint" json:"order_id"`
	Order     *Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductID string    `gorm:"not null;type:varchar(255)" json:"product_id"`
	StoreID   int       `gorm:"not null;type:int" json:"store_id"`
	OptionID  string    `gorm:"not null;type:varchar(250)" json:"option_id" `
	Quantity  int       `gorm:"not null;type:int" json:"quantity"`
	Price     int       `gorm:"not null;type:bigint" json:"price"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt time.Time `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
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
	Id             int               `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	UserId         string            `gorm:"not null;type:varchar(250)" json:"user_id"`
	Username       string            `gorm:"not null;type:varchar(250)" json:"email"`
	Amount         int               `gorm:"not null;type:bigint" json:"amount"`
	Discount       int               `gorm:"not null;type:int" json:"discount"`
	Total          int               `gorm:"not null;type:int" json:"total"`
	Status         int               `gorm:"not null;type:int" json:"status"`
	AddressID      int               `gorm:"not null;type:bigint" json:"address_id"`
	AddressDetail  string            `gorm:"not null;type:longtext" json:"address_detail"`
	VoucherCode    string            `gorm:"not null;type:varchar(250)" json:"voucher_code"`
	UpdatedAt      time.Time         `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt      time.Time         `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
	OrderItem      []*OrderItem      `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"order_items"`
	OrderStatusLog []*OrderStatusLog `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"order_status_logs"`
	PaymentLog     *PaymentLog       `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"payment_log"`
	Delivery       *DeliveryOrder    `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"delivery"`
}

func (Order) TableName() string {
	return "order"
}
