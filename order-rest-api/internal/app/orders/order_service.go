package orders

import (
	"context"
	"github.com/google/uuid"
	"order-rest-api/internal/common/errors"
	orderDTO "order-rest-api/internal/domain/dto/order"
	"order-rest-api/internal/domain/dto/order/store"
	"order-rest-api/internal/domain/entities/order"
	"order-rest-api/internal/infrastructure/adapter/deliveryserv"
	deliDto "order-rest-api/internal/infrastructure/adapter/deliveryserv/dto"
	"order-rest-api/internal/infrastructure/adapter/productserv"
	prodServDTO "order-rest-api/internal/infrastructure/adapter/productserv/dto"
	"order-rest-api/internal/infrastructure/adapter/userserv"
	userDTO "order-rest-api/internal/infrastructure/adapter/userserv/dto"
	voucherserv "order-rest-api/internal/infrastructure/adapter/vouchersev"
	voucherDTO "order-rest-api/internal/infrastructure/adapter/vouchersev/dto"
	"order-rest-api/internal/message"
	"order-rest-api/internal/middleware/auth"
	"order-rest-api/pkg/cache/redis"
	"order-rest-api/pkg/util/mapper"
)

type orderService struct {
	orderRepo   order.Repository
	cacheEngine *redis.CacheEngine
	productServ productserv.Service
	userServ    userserv.Service
	deliServ    deliveryserv.Service
	voucherSer  voucherserv.Service
}

func NewOrderService(orderRepo order.Repository, productServ productserv.Service,
	cacheEngine *redis.CacheEngine, userServ userserv.Service, deliServ deliveryserv.Service,
	voucherServ voucherserv.Service) Usecase {
	return orderService{
		orderRepo:   orderRepo,
		cacheEngine: cacheEngine,
		productServ: productServ,
		userServ:    userServ,
		deliServ:    deliServ,
		voucherSer:  voucherServ,
	}
}

func (o orderService) CancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
	dao, err := o.orderRepo.FindByUUID(dto.OrderUUID)
	if err != nil {
		return err
	}

	if dao.Status == order.ORDER_CREATED && dao.UserId == dto.UserId {
		if err := o.orderRepo.UpdateStatus(dao.Id, order.ORDER_CANCEL); err != nil {
			return err
		}
	}

	return nil
}

func (o orderService) ProcessCacheOrder(ctx context.Context, dto *orderDTO.CreateOrderRequest) (*orderDTO.CreateOrderResponse, error) {

	productReq := prodServDTO.OrderProductRequest{
		Items: MappingOrderItemToGetInfo(dto),
	}
	products, err := o.productServ.GetProductOrderInfo(ctx, &productReq)
	if err != nil {
		return nil, err
	}

	addressRequest := userDTO.GetDetailAddressRequest{
		AddressId:           dto.Address.AddressId,
		AuthorizationHeader: userDTO.AuthorizationHeader{BearerToken: dto.Header.BearerToken},
	}
	userAddress, err := o.userServ.GetAddressDetails(ctx, &addressRequest)
	if err != nil {
		return nil, err
	}

	shippingReq := deliDto.GetShippingCostRequest{
		SrcCode:    products.StoreProvinceCodes,
		DestCode:   userAddress.CityOrProvinceId,
		DeliveryId: dto.Delivery.DeliveryId,
	}
	shippingDetail, err := o.deliServ.CalculateShippingCost(ctx, &shippingReq)
	if err != nil {
		return nil, err
	}

	orderData := o.initOrderCacheData(products, userAddress, shippingDetail, dto)

	if len(dto.VoucherCode) > 0 {
		voucherReq := voucherDTO.ApplyVoucherRequest{}
		voucherReq.Vouchers = dto.VoucherCode
		voucherReq.AuthorizationHeader.BearerToken = dto.Header.BearerToken

		voucherDetail, err := o.voucherSer.ApplyVoucher(ctx, &voucherReq)
		if err != nil {
			return nil, err
		}

		if voucherDetail.IsSuccess == true {

			for _, v := range voucherDetail.Items {
				switch v.VoucherType {
				case voucherDTO.FREE_SHIP:
					if shippingDetail.Cost < v.DiscountValue {
						orderData.ShippingCost = 0
					} else {
						orderData.ShippingCost -= v.DiscountValue
					}

				case voucherDTO.DISCOUNT_ORDER:
					orderData.Discount += v.DiscountValue
				}

			}
			orderData.Vouchers = dto.VoucherCode
		}

	}
	//calculate amount order
	orderData.Amount = orderData.SubTotal + orderData.ShippingCost - orderData.Discount
	//gen key order
	keyGen := uuid.NewString()
	orderData.OrderUUID = keyGen

	if err := message.SendMessage(orderData, orderData.OrderUUID); err != nil {
		return nil, err
	}

	data := orderDTO.CreateOrderResponse{
		UserOrder: orderDTO.UserRequest{
			UserId:   dto.UserRequest.UserId,
			Username: dto.UserRequest.Username,
		},
		OrderKey:      keyGen,
		Amount:        orderData.Amount,
		Discount:      orderData.Discount,
		SubTotal:      orderData.SubTotal,
		PaymentMethod: orderData.PaymentMethod,
	}

	return &data, nil
}

