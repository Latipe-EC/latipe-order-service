package order

import (
	"order-service-rest-api/internal/domain/entities/order"
	"order-service-rest-api/internal/domain/entities/payment"
	"order-service-rest-api/pkg/db/gorm"
	"order-service-rest-api/pkg/util/pagable"
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
		&payment.PaymentLog{},
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

func (g GormRepository) Save(dao order.Order) error {
	err := g.client.DB().Model(&order.Order{}).Save(&dao)
	return err.Error
}

func (g GormRepository) Update(order order.Order) error {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) Total(query *pagable.Query) (int, error) {
	//TODO implement me
	panic("implement me")
}
