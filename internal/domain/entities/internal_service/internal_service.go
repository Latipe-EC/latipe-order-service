package internal_service

type InternalService struct {
	InternalServiceID int    `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	ServiceKey        string `gorm:"not null;type:bigint" json:"service_key"`
	ServiceValue      string `gorm:"not null;type:bigint" json:"service_value"`
	ServiceName       int    `gorm:"not null;type:bigint" json:"service_name"`
}

func (InternalService) TableName() string {
	return "internal_services"
}
