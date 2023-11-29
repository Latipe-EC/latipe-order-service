package orders

import (
	"context"
	"fmt"
	log2 "github.com/gofiber/fiber/v2/log"
	"log"
	messageDTO "order-worker/internal/domain/dto"
	dto "order-worker/internal/domain/dto/order"
	"order-worker/internal/domain/entities/order"
	"order-worker/internal/infrastructure/adapter/productserv"
	"order-worker/internal/infrastructure/adapter/storeserv"
	storeDTO "order-worker/internal/infrastructure/adapter/storeserv/dto"
	"order-worker/internal/publisher"
	"order-worker/pkg/util/mapper"
)

type orderService struct {
	orderRepo   order.Repository
	productServ productserv.Service
	storeServ   storeserv.Service
	message     *publisher.MessageProducer
}

func NewOrderService(orderRepo order.Repository, productServ productserv.Service, storeServ storeserv.Service,
	message *publisher.MessageProducer) Usecase {
	return orderService{
		orderRepo:   orderRepo,
		productServ: productServ,
		storeServ:   storeServ,
		message:     message,
	}
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
			ProdImg:     item.ProductItem.Image,
			Order:       &orderDAO,
		}
		if i.NetPrice != 0 {
			i.SubTotal = i.NetPrice * i.Quantity
		} else {
			i.SubTotal = i.Price * i.Quantity
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
		Message:      "Đơn hàng được tạo thành công",
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
	orderDAO.PaymentMethod = message.PaymentMethod

	orderDAO.Status = order.ORDER_CREATED
	err = o.orderRepo.Save(&orderDAO)
	if err != nil {
		//handle rollback
		log.Printf("[%s] The order created failed : %s", "ERROR", err)
		return err
	}

	//send notify email
	email := messageDTO.EmailRequest{
		EmailRecipient: message.UserRequest.Username,
		Name:           message.Address.ShippingName,
		OrderId:        message.OrderUUID,
		Code:           message.OrderUUID,
	}

	err = o.message.SendEmailMessage(email)
	if err != nil {
		log.Printf("[%s] sending message to email queue was failed : %s", "ERROR", err)
		return err
	}

	//delete cart
	var cartItemList []string
	for _, i := range message.OrderItems {
		if i.CartId != "" {
			cartItemList = append(cartItemList, i.CartId)
		}
	}
	if len(cartItemList) > 0 {
		msg := messageDTO.CartMessage{CartIdVmList: cartItemList}

		err := o.message.SendCartServiceMessage(&msg)
		if err != nil {
			log.Printf("[%s] sending message to cart queue was failed : %s", "ERROR", err)
			return err
		}
	}
	return nil
}

func (o orderService) CreateCommissionOrderComplete(ctx context.Context) error {
	idStr := ""
	rows := 0
	orders, err := o.orderRepo.FindAllFinishShippingOrder()
	if err != nil {
		return err
	}

	if len(orders) > 0 {
		//loop data in orders
		for _, i := range orders {
			//check time update after seven day
			if IsAfterSevenDays(i.UpdatedAt) {
				//handle commission
				if err := o.createCommissionOfOrder(ctx, &i); err != nil {
					break
				}

				idStr += fmt.Sprintf("%v;", i.Id)
				rows++
			}
		}
	}
	log2.Infof("total rows was update [%v]", rows)
	log2.Infof("order was update [%v]", idStr)
	return nil
}

func (o orderService) createCommissionOfOrder(ctx context.Context, dao *order.Order) error {
	amountFromStore, err := o.orderRepo.GetOrderAmountOfStore(dao.Id)
	if err != nil {
		return err
	}

	for _, i := range amountFromStore {
		req := storeDTO.GetStoreByIdRequest{
			StoreID: i.StoreId,
		}

		storeCms, err := o.storeServ.GetStoreByStoreId(ctx, &req)
		if err != nil {
			return err
		}

		systemFee := int(float64(i.OrderAmount) * storeCms.FeePerOrder)
		storeReceived := i.OrderAmount - systemFee
		oc := order.OrderCommission{
			OrderID:        dao.Id,
			StoreID:        storeCms.Id,
			AmountReceived: storeReceived,
			SystemFee:      systemFee,
		}

		dao.Status = order.ORDER_COMPLETED

		orderStatusLog := order.OrderStatusLog{
			OrderID:      dao.Id,
			Message:      "Đơn hàng hoàn thành",
			StatusChange: order.ORDER_COMPLETED,
		}

		if err := o.orderRepo.CreateOrderCommmsionTransaction(dao, &oc, &orderStatusLog); err != nil {
			return err
		}

		msg := messageDTO.StoreBillingMessage{
			StoreID:        oc.StoreID,
			OrderUUID:      dao.OrderUUID,
			AmountReceived: oc.AmountReceived,
			SystemFee:      oc.SystemFee,
		}

		if err := o.message.SendBillingServiceMessage(&msg); err != nil {
			return err
		}

	}

	return nil
}
