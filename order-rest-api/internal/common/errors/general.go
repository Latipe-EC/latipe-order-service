package errors

import (
	"order-rest-api/internal/domain/enum"
)

var (
	ErrInternalServer = &Error{
		Status:    500,
		Code:      enum.INTERNAL,
		ErrorCode: "GE001",
		Message:   "Internal server error",
	}

	ErrBadRequest = &Error{
		Status:    400,
		Code:      enum.INVALID_ARGUMENT,
		ErrorCode: "GE002",
		Message:   "Bad request",
	}

	ErrPermissionDenied = &Error{
		Status:    403,
		Code:      enum.PERMISSION_DENIED,
		ErrorCode: "GE003",
		Message:   "Permission denied",
	}

	ErrNotFound = &Error{
		Status:    404,
		Code:      enum.NOT_FOUND,
		ErrorCode: "GE004",
		Message:   "Not found",
	}

	ErrAlreadyExists = &Error{
		Status:    409,
		Code:      enum.ALREADY_EXISTS,
		ErrorCode: "GE005",
		Message:   "Already exists",
	}

	ErrUnauthenticated = &Error{
		Status:    401,
		Code:      enum.UNAUTHENTICATED,
		ErrorCode: "GE006",
		Message:   "Unauthorized",
	}

	ErrInvalidCredentials = &Error{
		Status:    401,
		Code:      enum.UNAUTHENTICATED,
		ErrorCode: "GE007",
		Message:   "Invalid credentials",
	}

	ErrNotFoundRecord = &Error{
		Status:    404,
		Code:      enum.NOT_FOUND,
		ErrorCode: "GE008",
		Message:   "Record does not exist",
	}

	ErrInvalidParameters = &Error{
		Status:    400,
		Code:      enum.FAILED_PRECONDITION,
		ErrorCode: "GE009",
		Message:   "Invalid parameters",
	}

	ErrTooManyRequest = &Error{
		Status:    429,
		Code:      enum.RESOURCE_EXHAUSTED,
		ErrorCode: "GE010",
		Message:   "Too Many Requests",
	}

	ErrInvalidFilter = &Error{
		Status:    400,
		Code:      enum.INVALID_ARGUMENT,
		ErrorCode: "GE011",
		Message:   "Invalid filters",
	}
)
