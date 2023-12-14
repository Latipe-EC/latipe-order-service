package order

import (
	"context"
	gormF "gorm.io/gorm"
	"order-rest-api/internal/domain/dto/custom_entity"
	entity "order-rest-api/internal/domain/entities/order"
	"order-rest-api/internal/domain/msg"
)

func (g GormRepository) UserCountingOrder(ctx context.Context, userId string) (int, error) {
	var count int64
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.Order{}).
			Where("user_id=?", userId).
			Count(&count).Error
	}, ctx)

	return int(count), result
}

func (g GormRepository) StoreCountingOrder(ctx context.Context, storeId string) (int, error) {
	queryResp := struct {
		Count int64 `json:"count"`
	}{}
	sql := `SELECT count(distinct order_id) as count
			FROM orders
    		Join order_items oi on orders.id = oi.order_id
			Where store_id = ?
			`
	err := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Raw(sql, storeId).Scan(&queryResp).Error
	}, ctx)
	if err != nil {
		return 0, err
	}

	return int(queryResp.Count), err

}

func (g GormRepository) DeliveryCountingOrder(ctx context.Context, deliveryId string) (int, error) {
	var count int64
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&msg.Delivery{}).
			Where("delivery_id=?", deliveryId).
			Count(&count).Error
	}, ctx)

	return int(count), result
}

func (g GormRepository) AdminCountingOrder(ctx context.Context) (int, error) {
	var count int64
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Select("*").Table(entity.Order{}.TableName()).
			Count(&count).Error
	}, ctx)

	return int(count), result
}

func (g GormRepository) GetTotalOrderInSystemInDay(ctx context.Context, date string) ([]custom_entity.TotalOrderInSystemInHours, error) {
	var result []custom_entity.TotalOrderInSystemInHours

	err := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Table("orders").
			Select("HOUR(orders.created_at) as hour, SUM(amount) as amount, COUNT(*) as count").
			Where("date(orders.created_at) = (?)", date).
			Group("HOUR(orders.created_at)").
			Order("HOUR(orders.created_at) DESC").
			Scan(&result).Error
	}, ctx)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetTotalOrderInSystemInMonth(ctx context.Context, date string) ([]custom_entity.TotalOrderInSystemInDay, error) {
	var result []custom_entity.TotalOrderInSystemInDay

	err := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Table("orders").
			Select("DAY(orders.created_at) as day, SUM(amount) as amount, COUNT(*) as count").
			Where("orders.created_at >= ?", date).
			Where("year(orders.created_at) = year(?)", date).
			Where("month(orders.created_at) = month(?)", date).
			Group("DAY(orders.created_at)").
			Order("DAY(orders.created_at) DESC").
			Scan(&result).Error
	}, ctx)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetTotalOrderInSystemInYear(ctx context.Context, year int) ([]custom_entity.TotalOrderInSystemInMonth, error) {
	var result []custom_entity.TotalOrderInSystemInMonth

	err := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Table("orders").
			Select("MONTH(orders.created_at) as month, SUM(amount) as amount, COUNT(*) as count").
			Where("year(orders.created_at) = ?", year).
			Group("MONTH(orders.created_at)").
			Order("MONTH(orders.created_at) DESC").
			Scan(&result).Error
	}, ctx)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetTotalCommissionOrderInYear(ctx context.Context, date string) ([]custom_entity.SystemOrderCommissionDetail, error) {
	var result []custom_entity.SystemOrderCommissionDetail

	err := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Table("orders").
			Select("MONTH(orders.created_at) as month, COUNT(*) as total_orders, "+
				"SUM(amount) as amount, "+
				"SUM(orders_commission.amount_received) as store_received, "+
				"SUM(orders_commission.system_fee) as system_received").
			Joins("INNER JOIN orders_commission ON orders.id = orders_commission.order_id").
			Where("orders.created_at >= ?", date).
			Where("year(orders_commission.created_at) = year(?)", date).
			Group("MONTH(orders.created_at)").
			Order("MONTH(orders.created_at) DESC").
			Scan(&result).Error
	}, ctx)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) TopOfProductSold(ctx context.Context, date string, count int) ([]custom_entity.TopOfProductSold, error) {
	var result []custom_entity.TopOfProductSold

	err := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Table("orders").
			Select("oi.product_id as product_id, oi.product_name as product_name, SUM(oi.quantity) as total").
			Joins("INNER JOIN order_items oi ON orders.id = oi.order_id").
			Where("orders.created_at >= ?", date).
			Where("year(orders.created_at) = year(?)", date).
			Where("month(orders.created_at) = month(?)", date).
			Group("oi.product_id, oi.product_name").
			Limit(count).
			Scan(&result).Error
	}, ctx)

	if err != nil {
		return nil, err
	}
	return result, err
}

