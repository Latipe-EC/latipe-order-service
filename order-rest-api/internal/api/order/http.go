package order

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"order-rest-api/internal/app/orders"
	"order-rest-api/internal/common/errors"
	dto "order-rest-api/internal/domain/dto/order"
	"order-rest-api/internal/domain/dto/order/delivery"
	internalDTO "order-rest-api/internal/domain/dto/order/internal-service"
	"order-rest-api/internal/domain/dto/order/store"
	"order-rest-api/internal/middleware/auth"
	"order-rest-api/internal/responses"
	"order-rest-api/pkg/util/pagable"
	"order-rest-api/pkg/util/valid"
	"strings"
)

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

func (o orderApiHandler) UserCancelOrder(ctx *fiber.Ctx) error {
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
		return err
	}

	resp := responses.DefaultSuccess
	return resp.JSON(ctx)
}

func (o orderApiHandler) UserRefundOrder(ctx *fiber.Ctx) error {
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
		return err
	}

	resp := responses.DefaultSuccess
	return resp.JSON(ctx)
}

func (o orderApiHandler) AdminCancelOrder(ctx *fiber.Ctx) error {
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

	err := o.orderUsecase.AdminCancelOrder(context, req)
	if err != nil {
		return err
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

func (o orderApiHandler) UserGetOrderByUUID(ctx *fiber.Ctx) error {
	context := ctx.Context()
	req := new(dto.GetOrderByUUIDRequest)

	if err := ctx.ParamsParser(req); err != nil {
		return errors.ErrInternalServer
	}

	role := fmt.Sprintf("%v", ctx.Locals(auth.ROLE))
	if role == "" {
		return errors.ErrPermissionDenied
	}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	req.OwnerId = userId
	req.Role = role

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

func (o orderApiHandler) DeliveryGetOrderByUUID(ctx *fiber.Ctx) error {
	context := ctx.Context()
	req := new(dto.GetOrderByUUIDRequest)

	if err := ctx.ParamsParser(req); err != nil {
		return errors.ErrInternalServer
	}

	role := fmt.Sprintf("%v", ctx.Locals(auth.ROLE))
	if role == "" {
		return errors.ErrPermissionDenied
	}

	deliID := fmt.Sprintf("%v", ctx.Locals(auth.DELIVERY_ID))
	if deliID == "" {
		return errors.ErrUnauthenticated
	}

	req.OwnerId = deliID
	req.Role = role

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

func (o orderApiHandler) InternalGetOrderByUUID(ctx *fiber.Ctx) error {
	context := ctx.Context()
	req := internalDTO.GetOrderRatingItemRequest{}

	if err := ctx.ParamsParser(&req); err != nil {
		return errors.ErrInternalServer
	}

	result, err := o.orderUsecase.InternalGetRatingID(context, &req)
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

	req := store.GetStoreOrderRequest{}
	req.Query = query

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	result, err := o.orderUsecase.GetOrdersOfStore(context, &req)
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
		case errors.Is(err, gorm.ErrRecordNotFound):
			return errors.ErrNotFoundRecord
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
		return err
	}

	data := responses.DefaultSuccess
	data.Data = resp

	return data.JSON(ctx)
}

func (o orderApiHandler) CancelOrderItemStatus(ctx *fiber.Ctx) error {
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

	resp, err := o.orderUsecase.CancelOrderItem(context, &req)
	if err != nil {
		return err
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
		return err
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

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

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

func (o orderApiHandler) SearchOrderIdByKeyword(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	req := store.FindStoreOrderRequest{}
	req.Query = query

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	result, err := o.orderUsecase.SearchStoreOrderId(context, &req)
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

func (o orderApiHandler) AdminCountingOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	result, err := o.orderUsecase.AdminCountingOrderAmount(context, &req)
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

func (o orderApiHandler) UserCountingOrder(ctx *fiber.Ctx) error {

	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	req.OwnerID = userId

	result, err := o.orderUsecase.UserCountingOrder(context, &req)
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

func (o orderApiHandler) StoreCountingOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	storeId := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.OwnerID = storeId

	result, err := o.orderUsecase.StoreCountingOrder(context, &req)
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

func (o orderApiHandler) DeliveryCountingOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	deli := fmt.Sprintf("%v", ctx.Locals(auth.DELIVERY_ID))
	req.OwnerID = deli

	result, err := o.orderUsecase.DeliveryCountingOrder(context, &req)
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
