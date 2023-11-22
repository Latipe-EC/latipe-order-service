package order

import (
	"order-worker/internal/domain/entities/order"
	"order-worker/pkg/db/gorm"
	"order-worker/pkg/util/pagable"
)

type GormRepository struct {
	client gorm.Gorm
}

func NewGormRepository(client gorm.Gorm) order.Repository {
	// auto migrate
	err := client.DB().AutoMigrate(
		&order.Order{},
		&order.OrderItem{},
		&order.OrderStatusLog{},
		&order.DeliveryOrder{},
		&order.OrderCommission{},
	)
	if err != nil {
		panic(err)
	}
	return &GormRepository{
		client: client,
	}
}

func (g GormRepository) FindById(Id string) (*order.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) FindAll(query *pagable.Query) ([]*order.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) FindByUserId(query *pagable.Query) ([]*order.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) FindOrderLogByOrderId(orderId int) ([]*order.OrderStatusLog, error) {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) Save(dao *order.Order) error {
	result := g.client.DB().Model(&order.Order{}).Create(&dao)
	return result.Error
}

func (g GormRepository) Update(order order.Order) error {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) Total(query *pagable.Query) (int, error) {
	//TODO implement me
	panic("implement me")
}
