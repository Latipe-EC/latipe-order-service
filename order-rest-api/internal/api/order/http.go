package order

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"order-rest-api/internal/app/orders"
	"order-rest-api/internal/common/errors"
	dto "order-rest-api/internal/domain/dto/order"
	"order-rest-api/internal/domain/dto/order/delivery"
	"order-rest-api/internal/domain/dto/order/store"
	"order-rest-api/internal/middleware/auth"
	"order-rest-api/internal/responses"
	"order-rest-api/pkg/util/pagable"
	"order-rest-api/pkg/util/valid"
	"strings"
)

type OrderApiHandler interface {
	CreateOrder(ctx *fiber.Ctx) error
	UpdateOrderStatus(ctx *fiber.Ctx) error
	ListOfOrder(ctx *fiber.Ctx) error
	GetOrderByUUID(ctx *fiber.Ctx) error
	GetMyOrder(ctx *fiber.Ctx) error
	CancelOrder(ctx *fiber.Ctx) error
	GetMyStoreOrder(ctx *fiber.Ctx) error
	GetStoreOrderDetail(ctx *fiber.Ctx) error
	UpdateOrderItemStatus(ctx *fiber.Ctx) error
	UpdateStatusByDelivery(ctx *fiber.Ctx) error
	GetOrdersByDelivery(ctx *fiber.Ctx) error
}

type orderApiHandler struct {
	orderUsecase orders.Usecase
}

func NewOrderHandler(orderUsecase orders.Usecase) OrderApiHandler {
	return orderApiHandler{
		orderUsecase: orderUsecase,
	}
}

func (o orderApiHandler) CreateOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	bodyReq := dto.CreateOrderRequest{}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	username := fmt.Sprintf("%v", ctx.Locals(auth.USERNAME))
	if username == "" {
		return errors.ErrUnauthenticated
	}

	bearerToken := fmt.Sprintf("%v", ctx.Locals(auth.BEARER_TOKEN))
	if bearerToken == "" {
		return errors.ErrUnauthenticated
	}

	if err := ctx.BodyParser(&bodyReq); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	if err := valid.GetValidator().Validate(bodyReq); err != nil {
		return errors.ErrBadRequest
	}

	bodyReq.Header.BearerToken = bearerToken
	bodyReq.UserRequest.UserId = userId
	bodyReq.UserRequest.Username = username

	dataResp, err := o.orderUsecase.ProcessCacheOrder(context, &bodyReq)
	if err != nil {
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (o orderApiHandler) GetMyOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	req := new(dto.GetByUserIdRequest)
	req.Query = query

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	req.UserId = userId

	result, err := o.orderUsecase.GetOrderByUserId(context, req)
	if err != nil {
		return errors.ErrInternalServer
	}

	resp := responses.DefaultSuccess
	resp.Data = result

	return resp.JSON(ctx)
}

func (o orderApiHandler) CancelOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := new(dto.CancelOrderRequest)

	if err := ctx.BodyParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	req.UserId = userId

	err := o.orderUsecase.CancelOrder(context, req)
	if err != nil {
		return errors.ErrInternalServer
	}

	resp := responses.DefaultSuccess
	return resp.JSON(ctx)
}

func (o orderApiHandler) UpdateOrderStatus(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := new(dto.UpdateOrderStatusRequest)

	if err := ctx.BodyParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	role := fmt.Sprintf("%v", ctx.Locals(auth.ROLE))
	if role == "" {
		return errors.ErrPermissionDenied
	}

	req.UserId = userId
	req.Role = role

	err := o.orderUsecase.UpdateStatusOrder(context, req)
	if err != nil {
		return errors.ErrInternalServer
	}

	return responses.DefaultSuccess.JSON(ctx)
}

func (o orderApiHandler) ListOfOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	req := new(dto.GetOrderListRequest)
	req.Query = query

	result, err := o.orderUsecase.GetOrderList(context, req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) GetOrderByUUID(ctx *fiber.Ctx) error {
	context := ctx.Context()
	req := new(dto.GetOrderByUUIDRequest)

	if err := ctx.ParamsParser(req); err != nil {
		return errors.ErrInternalServer
	}

	result, err := o.orderUsecase.GetOrderByUUID(context, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return errors.ErrNotFound
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) GetMyStoreOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	req := new(store.GetStoreOrderRequest)
	req.Query = query

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	result, err := o.orderUsecase.GetOrdersOfStore(context, req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) GetStoreOrderDetail(ctx *fiber.Ctx) error {
	context := ctx.Context()

	var req store.GetOrderOfStoreByIDRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	result, err := o.orderUsecase.ViewDetailStoreOrder(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) UpdateOrderItemStatus(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := store.UpdateOrderItemRequest{}

	if err := ctx.BodyParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	if err := ctx.ParamsParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	storeId := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	if storeId == "" {
		return errors.ErrUnauthenticated
	}

	req.StoreId = storeId

	resp, err := o.orderUsecase.UpdateOrderItem(context, &req)
	if err != nil {
		return errors.ErrInternalServer
	}

	data := responses.DefaultSuccess
	data.Data = resp

	return data.JSON(ctx)
}

func (o orderApiHandler) UpdateStatusByDelivery(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := delivery.UpdateOrderStatusRequest{}

	if err := ctx.BodyParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	if err := ctx.ParamsParser(&req); err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	deli := fmt.Sprintf("%v", ctx.Locals(auth.DELIVERY_ID))
	if deli == "" {
		return errors.ErrUnauthenticated
	}

	req.DeliveryID = deli

	if err := valid.GetValidator().Validate(req); err != nil {
		return errors.ErrBadRequest
	}

	resp, err := o.orderUsecase.DeliveryUpdateStatusOrder(context, req)
	if err != nil {
		return errors.ErrInternalServer
	}

	data := responses.DefaultSuccess
	data.Data = resp

	return data.JSON(ctx)
}

func (o orderApiHandler) GetOrdersByDelivery(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	req := delivery.GetOrderListRequest{}
	req.Query = query

	deliId := fmt.Sprintf("%v", ctx.Locals(auth.DELIVERY_ID))
	req.DeliveryID = deliId

	result, err := o.orderUsecase.GetOrdersOfDelivery(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}
