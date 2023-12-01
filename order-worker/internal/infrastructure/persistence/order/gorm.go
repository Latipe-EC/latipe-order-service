package order

import (
	gormF "gorm.io/gorm"
	"order-worker/internal/domain/entities/custom"
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

func (g GormRepository) Save(dao *order.Order) error {
	result := g.client.DB().Model(&order.Order{}).Create(&dao)
	return result.Error
}

func (g GormRepository) UpdateOrderRating(itemId string, ratingId string) error {
	result := g.client.DB().Model(&order.OrderItem{}).
		Where("id = ?", itemId).Update("rating_id", ratingId)

	if result.Error != nil || result.RowsAffected == 0 {
		return result.Error
	}
	return nil
}

func (g GormRepository) GetOrderAmountOfStore(orderId int) ([]custom.AmountItemOfStoreInOrder, error) {
	var result []custom.AmountItemOfStoreInOrder

	err := g.client.DB().Table("orders").
		Select("order_items.store_id as store_id, SUM(order_items.sub_total) as order_amount").
		Joins("INNER JOIN order_items ON orders.id = order_items.order_id").
		Where("orders.id = ?", orderId).
		Group("order_items.store_id").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) FindAllFinishShippingOrder() ([]order.Order, error) {
	var data []order.Order

	err := g.client.DB().
		Where("orders.status =?", order.ORDER_SHIPPING_FINISH).
		Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (g GormRepository) CreateOrderCommmsionTransaction(dao *order.Order, ocms *order.OrderCommission, log *order.OrderStatusLog) error {
	err := g.client.DB().Transaction(func(tx *gormF.DB) error {
		if err := tx.Create(&ocms).Error; err != nil {
			return err
		}

		if err := tx.Model(&order.Order{}).Where("id=?", dao.Id).Update("status", dao.Status).Error; err != nil {
			return err
		}

		if err := tx.Create(&log).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (g GormRepository) FindById(Id string) (*order.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) FindByUserId(query *pagable.Query) ([]*order.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) Update(order order.Order) error {
	//TODO implement me
	panic("implement me")
}

func (g GormRepository) Total(query *pagable.Query) (int, error) {
	//TODO implement me
	panic("implement me")
}
