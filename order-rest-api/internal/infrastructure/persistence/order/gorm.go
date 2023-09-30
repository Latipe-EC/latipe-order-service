package order

import (
	order2 "order-rest-api/internal/domain/entities/order"
	"order-rest-api/pkg/db/gorm"
	"order-rest-api/pkg/util/pagable"
)

type GormRepository struct {
	client gorm.Gorm
}

func NewGormRepository(client gorm.Gorm) order2.Repository {
	// auto migrate
	err := client.DB().AutoMigrate(
		&order2.Order{},
		&order2.OrderItem{},
		&order2.OrderStatusLog{},
		&order2.PaymentLog{},
	)
	if err != nil {
		panic(err)
	}
	return &GormRepository{
		client: client,
	}
}

func (g GormRepository) FindById(Id string) (*order2.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) FindAll(query *pagable.Query) ([]*order2.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) FindByUserId(query *pagable.Query) ([]*order2.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) FindOrderLogByOrderId(orderId int) ([]*order2.OrderStatusLog, error) {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) Save(dao *order2.Order) error {
	result := g.client.DB().Model(&order2.Order{}).Create(&dao)
	return result.Error
}

func (g GormRepository) Update(order order2.Order) error {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) Total(query *pagable.Query) (int, error) {
	//TODO implement me
	panic("implement me")
}