func (g GormRepository) GetTotalOrderInSystemInMonthOfStore(ctx context.Context, date string, storeId string) ([]custom_entity.TotalOrderInSystemInDay, error) {
	var result []custom_entity.TotalOrderInSystemInDay

	err := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Table("orders_commission").
			Select("DAY(orders_commission.created_at) as day, SUM(amount_received) as amount, COUNT(*) as count").
			Where("orders_commission.store_id = ?", storeId).
			Where("orders_commission.created_at >= ?", date).
			Where("year(orders_commission.created_at) = year(?)", date).
			Where("month(orders_commission.created_at) = month(?)", date).
			Group("DAY(orders_commission.created_at)").
			Order("DAY(orders_commission.created_at) DESC").
			Scan(&result).Error
	}, ctx)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetTotalOrderInSystemInYearOfStore(ctx context.Context, year int, storeId string) ([]custom_entity.TotalOrderInSystemInMonth, error) {
	var result []custom_entity.TotalOrderInSystemInMonth

	err := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Table("orders_commission").
			Select("MONTH(orders_commission.created_at) as month, SUM(amount_received) as amount, COUNT(*) as count").
			Where("orders_commission.store_id = ?", storeId).
			Where("year(orders_commission.created_at) = ?", year).
			Group("MONTH(orders_commission.created_at)").
			Order("MONTH(orders_commission.created_at) DESC").
			Scan(&result).Error
	}, ctx)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetTotalCommissionOrderInYearOfStore(ctx context.Context, date string, storeId string) ([]custom_entity.OrderCommissionDetail, error) {
	var result []custom_entity.OrderCommissionDetail

	err := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Table("orders_commission").
			Select("MONTH(orders_commission.created_at) as month,"+
				"COUNT(*) as total_orders, "+
				"SUM(orders_commission.amount_received) as total_received, "+
				"SUM(orders_commission.system_fee) as total_fee").
			Where("orders_commission.store_id = ?", storeId).
			Where("orders_commission.status >= ?", 1).
			Where("orders_commission.created_at >= ?", date).
			Where("year(orders_commission.created_at) = year(?)", date).
			Group("MONTH(orders_commission.created_at)").
			Order("MONTH(orders_commission.created_at) DESC").
			Scan(&result).Error
	}, ctx)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) TopOfProductSoldOfStore(ctx context.Context, date string, count int, storeId string) ([]custom_entity.TopOfProductSold, error) {
	var result []custom_entity.TopOfProductSold

	err := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Table("orders").
			Select("order_items.product_id as product_id, order_items.product_name as product_name, "+
				"SUM(order_items.quantity) as total").
			Joins("INNER JOIN order_items ON orders.id = order_items.order_id").
			Where("order_items.store_id = ?", storeId).
			Where("orders.created_at >= ?", date).
			Group("order_items.product_id, order_items.product_name").
			Limit(count).
			Scan(&result).Error
	}, ctx)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (g GormRepository) GetOrderAmountOfStore(ctx context.Context, orderId int) ([]custom_entity.AmountItemOfStoreInOrder, error) {
	var result []custom_entity.AmountItemOfStoreInOrder

	err := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Table("orders").
			Select("oi.store_id as store_id, SUM(order_items.sub_total) as order_amount").
			Joins("INNER JOIN order_items ON orders.id = order_items.order_id").
			Where("orders.id = ?", orderId).
			Group("order_items.store_id").
			Scan(&result).Error
	}, ctx)
	if err != nil {
		return nil, err
	}

	return result, err
}
