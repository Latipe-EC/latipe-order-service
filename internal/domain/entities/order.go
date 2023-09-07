package entities

import "time"

type Orders struct {
	OrderID       int           `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	Amount        int           `gorm:"not null;type:bigint" json:"amount"`
	Discount      int           `gorm:"not null;type:int" json:"discount"`
	Total         int           `gorm:"not null;type:int" json:"total"`
	Status        int           `gorm:"not null;type:int" json:"status"`
	AddressID     int           `gorm:"not null;type:bigint" json:"address_id"`
	AddressDetail string        `gorm:"not null;type:longtext" json:"address_detail"`
	VoucherCode   string        `gorm:"not null;type:varchar(250)" json:"voucher_code"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt     time.Time     `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
	OrderItems    []*OrderItems `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Orders;" json:"order_items"`
	PaymentMethod int           `gorm:"not null;type:int" json:"payments"`
}

type OrderItems struct {
	ItemID    int       `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	OrderID   int       `gorm:"not null;type:bigint" json:"order_id"`
	ProductID int       `gorm:"not null;type:bigint" json:"product_id"`
	SellerID  int       `gorm:"not null;type:int" json:"seller_id"`
	OptionID  int       `gorm:"not null;type:bigint" json:"option_id" `
	Quantity  int       `gorm:"not null;type:int" json:"quantity"`
	Price     int       `gorm:"not null;type:bigint" json:"price"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt time.Time `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
}

type OrderStatusLog struct {
	Id           int       `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	LogId        string    `gorm:"not null;type:string" json:"LogId"`
	OrderId      int       `gorm:"not null;type:bigint" json:"order_id"`
	Message      string    `gorm:"type:longtext" json:"message"`
	StatusChange int       `gorm:"type:int" json:"status_change"`
	Order        *Orders   `gorm:"polymorphic:Orders" json:"order"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
}
