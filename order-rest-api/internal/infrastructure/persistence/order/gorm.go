package order

import (
	"fmt"
	gormF "gorm.io/gorm"
	"order-rest-api/internal/domain/dto/custom_entity"
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
		&entity.OrderCommission{},
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
		Preload("Delivery").
		First(&order, Id).Error
	if result != nil {
		return nil, result
	}

	return &order, nil
}

func (g GormRepository) FindByUUID(uuid string) (*entity.Order, error) {
	order := entity.Order{}

	result := g.client.DB().Model(&entity.Order{}).
		Preload("OrderItem").
		Preload("Delivery").
		First(&order, "order_uuid = ?", uuid).Error
	if result != nil {
		return nil, result
	}

	return &order, nil
}

func (g GormRepository) FindAll(query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	whereState := query.ORMConditions().(string)
	result := g.client.DB().Model(&entity.Order{}).
		Preload("OrderItem").
		Preload("Delivery").
		Where(whereState).
		Limit(query.GetLimit()).Offset(query.GetOffset()).
		Find(&orders).Error
	if result != nil {
		return nil, result
	}

	return orders, nil
}

func (g GormRepository) FindByUserId(userId string, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	result := g.client.DB().Model(&entity.Order{}).
		Preload("Delivery").
		Where("orders.user_id", userId).
		Order("created_at desc").
		Limit(query.GetLimit()).Offset(query.GetOffset()).
		Find(&orders).Error

	if result != nil {
		return nil, result
	}

	return orders, nil
}

