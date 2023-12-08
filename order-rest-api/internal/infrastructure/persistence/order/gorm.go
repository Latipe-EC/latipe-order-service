package order

import (
	"context"
	"fmt"
	gormF "gorm.io/gorm"
	entity "order-rest-api/internal/domain/entities/order"
	"order-rest-api/pkg/db/gorm"
	"order-rest-api/pkg/util/pagable"
	"strings"
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
		&entity.OrderCommission{},
	)
	if err != nil {
		panic(err)
	}
	return &GormRepository{
		client: client,
	}
}

func (g GormRepository) FindById(ctx context.Context, Id int) (*entity.Order, error) {
	order := entity.Order{}

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.Order{}).
			Preload("OrderItem").
			Preload("Delivery").
			First(&order, Id).Error
	}, ctx)
	if result != nil {
		return nil, result
	}

	return &order, nil
}

func (g GormRepository) FindByItemId(ctx context.Context, itemId string) (*entity.OrderItem, error) {
	item := entity.OrderItem{}

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.OrderItem{}).
			First(&item, "id = ?", itemId).Error
	}, ctx)
	if result != nil {
		return nil, result
	}

	return &item, nil
}

func (g GormRepository) FindByUUID(ctx context.Context, uuid string) (*entity.Order, error) {
	order := entity.Order{}

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.Order{}).
			Preload("OrderItem").
			Preload("Delivery").
			Preload("OrderStatusLog").
			Preload("OrderCommission").
			First(&order, "order_uuid = ?", uuid).Error
	}, ctx)
	if result != nil {
		return nil, result
	}

	return &order, nil
}

func (g GormRepository) FindAll(ctx context.Context, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	whereState := query.ORMConditions().(string)

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.Order{}).
			Preload("OrderItem").
			Preload("Delivery").
			Where(whereState).
			Limit(query.GetLimit()).Offset(query.GetOffset()).
			Find(&orders).Error
	}, ctx)

	if result != nil {
		return nil, result
	}

	return orders, nil
}

func (g GormRepository) FindByUserId(ctx context.Context, userId string, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order

	whereState := query.UserORMConditions().(string)

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.Order{}).
			Preload("Delivery").
			Where(whereState).
			Where("orders.user_id", userId).
			Order("created_at desc").
			Limit(query.GetLimit()).Offset(query.GetOffset()).
			Find(&orders).Error
	}, ctx)

	if result != nil {
		return nil, result
	}

	return orders, nil
}

func (g GormRepository) FindOrderByStoreID(ctx context.Context, storeId string, query *pagable.Query, keyword string) ([]entity.Order, error) {
	var orders []entity.Order
	whereState := query.UserORMConditions().(string)
	var likeState string

	if strings.Contains(whereState, "status") {
		whereState = strings.Replace(whereState, "status", "orders.status", 1)
	}

	if len(keyword) > 2 {
		likeState = fmt.Sprintf("orders.order_uuid like %v", fmt.Sprintf("'%%%v%%'", keyword))
	}

	err := g.client.DB().Model(&entity.Order{}).
		Preload("Delivery").
		Joins("inner join orders_commission on orders_commission.order_id = orders.id").
		Where("orders_commission.store_id=?", storeId).
		Where(whereState).
		Where(likeState).
		Limit(query.GetLimit()).Offset(query.GetOffset()).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, err
}

func (g GormRepository) SearchOrderByStoreID(ctx context.Context, storeId string, keyword string, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	err := g.client.DB().Model(&entity.Order{}).
		Preload("Delivery").
		Joins("inner join orders_commission on orders_commission.order_id = orders.id").
		Where("orders_commission.store_id=?", storeId).
		Where("orders.order_uuid like ?", fmt.Sprintf("'%%%v%%'", keyword)).
		Limit(query.GetLimit()).Offset(query.GetOffset()).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, err
}

func (g GormRepository) TotalSearchOrderByStoreID(ctx context.Context, storeId string, keyword string) (int, error) {
	var count int64

	err := g.client.DB().Select("*").Model(&entity.Order{}).
		Joins("inner join orders_commission on orders_commission.order_id = orders.id").
		Where("orders_commission.store_id=?", storeId).
		Where("orders.order_uuid like ?", fmt.Sprintf("'%%%v%%'", keyword)).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return int(count), err
}

