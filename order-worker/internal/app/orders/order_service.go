package orders

import (
	"context"
	"fmt"
	"log"
	dto "order-worker/internal/domain/dto/order"
	"order-worker/internal/domain/entities/order"
	"order-worker/internal/infrastructure/adapter/productserv"
	"order-worker/pkg/util/mapper"
)

type orderService struct {
	orderRepo   order.Repository
	productServ productserv.Service
}

func NewOrderService(orderRepo order.Repository, productServ productserv.Service) Usecase {
	return orderService{
		orderRepo:   orderRepo,
		productServ: productServ,
	}
}

func (o orderService) RollBackQuantity(ctx context.Context, dto *dto.OrderMessage) error {
	return nil
}

func (o orderService) CreateOrder(ctx context.Context, message *dto.OrderMessage) error {
	orderDAO := order.Order{}
	orderDAO.OrderUUID = message.OrderUUID
	orderDAO.Username = message.UserRequest.Username
	orderDAO.UserId = message.UserRequest.UserId

	if err := mapper.BindingStruct(message, &orderDAO); err != nil {
		log.Printf("[%s] Mapping value from dto to dao failed cause: %s", "ERROR", err)
		return err
	}

	//create items in order
	var orderItems []*order.OrderItem
	for _, item := range message.OrderItems {
		i := order.OrderItem{
			ProductID:   item.ProductItem.ProductID,
			ProductName: item.ProductItem.ProductName,
			StoreID:     item.ProductItem.StoreID,
			OptionID:    item.ProductItem.OptionID,
			Quantity:    item.ProductItem.Quantity,
			Price:       item.ProductItem.Price,
			NetPrice:    item.ProductItem.NetPrice,
			Order:       &orderDAO,
		}
		orderItems = append(orderItems, &i)
	}
	orderDAO.OrderItem = orderItems

	//calculate order price
	orderDAO.SubTotal = message.SubTotal
	orderDAO.Amount = message.Amount
	orderDAO.Discount = message.Discount

	//create log
	var logs []*order.OrderStatusLog
	orderLog := order.OrderStatusLog{
		Order:        &orderDAO,
		Message:      "order created",
		StatusChange: order.ORDER_CREATED,
	}
	orderDAO.OrderStatusLog = append(logs, &orderLog)

	//create delivery
	recvTime, err := order.ParseStringToDate(message.Delivery.ReceivingDate)
	if err != nil {
		return err
	}
	deli := order.DeliveryOrder{
		DeliveryId:      message.Delivery.DeliveryId,
		DeliveryName:    message.Delivery.Name,
		Cost:            message.Delivery.Cost,
		AddressId:       message.Address.AddressId,
		ShippingName:    message.Address.ShippingName,
		ShippingPhone:   message.Address.ShippingPhone,
		ShippingAddress: message.Address.ShippingAddress,
		ReceivingDate:   *recvTime,
		Order:           &orderDAO,
	}
	orderDAO.Delivery = &deli

	vouchers := ""
	for _, i := range message.Vouchers {
		vouchers += fmt.Sprintf("%v;", i)
	}
	orderDAO.VoucherCode = vouchers

	err = o.orderRepo.Save(&orderDAO)
	if err != nil {
		//handle rollback
		log.Printf("[%s] The order created failed : %s", "ERROR", err)
		return err
	}

	return nil
}
