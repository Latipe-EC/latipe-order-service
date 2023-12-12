package orders

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"order-rest-api/config"
	"order-rest-api/internal/common/errors"
	orderDTO "order-rest-api/internal/domain/dto/order"
	"order-rest-api/internal/domain/dto/order/delivery"
	internalDTO "order-rest-api/internal/domain/dto/order/internal-service"
	"order-rest-api/internal/domain/dto/order/store"
	"order-rest-api/internal/domain/entities/order"
	"order-rest-api/internal/domain/msg"
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
	"strings"
)

type orderService struct {
	orderRepo   order.Repository
	cacheEngine *redis.CacheEngine
	productServ productserv.Service
	userServ    userserv.Service
	deliServ    deliveryserv.Service
	voucherSer  voucherserv.Service
	cfg         *config.Config
}

func NewOrderService(cfg *config.Config, orderRepo order.Repository, productServ productserv.Service,
	cacheEngine *redis.CacheEngine, userServ userserv.Service, deliServ deliveryserv.Service,
	voucherServ voucherserv.Service) Usecase {
	return orderService{
		orderRepo:   orderRepo,
		cacheEngine: cacheEngine,
		productServ: productServ,
		userServ:    userServ,
		deliServ:    deliServ,
		voucherSer:  voucherServ,
		cfg:         cfg,
	}
}

func (o orderService) InternalGetRatingID(ctx context.Context, dto *internalDTO.GetOrderRatingItemRequest) (*internalDTO.GetOrderRatingItemResponse, error) {
	resp := internalDTO.GetOrderRatingItemResponse{}
	orderDAO, err := o.orderRepo.FindByItemId(ctx, dto.ItemID)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orderDAO, &resp); err != nil {
		return nil, err
	}

	return &resp, err
}

func (o orderService) CancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
	dao, err := o.orderRepo.FindByUUID(ctx, dto.OrderUUID)
	if err != nil {
		return err
	}

	if dao.UserId != dto.UserId {
		return errors.ErrNotFoundRecord
	}

	if dao.Status == order.ORDER_CANCEL {
		return errors.ErrNotChange
	}

	if dao.Status != order.ORDER_CREATED {
		return errors.OrderCannotCancel
	}

	if err := o.orderRepo.UpdateStatus(ctx, dao.Id, order.ORDER_CANCEL); err != nil {
		return err
	}

	mess := msg.OrderMessage{
		OrderUUID:     dao.OrderUUID,
		Status:        order.ORDER_CANCEL,
		PaymentMethod: dao.PaymentMethod,
	}

	if err := message.SendOrderMessage(&mess); err != nil {
		return err
	}

	return nil
}

func (o orderService) UserRefundOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
	dao, err := o.orderRepo.FindByUUID(ctx, dto.OrderUUID)
	if err != nil {
		return err
	}

	if dao.UserId != dto.UserId {
		return errors.ErrNotFoundRecord
	}

	if dao.Status == order.ORDER_REFUND {
		return errors.ErrNotChange
	}

	if dao.Status != order.ORDER_SHIPPING_FINISH {
		return errors.OrderCannotCancel
	}

	if err := o.orderRepo.UpdateStatus(ctx, dao.Id, order.ORDER_REFUND); err != nil {
		return err
	}

	mess := msg.OrderMessage{
		OrderUUID:     dao.OrderUUID,
		Status:        order.ORDER_REFUND,
		PaymentMethod: dao.PaymentMethod,
	}

	if err := message.SendOrderMessage(&mess); err != nil {
		return err
	}

	return nil
}

