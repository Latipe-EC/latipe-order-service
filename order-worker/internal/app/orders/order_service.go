package orders

import (
	"context"
	"encoding/json"
	"log"
	dto "order-worker/internal/domain/dto/order"
	"order-worker/internal/domain/entities/order"
	"order-worker/internal/infrastructure/adapter/productserv"
	prodServDTO "order-worker/internal/infrastructure/adapter/productserv/dto"
	"order-worker/internal/message"
	"order-worker/pkg/cache/redis"
	"order-worker/pkg/util/mapper"
)

type orderService struct {
	orderRepo   order.Repository
	cacheRepo   *redis.CacheEngine
	productServ productserv.Service
}

func NewOrderService(orderRepo order.Repository, productServ productserv.Service, cacheRepo *redis.CacheEngine) Usecase {
	return orderService{
		orderRepo:   orderRepo,
		productServ: productServ,
		cacheRepo:   cacheRepo,
	}
}

func (o orderService) RollBackQuantity(ctx context.Context, dto *dto.CreateOrderRequest) error {
	reduceReq := prodServDTO.ReduceProductRequest{
		Items: MappingOrderItemForReduce(dto),
	}

	if _, err := o.productServ.ReduceProductQuantity(&reduceReq); err != nil {
		return err
	}

	return nil
}

func (o orderService) CreateOrder(ctx context.Context, orderCacheKey string) error {
	cacheData, err := o.cacheRepo.Get(orderCacheKey)
	if err != nil {
		return err
	}

	dto := new(dto.CreateOrderRequest)
	if cacheData != nil {
		if err := json.Unmarshal(cacheData, &dto); err != nil {
			return err
		}
	}

	orderDAO := order.Order{}
	orderDAO.OrderCacheKey = orderCacheKey
	orderDAO.Username = dto.UserRequest.Username
	orderDAO.UserId = dto.UserRequest.UserId

	if err := mapper.BindingStruct(dto, &orderDAO); err != nil {
		log.Printf("[%s] Mapping value from dto to dao failed cause: %s", "ERROR", err)
		return err
	}

	itemReq := MappingOrderItemToValidateItems(dto.OrderItems)
	productReq := prodServDTO.OrderProductRequest{Items: itemReq}

	products, err := o.productServ.GetProductOrderInfo(&productReq)
	if err != nil {
		return err
	}

	//create items in order
	total := 0
	var orderItems []*order.OrderItem
	for _, item := range products.Products {
		total += item.Price
		i := order.OrderItem{
			ProductID:   item.ProductId,
			ProductName: item.Name,
			StoreID:     item.StoreId,
			OptionID:    item.OptionId,
			Quantity:    item.Quantity,
			Price:       item.Price,
			NetPrice:    item.PromotionalPrice,
			Order:       &orderDAO,
		}
		orderItems = append(orderItems, &i)
	}
	orderDAO.OrderItem = orderItems

	//calculate order price
	orderDAO.Total = total
	orderDAO.Amount = products.TotalPrice
	orderDAO.Discount = total - products.TotalPrice

	//create log
	var logs []*order.OrderStatusLog
	orderLog := order.OrderStatusLog{
		Order:        &orderDAO,
		Message:      "order created",
		StatusChange: order.ORDER_CREATED,
	}
	orderDAO.OrderStatusLog = append(logs, &orderLog)

	//create delivery
	recvTime, err := order.ParseStringToDate(dto.Delivery.ReceivingDate)
	if err != nil {
		return err
	}
	deli := order.DeliveryOrder{
		DeliveryId:      dto.Delivery.DeliveryId,
		DeliveryName:    dto.Delivery.Name,
		Cost:            dto.Delivery.Cost,
		AddressId:       dto.Address.AddressId,
		ShippingName:    dto.Address.ShippingName,
		ShippingPhone:   dto.Address.ShippingPhone,
		ShippingAddress: dto.Address.ShippingAddress,
		ReceivingDate:   *recvTime,
		Order:           &orderDAO,
	}
	orderDAO.Delivery = &deli

	err = o.orderRepo.Save(&orderDAO)
	if err != nil {
		reduceReq := prodServDTO.RollbackQuantityRequest{
			Items: MappingOrderItemForRollback(dto),
		}

		if _, err := o.productServ.RollBackQuantityOrder(&reduceReq); err != nil {
			log.Printf("[%s] Rollback product quantity was failed : %v", "ERROR", err)
			return err
		}

		log.Printf("[%s] The order created failed : %s", "ERROR", err)
		return err
	}

	//handle cache value
	switch dto.PaymentMethod {
	case order.PAYMENT_COD:
		if err := o.cacheRepo.Delete(orderCacheKey); err != nil {
			return err
		}
	case order.PAYMENT_VIA_PAYPAL:
		if err := o.cacheRepo.Expire(orderCacheKey); err != nil {
			return err
		}
	}

	//reset cart-items
	var cartItemList []string
	for _, i := range dto.OrderItems {
		if i.CartItemId != "" {
			cartItemList = append(cartItemList, i.CartItemId)
		}
	}

	if len(cartItemList) > 0 {
		err := message.SendCartServiceMessage(cartItemList)
		if err != nil {
			log.Printf("[%s] sending message to cart queue was failed : %s", "ERROR", err)
			return err
		}
	}

	return nil
}