func (g GormRepository) FindOrderByStoreID(storeId string, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	sql := `
	SELECT orders.id,orders.created_at,
       order_uuid,username,amount,payment_method,orders.status,orders.created_at
	FROM orders
	join order_items oi on orders.id = oi.order_id
	where store_id=?
	group by orders.id
	order by orders.created_at desc
	limit ?,?
	`
	err := g.client.DB().Raw(sql, storeId, query.Page, query.Size).Scan(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, err
}

func (g GormRepository) FindOrderByDelivery(deliID string, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	err := g.client.DB().Model(&entity.Order{}).Preload("Delivery").
		Joins("inner join delivery_orders ON orders.id = delivery_orders.order_id").
		Where("delivery_orders.id=?", deliID).
		Order("orders.created_at DESC").
		Limit(query.GetLimit()).Offset(query.GetOffset()).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, err
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
		Order("created_at desc").
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

func (g GormRepository) UpdateStatus(orderID int, status int) error {

	updateLog := entity.OrderStatusLog{
		OrderID:      orderID,
		Message:      "",
		StatusChange: status,
	}

	result := g.client.Transaction(func(tx *gormF.DB) error {
		if err := tx.Model(&entity.Order{}).
			Where("id = ?", orderID).Update("status", status).Error; err != nil {
			return err
		}

		if err := tx.Create(&updateLog).Error; err != nil {
			return err
		}

		return nil
	})

	if result != nil {
		return result
	}

	return nil
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

func (g GormRepository) UpdateOrderItem(orderItemID int, status int) error {
	result := g.client.DB().Model(&entity.OrderItem{}).
		Where("id = ?", orderItemID).Update("status", status)

	if result.Error != nil || result.RowsAffected == 0 {
		return result.Error
	}
	return nil
}

func (g GormRepository) GetTotalOrderInSystemInDay(date string) ([]custom_entity.TotalOrderInSystemInHours, error) {
	var result []custom_entity.TotalOrderInSystemInHours

	err := g.client.DB().Table("orders").
		Select("HOUR(created_at) as hour, SUM(amount) as amount, COUNT(*) as count").
		Where("created_at > ?", date).
		Group("HOUR(created_at)").
		Order("HOUR(created_at) DESC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetTotalOrderInSystemInMonth(month int, year int) ([]custom_entity.TotalOrderInSystemInDay, error) {
	var result []custom_entity.TotalOrderInSystemInDay

	err := g.client.DB().Table("orders").
		Select("DAY(created_at) as day, SUM(amount) as amount, COUNT(*) as count").
		Where("created_at > ?", fmt.Sprintf("%v-%v-01", year, month)).
		Group("DAY(created_at)").
		Order("DAY(created_at) DESC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetTotalOrderInSystemInYear(year int) ([]custom_entity.TotalOrderInSystemInMonth, error) {
	var result []custom_entity.TotalOrderInSystemInMonth

	err := g.client.DB().Table("orders").
		Select("YEAR(created_at) as month, SUM(amount) as amount, COUNT(*) as count").
		Where("created_at > ?", fmt.Sprintf("%v-01-01", year)).
		Group("YEAR(created_at)").
		Order("YEAR(created_at) DESC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetTotalCommissionOrderInYear(month int, year int, count int) ([]custom_entity.OrderCommissionDetail, error) {
	var result []custom_entity.OrderCommissionDetail

	err := g.client.DB().Table("orders").
		Select("MONTH(orders.created_at) as month, COUNT(*) as total_orders, "+
			"SUM(amount) as amount, "+
			"SUM(orders_commission.amount_received) as total_store_received, "+
			"SUM(orders_commission.system_fee) as total_fee").
		Joins("INNER JOIN orders_commission ON orders.id = orders_commission.order_id").
		Where("where orders.created_at >= ?", fmt.Sprintf("%v-%v-01", year, month)).
		Group("MONTH(orders.created_at)").
		Order("MONTH(orders.created_at) DESC").
		Limit(count).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) ListOfProductSelledOnMonth(month int, year int, count int) ([]custom_entity.TopOfProductSold, error) {
	var result []custom_entity.TopOfProductSold

	err := g.client.DB().Table("orders").
		Select("oi.product_id as product_id, oi.product_name as product_name, SUM(oi.quantity) as total").
		Joins("INNER JOIN order_items oi ON orders.id = oi.order_id").
		Where("orders.created_at >= ?", fmt.Sprintf("%v-%v-01", year, month)).
		Group("oi.product_id, oi.product_name").
		Limit(count).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetTotalOrderInSystemInMonthOfStore(month int, year int, storeId string) ([]custom_entity.TotalOrderInSystemInDay, error) {
	var result []custom_entity.TotalOrderInSystemInDay

	err := g.client.DB().Table("orders").
		Select("DAY(created_at) as day, SUM(amount) as amount, COUNT(*) as count").
		Joins("INNER JOIN order_items ON orders.id = order_items.order_id").
		Where("order_items.store_id = ?", storeId).
		Where("created_at > ?", fmt.Sprintf("%v-%v-01", year, month)).
		Group("DAY(created_at)").
		Order("DAY(created_at) DESC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetTotalOrderInSystemInYearOfStore(year int, storeId string) ([]custom_entity.TotalOrderInSystemInMonth, error) {
	var result []custom_entity.TotalOrderInSystemInMonth

	err := g.client.DB().Table("orders").
		Select("YEAR(created_at) as month, SUM(amount) as amount, COUNT(*) as count").
		Joins("INNER JOIN order_items ON orders.id = order_items.order_id").
		Where("order_items.store_id = ?", storeId).
		Where("created_at > ?", fmt.Sprintf("%v-01-01", year)).
		Group("YEAR(created_at)").
		Order("YEAR(created_at) DESC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetTotalCommissionOrderInYearOfStore(month int, year int, count int, storeId string) ([]custom_entity.OrderCommissionDetail, error) {
	var result []custom_entity.OrderCommissionDetail

	err := g.client.DB().Table("orders").
		Select("MONTH(orders.created_at) as month, COUNT(*) as total_orders, "+
			"SUM(amount) as amount, "+
			"SUM(orders_commission.amount_received) as total_store_received, "+
			"SUM(orders_commission.system_fee) as total_fee").
		Joins("INNER JOIN orders_commission ON orders.id = orders_commission.order_id").
		Where("orders_commission.store_id = ?", storeId).
		Where("where orders.created_at >= ?", fmt.Sprintf("%v-%v-01", year, month)).
		Group("MONTH(orders.created_at)").
		Order("MONTH(orders.created_at) DESC").
		Limit(count).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) ListOfProductSelledOnMonthStore(month int, year int, count int, storeId string) ([]custom_entity.TopOfProductSold, error) {
	var result []custom_entity.TopOfProductSold

	err := g.client.DB().Table("orders").
		Select("oi.product_id as product_id, oi.product_name as product_name, SUM(oi.quantity) as total").
		Joins("INNER JOIN order_items oi ON orders.id = oi.order_id").
		Where("order_items.store_id = ?", storeId).
		Where("orders.created_at >= ?", fmt.Sprintf("%v-%v-01", year, month)).
		Group("oi.product_id, oi.product_name").
		Limit(count).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetOrderAmountOfStore(orderId int) ([]custom_entity.AmountItemOfStoreInOrder, error) {
	var result []custom_entity.AmountItemOfStoreInOrder

	err := g.client.DB().Table("orders").
		Select("oi.store_id as store_id, SUM(order_items.sub_total) as order_amount").
		Joins("INNER JOIN order_items ON orders.id = order_items.order_id").
		Where("orders.id = ?", orderId).
		Group("order_items.store_id").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, err
}
