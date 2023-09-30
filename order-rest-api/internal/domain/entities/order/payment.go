package order

import (
	"time"
)

type PaymentLog struct {
	Id                 int `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	OrderType          string
	OrderID            int       `gorm:"not null;type:bigint" json:"order_id"`
	Order              *Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PaymentTransaction string    `gorm:"not null;type:varchar(255)" json:"payment_transaction"`
	PaymentType        int       `gorm:"not null;type:varchar(255)" json:"payment_type"`
	Total              int       `gorm:"not null;" json:"total"`
	ThirdPartyLog      string    `gorm:"TEXT" json:"third_party_log"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt          time.Time `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
}

func (PaymentLog) TableName() string {
	return "payments"
}
