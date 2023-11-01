package order

import (
	entity "order-rest-api/internal/domain/entities/order"
	"order-rest-api/pkg/db/gorm"
	"order-rest-api/pkg/util/pagable"
)

type GormRepository struct {
	client gorm.Gorm
}

func NewGormRepository(client gorm.Gorm) entity.Repository {
	// auto migrate
	err := client.DB().AutoMigrate(
		&entity.Order{},
		&entity.OrderItem{},
		&entity.OrderStatusLog{},
		&entity.DeliveryOrder{},
	)
	if err != nil {
		panic(err)
	}
	return &GormRepository{
		client: client,
	}
}

func (g GormRepository) FindById(Id int) (*entity.Order, error) {
	order := entity.Order{}

	result := g.client.DB().Model(&entity.Order{}).
		Preload("OrderItem").
		First(&order, Id).Error
	if result != nil {
		return nil, result
	}

	return &order, nil
}

func (g GormRepository) FindAll(query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	whereState := query.ORMConditions().(string)
	result := g.client.DB().Model(&entity.Order{}).
		Where(whereState).
		Limit(query.GetLimit()).Offset(query.GetOffset()).
		Find(&orders).Error
	if result != nil {
		return nil, result
	}

	return orders, nil
}

func (g GormRepository) FindByUserId(userId int, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	result := g.client.DB().Model(&entity.Order{}).
		Where("order.user_id", userId).
		Order("create_at desc").
		Limit(query.GetLimit()).Offset(query.GetOffset()).
		Find(&orders).Error
	if result != nil {
		return nil, result
	}

	return orders, nil
}

func (g GormRepository) FindOrderByUserAndProduct(userId string, productId string) ([]entity.Order, error) {
	var orders []entity.Order
	err := g.client.DB().Raw("select * from orders inner join order_items on orders.id = order_items.order_id "+
		"where orders.user_id= ? and order_items.product_id = ?", userId, productId).Scan(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, err
}

func (g GormRepository) FindOrderLogByOrderId(orderId int) ([]entity.OrderStatusLog, error) {
	var orderStatus []entity.OrderStatusLog
	result := g.client.DB().Model(&entity.OrderStatusLog{}).
		Where("order_id", orderId).
		Order("create_at desc").
		Find(&orderStatus).Error
	if result != nil {
		return nil, result
	}

	return orderStatus, nil
}

func (g GormRepository) Save(dao *entity.Order) error {
	result := g.client.DB().Model(&entity.Order{}).Create(&dao)
	return result.Error
}

func (g GormRepository) Update(order entity.Order) error {
	result := g.client.DB().Updates(order)

	if result.Error != nil || result.RowsAffected == 0 {
		return result.Error
	}
	return nil
}

func (g GormRepository) Total(query *pagable.Query) (int, error) {
	var count int64
	whereState := query.ORMConditions().(string)
	result := g.client.DB().Select("*").Table(entity.Order{}.TableName()).
		Where(whereState).
		Count(&count).Error

	return int(count), result
}
