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
	items, err := o.orderRepo.GetTotalOrderInSystemInDay(ctx, dto.Date)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.AdminTotalOrderInDayResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}

func (o orderService) AdminGetTotalOrderInSystemInMonth(ctx context.Context, dto *statistic.AdminTotalOrderInMonthRequest) (*statistic.AdminTotalOrderInMonthResponse, error) {
	_, m, y := ParseDate(dto.Date)
	items, err := o.orderRepo.GetTotalOrderInSystemInMonth(ctx, m, y)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.AdminTotalOrderInMonthResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}

func (o orderService) AdminGetTotalOrderInSystemInYear(ctx context.Context, dto *statistic.AdminGetTotalOrderInYearRequest) (*statistic.AdminGetTotalOrderInYearResponse, error) {
	items, err := o.orderRepo.GetTotalOrderInSystemInYear(ctx, dto.Year, dto.Count)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.AdminGetTotalOrderInYearResponse{Items: items}
	return &dataResp, nil
}

func (o orderService) AdminGetTotalCommissionOrderInYear(ctx context.Context, dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error) {

	items, err := o.orderRepo.GetTotalCommissionOrderInYear(ctx, dto.Date)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.OrderCommissionDetailResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}

func (o orderService) AdminListOfProductSoldOnMonth(ctx context.Context, dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error) {
	items, err := o.orderRepo.ListOfProductSelledOnMonth(ctx, dto.Date, dto.Count)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.ListOfProductSoldResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}

func (o orderService) GetTotalOrderInMonthOfStore(ctx context.Context, dto *statistic.GetTotalStoreOrderInMonthRequest) (*statistic.GetTotalOrderInMonthResponse, error) {
	items, err := o.orderRepo.GetTotalOrderInSystemInMonthOfStore(ctx, dto.Date, dto.StoreId)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.GetTotalOrderInMonthResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}
func (o orderService) GetTotalOrderInYearOfStore(ctx context.Context, dto *statistic.GetTotalOrderInYearOfStoreRequest) (*statistic.GetTotalOrderInYearOfStoreResponse, error) {
	items, err := o.orderRepo.GetTotalOrderInSystemInYearOfStore(ctx, dto.Year, dto.StoreID, dto.Count)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.GetTotalOrderInYearOfStoreResponse{Items: items}
	return &dataResp, nil
}
func (o orderService) GetTotalStoreCommissionInYear(ctx context.Context, dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error) {
	items, err := o.orderRepo.GetTotalCommissionOrderInYearOfStore(ctx, dto.Date, dto.StoreId)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.OrderCommissionDetailResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}

func (o orderService) ListOfProductSoldOnMonthStore(ctx context.Context, dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error) {
	items, err := o.orderRepo.ListOfProductSelledOnMonthStore(ctx, dto.Date, dto.Count, dto.StoreId)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.ListOfProductSoldResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}
