package orders

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"order-rest-api/internal/common/errors"
	orderDTO "order-rest-api/internal/domain/dto/order"
	"order-rest-api/internal/domain/entities/order"
	"order-rest-api/internal/infrastructure/adapter/productserv"
	prodServDTO "order-rest-api/internal/infrastructure/adapter/productserv/dto"
	"order-rest-api/internal/middleware/auth"
	"order-rest-api/pkg/cache/redis"
	"order-rest-api/pkg/util/mapper"
)

type orderService struct {
	orderRepo   order.Repository
	cacheEngine *redis.CacheEngine
	productServ productserv.Service
}

func NewOrderService(orderRepo order.Repository, productServ productserv.Service, cacheEngine *redis.CacheEngine) Usecase {
	return orderService{
		orderRepo:   orderRepo,
		cacheEngine: cacheEngine,
		productServ: productServ,
	}
}

func (o orderService) ProcessCacheOrder(ctx context.Context, dto *orderDTO.CreateOrderRequest) (string, error) {
	reduceReq := prodServDTO.ReduceProductRequest{
		Items: MappingOrderItemForReduce(dto),
	}

	if _, err := o.productServ.ReduceProductQuantity(ctx, &reduceReq); err != nil {
		return "", errors.NotAvailableQuantity
	}

	tempId := uuid.NewString()

	cacheData, err := json.Marshal(dto)
	if err != nil {
		return "", err
	}

	if err := o.cacheEngine.Set(tempId, cacheData); err != nil {
		return "", err
	}
	return tempId, nil
}

func (o orderService) GetOrderById(ctx context.Context, dto *orderDTO.GetOrderRequest) (*orderDTO.GetOrderResponse, error) {
	orderResp := orderDTO.OrderResponse{}

	orderDAO, err := o.orderRepo.FindById(dto.OrderId)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orderDAO, &orderResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderResponse{Order: orderResp}

	return &resp, err
}

func (o orderService) GetOrderList(ctx context.Context, dto *orderDTO.GetOrderListRequest) (*orderDTO.GetOrderListResponse, error) {
	var dataResp []orderDTO.OrderResponse

	orders, err := o.orderRepo.FindAll(dto.Query)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.Total(dto.Query)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderListResponse{}
	resp.Data = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderService) GetOrderByUserId(ctx context.Context, dto *orderDTO.GetByUserIdRequest) (*orderDTO.GetByUserIdResponse, error) {
	dataResp := orderDTO.OrderResponse{}

	orders, err := o.orderRepo.FindByUserId(dto.UserId, dto.Query)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.Total(dto.Query)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetByUserIdResponse{}
	resp.Data = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderService) UpdateStatusOrder(ctx context.Context, dto *orderDTO.UpdateOrderStatusRequest) error {

	orderDAO, err := o.orderRepo.FindById(dto.OrderId)
	if err != nil {
		return err
	}

	if orderDAO.Status == order.ORDER_CANCEL {
		return errors.ErrBadRequest
	}

	switch dto.Role {
	case auth.ROLE_USER:
		if dto.Status != order.ORDER_CANCEL || orderDAO.Status != order.ORDER_CREATED {
			return errors.ErrPermissionDenied
		}
	case auth.ROLE_STORE:
		if dto.Status != order.ORDER_PENDING && dto.Status != order.ORDER_DELIVERY {
			return errors.ErrPermissionDenied
		}
	case auth.ROLE_SHIPPER:
		if orderDAO.Status != order.ORDER_DELIVERY {
			return errors.ErrPermissionDenied
		}
		if dto.Status != order.ORDER_PENDING && dto.Status != order.ORDER_REFUND {
			return errors.ErrPermissionDenied
		}
	}

	orderDAO.Status = dto.Status

	if err := o.orderRepo.Update(*orderDAO); err != nil {
		return err
	}

	return nil
}

func (o orderService) UpdateOrder(ctx context.Context, dto *orderDTO.UpdateOrderRequest) error {
	//TODO implement me
	panic("implement me")
}
