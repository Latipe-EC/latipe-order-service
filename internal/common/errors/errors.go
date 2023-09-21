package errors

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"order-service-rest-api/internal/domain/enum"
	"order-service-rest-api/internal/responses"
)

type Error struct {
	Status               int    `json:"-"`
	InternalErrorMessage string `json:"-"`
	Code                 int    `json:"code"`
	ErrorCode            string `json:"error_code"`
	Message              string `json:"message"`
}

func (e *Error) Error() string {
	if e.InternalErrorMessage != "" {
		return e.InternalErrorMessage
	}
	return e.Message
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func New(msg string) error {
	return errors.New(msg)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func (e *Error) WithInternalError(inErr error) *Error {
	e.InternalErrorMessage = inErr.Error()
	return e
}

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	msg := responses.DefaultError

	var (
		customErr *Error
		fiberErr  *fiber.Error
	)

	switch {
	// trieve the custom status code if it's an fiber.*Error
	case errors.As(err, &fiberErr):
		msg.Status = fiberErr.Code
		msg.Code = fiberErr.Code
		msg.Message = fiberErr.Message
		// TODO: handle fiber errors
	case errors.As(err, &customErr):
		msg.Status = customErr.Status
		msg.Code = customErr.Code
		msg.Message = customErr.Message
		msg.ErrorCode = customErr.ErrorCode
	default:
		msg.Status = http.StatusInternalServerError
		msg.Code = enum.INTERNAL
		msg.Message = "Internal server error"
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return msg.JSON(ctx)
}