func (g GormRepository) FindOrderByDelivery(ctx context.Context, deliID string, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	err := g.client.DB().Model(&entity.Order{}).Preload("Delivery").
		Joins("inner join delivery_orders ON orders.id = delivery_orders.order_id").
		Where("delivery_orders.delivery_id=?", deliID).
		Order("orders.created_at DESC").
		Limit(query.GetLimit()).Offset(query.GetOffset()).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, err
}

func (g GormRepository) FindOrderByUserAndProduct(ctx context.Context, userId string, productId string) ([]entity.Order, error) {
	var orders []entity.Order
	err := g.client.DB().Raw("select * from orders inner join order_items on orders.id = order_items.order_id "+
		"where orders.user_id= ? and order_items.product_id = ?", userId, productId).Scan(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, err
}

func (g GormRepository) FindOrderLogByOrderId(ctx context.Context, orderId int) ([]entity.OrderStatusLog, error) {
	var orderStatus []entity.OrderStatusLog
	result := g.client.DB().Model(&entity.OrderStatusLog{}).
		Where("order_id", orderId).
		Order("created_at desc").
		Find(&orderStatus).Error
	if result != nil {
		return nil, result
	}

	return orderStatus, nil
}

func (g GormRepository) Save(ctx context.Context, dao *entity.Order) error {
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.Order{}).Create(&dao).Error
	}, ctx)

	return result
}

func (g GormRepository) UpdateStatus(ctx context.Context, orderID int, status int, message ...string) error {

	updateLog := entity.OrderStatusLog{
		OrderID:      orderID,
		OrderType:    "orders",
		StatusChange: status,
	}

	if len(message) == 0 {
		switch status {
		case entity.ORDER_PENDING:
			updateLog.Message = "Đơn hàng đang xử lý bởi nhà bán hàng"
		case entity.ORDER_DELIVERY:
			updateLog.Message = "Đơn hàng đang được vận chuyển"
		case entity.ORDER_SHIPPING_FINISH:
			updateLog.Message = "Đơn hàng được giao thành công"
		case entity.ORDER_CANCEL:
			updateLog.Message = "Đơn hàng bị hủy"
		}
	} else {
		updateLog.Message = message[0]
	}

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Transaction(func(tx *gormF.DB) error {
			if err := tx.Model(&entity.Order{}).
				Where("id = ?", orderID).Update("status", status).Error; err != nil {
				return err
			}

			if err := tx.Create(&updateLog).Error; err != nil {
				return err
			}
			return nil
		})
	}, ctx)

	if result != nil {
		return result
	}

	return nil
}

func (g GormRepository) Update(ctx context.Context, order entity.Order) error {
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Updates(order).Error
	}, ctx)

	if result.Error != nil {
		return result
	}
	return nil
}

func (g GormRepository) Total(ctx context.Context, query *pagable.Query) (int, error) {
	var count int64
	whereState := query.ORMConditions().(string)
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Select("*").Table(entity.Order{}.TableName()).
			Where(whereState).
			Count(&count).Error
	}, ctx)

	return int(count), result
}

func (g GormRepository) TotalStoreOrder(ctx context.Context, storeId string, query *pagable.Query, keyword string) (int, error) {
	var count int64
	var likeState string

	whereState := query.UserORMConditions().(string)
	if strings.Contains(whereState, "status") {
		whereState = strings.Replace(whereState, "status", "orders.status", 1)
	}

	if len(keyword) > 2 {
		likeState = fmt.Sprintf("orders.order_uuid like %v", fmt.Sprintf("'%%%v%%'", keyword))
	}

	err := g.client.DB().Select("*").Model(&entity.Order{}).
		Joins("inner join orders_commission on orders_commission.order_id = orders.id").
		Where("orders_commission.store_id=?", storeId).
		Where(whereState).
		Where(likeState).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (g GormRepository) UpdateOrderItem(ctx context.Context, orderItemID string, status int) error {
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.OrderItem{}).
			Where("id = ?", orderItemID).Update("status", status).Error
	}, ctx)

	if result != nil {
		return result
	}
	return nil
}
