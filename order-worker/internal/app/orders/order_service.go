package orders

import (
	"context"
	"github.com/google/uuid"
	"log"
	dto "order-worker/internal/domain/dto/order"
	"order-worker/internal/domain/entities/order"
	"order-worker/internal/infrastructure/adapter/productserv"
	prodServDTO "order-worker/internal/infrastructure/adapter/productserv/dto"
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

func (o orderService) PendingOrder(ctx context.Context, dto *dto.CreateOrderRequest) error {
	reduceReq := prodServDTO.ReduceProductRequest{
		Items: MappingOrderItemForReduce(dto),
	}

	if _, err := o.productServ.ReduceProductQuantity(ctx, &reduceReq); err != nil {
		return err
	}

	return nil

}

func (o orderService) RollBackQuantity(ctx context.Context, dto *dto.CreateOrderRequest) error {
	reduceReq := prodServDTO.ReduceProductRequest{
		Items: MappingOrderItemForReduce(dto),
	}

	if _, err := o.productServ.ReduceProductQuantity(ctx, &reduceReq); err != nil {
		return err
	}

	return nil
}

func (o orderService) CreateOrder(ctx context.Context, dto *dto.CreateOrderRequest) error {
	orderDAO := order.Order{}

	if err := mapper.BindingStruct(dto, &orderDAO); err != nil {
		log.Printf("[%s] Mapping value from dto to dao failed cause: %s", "ERROR", err)
		return err
	}

	//create items
	var orderItems []*order.OrderItem
	for _, item := range dto.OrderItems {
		i := order.OrderItem{
			ProductID: item.ProductId,
			SellerID:  0,
			OptionID:  item.OptionId,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Order:     &orderDAO,
		}
		orderItems = append(orderItems, &i)
	}
	orderDAO.OrderItem = orderItems

	//create log
	var logs []*order.OrderStatusLog
	orderLog := order.OrderStatusLog{
		Order:        &orderDAO,
		Message:      "order created",
		StatusChange: order.ORDER_CREATED,
	}
	orderDAO.OrderStatusLog = append(logs, &orderLog)

	//create payment
	paymentLog := order.PaymentLog{
		PaymentTransaction: uuid.New().String(),
		PaymentType:        dto.PaymentMethod,
		Total:              0,
		ThirdPartyLog:      "",
	}
	orderDAO.PaymentLog = &paymentLog
	orderDAO.Username = dto.UserRequest.Username
	orderDAO.UserId = dto.UserRequest.UserId

	err := o.orderRepo.Save(&orderDAO)
	if err != nil {
		reduceReq := prodServDTO.RollbackQuantityRequest{
			Items: MappingOrderItemForRollback(dto),
		}

		if _, err := o.productServ.RollBackQuantityOrder(ctx, &reduceReq); err != nil {
			log.Printf("[%s] Rollback product quantity was failed : %v", "ERROR", err)
			return err
		}

		log.Printf("[%s] The order created failed : %s", "ERROR")
		return err
	}
	return nil
}
