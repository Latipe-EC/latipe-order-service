package order

import "time"

type DeliveryOrder struct {
	Id              int `gorm:"not null;autoIncrement;primaryKey;type:bigint"`
	OrderType       string
	DeliveryId      string    `gorm:"not null;type:varchar(255)" json:"delivery_id"`
	DeliveryName    string    `gorm:"not null;type:varchar(255)" json:"delivery_name"`
	Cost            int       `gorm:"not null;type:int" json:"cost"`
	AddressId       string    `json:"address_id" gorm:"type:varchar(255)"`
	ShippingName    string    `json:"shipping_name" gorm:"not null;type:varchar(255)"`
	ShippingPhone   string    `json:"shipping_phone" gorm:"not null;type:varchar(255)"`
	ShippingAddress string    `json:"shipping_address" gorm:"not null;type:varchar(255)"`
	ReceivingDate   time.Time `json:"receiving_date" gorm:"not null;type:date"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt       time.Time `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
	OrderID         int       `gorm:"not null;type:bigint" json:"order_id"`
	Order           *Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (DeliveryOrder) TableName() string {
	return "delivery_orders"
}
