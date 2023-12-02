package orders

import (
	"context"
	orderDTO "order-rest-api/internal/domain/dto/order"
	"order-rest-api/internal/domain/dto/order/statistic"
)

func (o orderService) AdminCountingOrderAmount(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.AdminCountingOrder()
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderService) UserCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.UserCountingOrder(dto.OwnerID)
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderService) StoreCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.StoreCountingOrder(dto.OwnerID)
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderService) DeliveryCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.StoreCountingOrder(dto.OwnerID)
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderService) AdminGetTotalOrderInSystemInDay(dto *statistic.AdminTotalOrderInDayRequest) (*statistic.AdminTotalOrderInDayResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) AdminGetTotalOrderInSystemInMonth(dto *statistic.AdminTotalOrderInMonthRequest) (*statistic.AdminTotalOrderInMonthResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) AdminGetTotalOrderInSystemInYear(dto *statistic.AdminGetTotalOrderInYearRequest) (*statistic.AdminGetTotalOrderInYearResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) AdminGetTotalCommissionOrderInYear(dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) AdminListOfProductSoldOnMonth(dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) GetTotalOrderInMonthOfStore(dto *statistic.GetTotalStoreOrderInMonthRequest) (statistic.GetTotalOrderInMonthResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) GetTotalOrderInYearOfStore(dto *statistic.GetTotalOrderInYearOfStoreRequest) (*statistic.GetTotalOrderInYearOfStoreResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) GetTotalStoreCommissionInYear(dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o orderService) ListOfProductSoldOnMonthStore(dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error) {
	//TODO implement me
	panic("implement me")
}
