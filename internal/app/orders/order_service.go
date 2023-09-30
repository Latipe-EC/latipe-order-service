package orders

import (
	"context"
	"order-service-rest-api/internal/common/errors"
	dto "order-service-rest-api/internal/domain/dto/order"
	"order-service-rest-api/internal/domain/entities/order"
	"order-service-rest-api/internal/infrastructure/adapter/productserv"
	prodServDTO "order-service-rest-api/internal/infrastructure/adapter/productserv/dto"
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
		return errors.NotAvailableQuantity
	}

	return nil

}

func (o orderService) RollBackQuantity(ctx context.Context, dto *dto.CreateOrderRequest) error {
	reduceReq := prodServDTO.ReduceProductRequest{
		Items: MappingOrderItemForReduce(dto),
	}

	if _, err := o.productServ.ReduceProductQuantity(ctx, &reduceReq); err != nil {
		return errors.NotAvailableQuantity
	}

	return nil
}

func (o orderService) CreateOrder(ctx context.Context, dto *dto.CreateOrderRequest) error {
	//TODO implement me
	panic("implement me")
}

func (o orderService) GetOrderById(ctx context.Context, dto *dto.GetOrderRequest) (*dto.GetOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) GetOrderList(ctx context.Context, dto *dto.GetOrderListRequest) (*dto.GetOrderListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) UpdateStatusOrder(ctx context.Context, dto *dto.UpdateOrderStatusRequest) error {
	//TODO implement me
	panic("implement me")
}

func (o orderService) UpdateOrder(ctx context.Context, dto *dto.UpdateOrderRequest) error {
	//TODO implement me
	panic("implement me")
}
