package entities

type InternalService struct {
	InternalService int `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	ServiceKey      int `gorm:"not null;type:bigint" json:"service_key"`
	ServiceName     int `gorm:"not null;type:bigint" json:"service_name"`
}
