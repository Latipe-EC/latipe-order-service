package orders

import (
	"context"
	dto "order-service-rest-api/internal/domain/dto/order"
	"order-service-rest-api/internal/domain/entities/order"
	"order-service-rest-api/internal/message"
)

type OrderService struct {
	orderRepo       order.Repository
	producerMessage message.ProducerOrderMessage
}

func (o OrderService) CreateOrder(ctx context.Context, dto *dto.CreateOrderRequest) error {
	//TODO implement me
	panic("implement me")
}

func (o OrderService) GetOrderById(ctx context.Context, dto *dto.GetOrderRequest) (*dto.GetOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderService) GetOrderList(ctx context.Context, dto *dto.GetOrderListRequest) (*dto.GetOrderListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderService) UpdateStatusOrder(ctx context.Context, dto *dto.UpdateOrderStatusRequest) error {
	//TODO implement me
	panic("implement me")
}

func (o OrderService) UpdateOrder(ctx context.Context, dto *dto.UpdateOrderRequest) error {
	//TODO implement me
	panic("implement me")
}