func (o orderService) AdminCancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error {
	dao, err := o.orderRepo.FindByUUID(ctx, dto.OrderUUID)
	if err != nil {
		return err
	}

	if dao.Status == order.ORDER_CANCEL {
		return errors.ErrNotChange
	}

	if err := o.orderRepo.UpdateStatus(ctx, dao.Id, order.ORDER_CANCEL); err != nil {
		return err
	}

	mess := msg.OrderMessage{
		OrderUUID:     dao.OrderUUID,
		Status:        order.ORDER_CANCEL,
		PaymentMethod: dao.PaymentMethod,
	}

	if err := message.SendOrderMessage(&mess); err != nil {
		return err
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
		voucherReq := voucherDTO.CheckingVoucherRequest{}
		voucherReq.Vouchers = dto.VoucherCode
		voucherReq.AuthorizationHeader.BearerToken = dto.Header.BearerToken
		voucherReq.OrderTotalAmount = orderData.SubTotal - orderData.ShippingCost
		voucherReq.PaymentMethod = orderData.PaymentMethod
		voucherReq.UserId = orderData.UserRequest.UserId

		voucherDetail, err := o.voucherSer.CheckingVoucher(ctx, &voucherReq)
		if err != nil {
			return nil, err
		}

		if voucherDetail.IsSuccess == true {

			for _, v := range voucherDetail.Items {
				switch v.VoucherType {
				case voucherDTO.FREE_SHIP:
					if shippingDetail.Cost < v.DiscountValue {
						orderData.ShippingDiscount = shippingDetail.Cost
					} else {
						orderData.ShippingDiscount = v.DiscountValue
					}

				case voucherDTO.DISCOUNT_ORDER:
					orderData.ItemDiscount = v.DiscountValue
				}

			}
			orderData.Vouchers = dto.VoucherCode
		}

	}
	//calculate amount order

	orderData.SubTotal += orderData.ShippingCost
	orderData.Amount = orderData.SubTotal - (orderData.ItemDiscount + orderData.ShippingDiscount)
	orderData.Status = order.ORDER_SYSTEM_PROCESS
	//gen key order

	orderKey := o.genOrderKey(orderData.UserRequest.UserId)
	orderData.OrderUUID = orderKey

	if err := message.SendOrderMessage(orderData); err != nil {
		return nil, err
	}

	data := orderDTO.CreateOrderResponse{
		UserOrder: orderDTO.UserRequest{
			UserId:   dto.UserRequest.UserId,
			Username: dto.UserRequest.Username,
		},
		OrderKey:         orderKey,
		ShippingCost:     orderData.ShippingCost,
		Amount:           orderData.Amount,
		ShippingDiscount: orderData.ShippingDiscount,
		ItemDiscount:     orderData.ItemDiscount,
		SubTotal:         orderData.SubTotal,
		PaymentMethod:    orderData.PaymentMethod,
	}

	return &data, nil
}

