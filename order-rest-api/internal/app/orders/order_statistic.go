package orders

import (
	"context"
	orderDTO "order-rest-api/internal/domain/dto/order"
	"order-rest-api/internal/domain/dto/order/statistic"
)

func (o orderService) AdminCountingOrderAmount(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.AdminCountingOrder(ctx)
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderService) UserCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.UserCountingOrder(ctx, dto.OwnerID)
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderService) StoreCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.StoreCountingOrder(ctx, dto.OwnerID)
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderService) DeliveryCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.StoreCountingOrder(ctx, dto.OwnerID)
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderService) AdminGetTotalOrderInSystemInDay(ctx context.Context, dto *statistic.AdminTotalOrderInDayRequest) (*statistic.AdminTotalOrderInDayResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) AdminGetTotalOrderInSystemInMonth(ctx context.Context, dto *statistic.AdminTotalOrderInMonthRequest) (*statistic.AdminTotalOrderInMonthResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) AdminGetTotalOrderInSystemInYear(ctx context.Context, dto *statistic.AdminGetTotalOrderInYearRequest) (*statistic.AdminGetTotalOrderInYearResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) AdminGetTotalCommissionOrderInYear(ctx context.Context, dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) AdminListOfProductSoldOnMonth(ctx context.Context, dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) GetTotalOrderInMonthOfStore(ctx context.Context, dto *statistic.GetTotalStoreOrderInMonthRequest) (statistic.GetTotalOrderInMonthResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) GetTotalOrderInYearOfStore(ctx context.Context, dto *statistic.GetTotalOrderInYearOfStoreRequest) (*statistic.GetTotalOrderInYearOfStoreResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) GetTotalStoreCommissionInYear(ctx context.Context, dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) ListOfProductSoldOnMonthStore(ctx context.Context, dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error) {
	//TODO implement me
	panic("implement me")
}