func (o orderService) GetOrderById(ctx context.Context, dto *orderDTO.GetOrderByIDRequest) (*orderDTO.GetOrderResponse, error) {
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

func (o orderService) GetOrderByUUID(ctx context.Context, dto *orderDTO.GetOrderByUUIDRequest) (*orderDTO.GetOrderResponse, error) {
	orderResp := orderDTO.OrderResponse{}

	orderDAO, err := o.orderRepo.FindByUUID(dto.OrderId)
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
	var dataResp []orderDTO.OrderResponse

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

func (o orderService) GetOrdersOfStore(ctx context.Context, dto *store.GetStoreOrderRequest) (*orderDTO.GetOrderListResponse, error) {
	var dataResp []store.StoreOrderResponse

	orders, err := o.orderRepo.FindOrderByStoreID(dto.StoreID, dto.Query)
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

func (o orderService) ViewDetailStoreOrder(ctx context.Context, dto *store.GetOrderOfStoreByIDRequest) (*store.GetOrderOfStoreByIDResponse, error) {
	orderResp := store.GetOrderOfStoreByIDResponse{}

	orderDAO, err := o.orderRepo.FindByUUID(dto.OrderUUID)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orderDAO, &orderResp); err != nil {
		return nil, err
	}

	storeAmount := 0
	var items []store.OrderStoreItem
	for _, o := range orderDAO.OrderItem {
		if o.StoreID == dto.StoreID {
			i := store.OrderStoreItem{
				ProductId: o.ProductID,
				OptionId:  o.OptionID,
				Quantity:  o.Quantity,
				Price:     o.Price,
				Status:    o.Status,
			}
			items = append(items, i)
			storeAmount += o.Price
		}
	}

	if len(items) < 1 {
		return nil, errors.ErrNotFoundRecord
	}

	orderResp.StoreOrderAmount = storeAmount
	orderResp.OrderItems = items

	return &orderResp, err
}

func (o orderService) UpdateStatusOrder(ctx context.Context, dto *orderDTO.UpdateOrderStatusRequest) error {

	orderDAO, err := o.orderRepo.FindByUUID(dto.OrderUUID)
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

func (o orderService) initOrderCacheData(products *prodServDTO.OrderProductResponse,
	address *userDTO.GetDetailAddressResponse, deli *deliDto.GetShippingCostResponse, dto *orderDTO.CreateOrderRequest) *order.OrderMessage {

	orderCache := order.OrderMessage{
		Header: order.BaseHeader{dto.Header.BearerToken},
		UserRequest: order.UserRequest{
			UserId:   dto.UserRequest.UserId,
			Username: dto.UserRequest.Username,
		},
		SubTotal:      products.TotalPrice,
		PaymentMethod: dto.PaymentMethod,
		Vouchers:      nil,
		Address: order.OrderAddress{
			AddressId:       address.Id,
			ShippingName:    address.ContactName,
			ShippingPhone:   address.Phone,
			ShippingAddress: address.DetailAddress,
		},
		Delivery: order.Delivery{
			DeliveryId:    deli.DeliveryId,
			Name:          deli.DeliveryName,
			Cost:          deli.Cost,
			ReceivingDate: deli.ReceiveDate,
		},
	}

	discount := 0
	//order detail
	var orderItems []order.OrderItemsCache
	for index, i := range products.Products {
		item := order.OrderItemsCache{
			CartItemId: dto.OrderItems[index].CartItemId,
			ProductItem: order.ProductItem{
				ProductID:   i.ProductId,
				ProductName: i.Name,
				StoreID:     i.StoreId,
				OptionID:    i.OptionId,
				Quantity:    i.Quantity,
				Price:       int(i.Price),
				NetPrice:    int(i.PromotionalPrice),
				Image:       i.Image,
			},
		}
		orderItems = append(orderItems, item)
		discount += int(i.PromotionalPrice)
	}
	orderCache.OrderItems = orderItems
	orderCache.ShippingCost = deli.Cost

	return &orderCache
}

func (o orderService) CheckProductPurchased(ctx context.Context, dto *orderDTO.CheckUserOrderRequest) (*orderDTO.CheckUserOrderResponse, error) {
	orders, err := o.orderRepo.FindOrderByUserAndProduct(dto.UserId, dto.ProductId)
	if err != nil {
		return nil, err
	}
	data := orderDTO.CheckUserOrderResponse{}

	if len(orders) > 0 {
		var ordersKeys []string
		for _, i := range orders {
			ordersKeys = append(ordersKeys, i.OrderUUID)
		}
		data.IsPurchased = true
		data.Orders = ordersKeys
	}

	return &data, err
}