func (o orderService) GetOrderById(ctx context.Context, dto *orderDTO.GetOrderByIDRequest) (*orderDTO.GetOrderResponse, error) {
	orderResp := orderDTO.OrderResponse{}

	orderDAO, err := o.orderRepo.FindById(ctx, dto.OrderId)
	if err != nil {
		return nil, err
	}

	switch dto.Role {
	case auth.ROLE_USER:
		if orderDAO.UserId != dto.OwnerId {
			return nil, errors.ErrNotFound
		}
	case auth.ROLE_STORE:
		if !CheckStoreHaveOrder(*orderDAO, dto.OwnerId) {
			return nil, errors.ErrNotFound
		}
	}

	if err = mapper.BindingStruct(orderDAO, &orderResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderResponse{Order: orderResp}

	return &resp, err
}

func (o orderService) GetOrderByUUID(ctx context.Context, dto *orderDTO.GetOrderByUUIDRequest) (*orderDTO.GetOrderResponse, error) {
	orderResp := orderDTO.OrderResponse{}

	orderDAO, err := o.orderRepo.FindByUUID(ctx, dto.OrderId)
	if err != nil {
		return nil, err
	}

	switch dto.Role {
	case auth.ROLE_USER:
		if orderDAO.UserId != dto.OwnerId {
			return nil, errors.ErrNotFound
		}
	case auth.ROLE_STORE:
		if !CheckStoreHaveOrder(*orderDAO, dto.OwnerId) {
			return nil, errors.ErrNotFound
		}
	case auth.ROLE_DELIVERY:
		if orderDAO.Delivery.DeliveryId != dto.OwnerId {
			return nil, errors.ErrNotFound
		}
	}

	if err = mapper.BindingStruct(orderDAO, &orderResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderResponse{Order: orderResp}

	return &resp, err
}

func (o orderService) GetOrderList(ctx context.Context, dto *orderDTO.GetOrderListRequest) (*orderDTO.GetOrderListResponse, error) {
	var dataResp []orderDTO.OrderResponse

	orders, err := o.orderRepo.FindAll(ctx, dto.Query)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.Total(ctx, dto.Query)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderListResponse{}
	resp.Items = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderService) GetOrderByUserId(ctx context.Context, dto *orderDTO.GetByUserIdRequest) (*orderDTO.GetByUserIdResponse, error) {
	var dataResp []orderDTO.OrderResponse

	orders, err := o.orderRepo.FindByUserId(ctx, dto.UserId, dto.Query)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.Total(ctx, dto.Query)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetByUserIdResponse{}
	resp.Items = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderService) SearchStoreOrderId(ctx context.Context, dto *store.FindStoreOrderRequest) (*orderDTO.GetOrderListResponse, error) {
	var dataResp []store.StoreOrderResponse

	orders, err := o.orderRepo.SearchOrderByStoreID(ctx, dto.StoreID, dto.Keyword, dto.Query)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.TotalSearchOrderByStoreID(ctx, dto.StoreID, dto.Keyword)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderListResponse{}
	resp.Items = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderService) GetOrdersOfStore(ctx context.Context, dto *store.GetStoreOrderRequest) (*orderDTO.GetOrderListResponse, error) {
	var dataResp []store.StoreOrderResponse

	orders, err := o.orderRepo.FindOrderByStoreID(ctx, dto.StoreID, dto.Query, dto.Keyword)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.TotalStoreOrder(ctx, dto.StoreID, dto.Query, dto.Keyword)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderListResponse{}
	resp.Items = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderService) GetOrdersOfDelivery(ctx context.Context, dto *delivery.GetOrderListRequest) (*delivery.GetOrderListResponse, error) {
	var dataResp []store.StoreOrderResponse

	orders, err := o.orderRepo.FindOrderByDelivery(ctx, dto.DeliveryID, dto.Query)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.TotalOrdersOfDelivery(ctx, dto.DeliveryID, dto.Keyword, dto.Query)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := delivery.GetOrderListResponse{}
	resp.Items = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderService) ViewDetailStoreOrder(ctx context.Context, dto *store.GetOrderOfStoreByIDRequest) (*store.GetOrderOfStoreByIDResponse, error) {
	orderResp := store.GetOrderOfStoreByIDResponse{}

	orderDAO, err := o.orderRepo.FindByUUID(ctx, dto.OrderUUID)
	if err != nil {
		return nil, err
	}

	if !CheckStoreHaveOrder(*orderDAO, dto.StoreID) {
		return nil, errors.ErrNotFound
	}

	if err = mapper.BindingStruct(orderDAO, &orderResp); err != nil {
		return nil, err
	}

	storeAmount := 0
	var items []store.OrderStoreItem
	for _, o := range orderDAO.OrderItem {
		if o.StoreID == dto.StoreID {
			i := store.OrderStoreItem{
				ProductId:   o.ProductID,
				OptionId:    o.OptionID,
				Quantity:    o.Quantity,
				Price:       o.Price,
				Status:      o.Status,
				Id:          o.Id,
				SubTotal:    o.SubTotal,
				ProdImg:     o.ProdImg,
				ProductName: o.ProductName,
				NetPrice:    o.NetPrice,
			}
			items = append(items, i)
			storeAmount += o.SubTotal
		}
	}

	if len(items) < 1 {
		return nil, errors.ErrNotFoundRecord
	}

	orderResp.CommissionDetail.SystemFee = orderDAO.OrderCommission.SystemFee
	orderResp.CommissionDetail.AmountReceived = orderDAO.OrderCommission.AmountReceived
	orderResp.StoreOrderAmount = storeAmount
	orderResp.OrderItems = items

	return &orderResp, err
}

func (o orderService) UpdateStatusOrder(ctx context.Context, dto *orderDTO.UpdateOrderStatusRequest) error {

	orderDAO, err := o.orderRepo.FindByUUID(ctx, dto.OrderUUID)
	if err != nil {
		return err
	}

	if orderDAO.Status == order.ORDER_CANCEL {
		return errors.ErrBadRequest
	}

	orderDAO.Status = dto.Status

	if err := o.orderRepo.Update(ctx, *orderDAO); err != nil {
		return err
	}

	return nil
}

func (o orderService) DeliveryUpdateStatusOrder(ctx context.Context, dto delivery.UpdateOrderStatusRequest) (*delivery.UpdateOrderStatusResponse, error) {
	orderDAO, err := o.orderRepo.FindByUUID(ctx, dto.OrderUUID)
	if err != nil {
		return nil, err
	}

	if orderDAO.Status != order.ORDER_DELIVERY {
		return nil, errors.OrderStatusNotValid
	}

	if orderDAO.Delivery.DeliveryId != dto.DeliveryID {
		return nil, errors.ErrNotFoundRecord
	}

	orderDAO.Status = dto.Status

	if dto.Status == order.ORDER_CANCEL || dto.Status == order.ORDER_SHIPPING_FINISH {
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.Id, dto.Status); err != nil {
			return nil, err
		}
	}

	msg := msg.OrderMessage{
		Status:    dto.Status,
		OrderUUID: orderDAO.OrderUUID,
	}

	if dto.Status == order.ORDER_SHIPPING_FINISH {
		if err := message.SendOrderMessage(&msg); err != nil {
			return nil, err
		}
	}

	if dto.Status == order.ORDER_CANCEL {
		if err := message.SendOrderMessage(&msg); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (o orderService) UpdateOrderItem(ctx context.Context, dto *store.UpdateOrderItemRequest) (*store.UpdateOrderItemResponse, error) {
	orderDAO, err := o.orderRepo.FindByUUID(ctx, dto.OrderUUID)
	if err != nil {
		return nil, err
	}

	notFound := true
	itemPreparedCount := 0

	for _, i := range orderDAO.OrderItem {

		if i.StoreID != dto.StoreId {
			continue
		}

		if i.Id == dto.ItemID {
			notFound = false
			if i.Status != order.OI_PREPARED && i.Status != order.OI_CANCEL {
				if err := o.orderRepo.UpdateOrderItem(ctx, i.Id, order.OI_PREPARED); err != nil {
					return nil, err
				}
				i.Status = order.OI_PREPARED
			} else {
				return nil, errors.ErrNotChange
			}
		}

		if i.Status == order.OI_PREPARED {
			itemPreparedCount++
		}
	}

	if notFound {
		return nil, errors.ErrNotFoundRecord
	}

	if orderDAO.Status == order.ORDER_CREATED {
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.Id, order.ORDER_PENDING); err != nil {
			return nil, err
		}
	}

	if len(orderDAO.OrderItem) == itemPreparedCount {
		if err := o.orderRepo.UpdateStatus(ctx, orderDAO.Id, order.ORDER_DELIVERY); err != nil {
			return nil, err
		}
	}

	resp := store.UpdateOrderItemResponse{
		OrderUUID: dto.OrderUUID,
		ItemID:    dto.ItemID,
		Status:    order.OI_PREPARED,
	}

	return &resp, nil
}

func (o orderService) CancelOrderItem(ctx context.Context, dto *store.UpdateOrderItemRequest) (*store.UpdateOrderItemResponse, error) {
	orderDAO, err := o.orderRepo.FindByUUID(ctx, dto.OrderUUID)
	if err != nil {
		return nil, err
	}

	notFound := true

	for _, i := range orderDAO.OrderItem {

		if i.StoreID != dto.StoreId {
			continue
		}

		if i.Status != order.OI_PREPARED && i.Id == dto.ItemID {
			notFound = false
			if err := o.orderRepo.UpdateOrderItem(ctx, i.Id, order.ORDER_CANCEL); err != nil {
				return nil, err
			}
			i.Status = order.OI_PREPARED
		}

	}

	if notFound {
		return nil, errors.ErrNotFoundRecord
	}

	if err := o.orderRepo.UpdateStatus(ctx, orderDAO.Id, order.ORDER_CANCEL,
		"Đơn hàng bị hủy do nhà cung cấp không thể chuẩn bị sản phẩm"); err != nil {
		return nil, err
	}

	resp := store.UpdateOrderItemResponse{
		OrderUUID: dto.OrderUUID,
		ItemID:    dto.ItemID,
		Status:    order.ORDER_CANCEL,
	}

	return &resp, nil
}

func (o orderService) UpdateOrder(ctx context.Context, dto *orderDTO.UpdateOrderRequest) error {
	//TODO implement me
	panic("implement me")
}

func (o orderService) initOrderCacheData(products *prodServDTO.OrderProductResponse,
	address *userDTO.GetDetailAddressResponse, deli *deliDto.GetShippingCostResponse, dto *orderDTO.CreateOrderRequest) *msg.OrderMessage {

	orderCache := msg.OrderMessage{
		Status: order.ORDER_SYSTEM_PROCESS,
		Header: msg.BaseHeader{dto.Header.BearerToken},
		UserRequest: msg.UserRequest{
			UserId:   dto.UserRequest.UserId,
			Username: dto.UserRequest.Username,
		},
		SubTotal:      products.TotalPrice,
		PaymentMethod: dto.PaymentMethod,
		Vouchers:      nil,
		Address: msg.OrderAddress{
			AddressId:       address.Id,
			ShippingName:    address.ContactName,
			ShippingPhone:   address.Phone,
			ShippingAddress: address.DetailAddress,
		},
		Delivery: msg.Delivery{
			DeliveryId:    deli.DeliveryId,
			Name:          deli.DeliveryName,
			Cost:          deli.Cost,
			ReceivingDate: deli.ReceiveDate,
		},
	}

	discount := 0
	//order detail
	var orderItems []msg.OrderItemsCache
	for index, i := range products.Products {
		item := msg.OrderItemsCache{
			CartId: dto.OrderItems[index].CartId,
			ProductItem: msg.ProductItem{
				ProductID:   i.ProductId,
				ProductName: i.Name,
				StoreID:     i.StoreId,
				NameOption:  i.NameOption,
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
	orders, err := o.orderRepo.FindOrderByUserAndProduct(ctx, dto.UserId, dto.ProductId)
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

func (o orderService) genOrderKey(userId string) string {
	keyGen := strings.ReplaceAll(uuid.NewString(), "-", "")[:10]
	key := fmt.Sprintf("%v%v%v", o.cfg.Server.KeyID, userId[:4], keyGen)

	return key
}
