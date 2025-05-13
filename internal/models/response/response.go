package response

import "fmt"

func SuccessValueResponse(value any) BaseValueResponse {
	return BaseValueResponse{
		Success: true,
		Value:   value,
		Error:   nil,
	}
}

func ErrorValueResponse(code int, format string, args ...any) BaseValueResponse {
	return BaseValueResponse{
		Success: false,
		Value:   nil,
		Error: &ErrorInfo{
			Code:    code,
			Message: fmt.Sprintf(format, args...),
		},
	}
}

func SuccessListResponse[T any](list []T) BaseListResponse[T] {
	return BaseListResponse[T]{
		Success: true,
		List:    list,
		Error:   nil,
	}
}

func ErrorListResponse[T any](code int, format string, args ...any) BaseListResponse[T] {
	return BaseListResponse[T]{
		Success: false,
		List:    nil,
		Error: &ErrorInfo{
			Code:    code,
			Message: fmt.Sprintf(format, args...),
		},
	}
}
