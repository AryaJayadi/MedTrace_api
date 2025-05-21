package response

import "fmt"

func SuccessValueResponse[T any](value T) BaseValueResponse[T] {
	return BaseValueResponse[T]{
		Success: true,
		Value:   &value,
		Error:   nil,
	}
}

func ErrorValueResponse[T any](code int, format string, args ...any) BaseValueResponse[T] {
	return BaseValueResponse[T]{
		Success: false,
		Value:   nil,
		Error: &ErrorInfo{
			Code:    code,
			Message: fmt.Sprintf(format, args...),
		},
	}
}

func SuccessListResponse[T any](list []*T) BaseListResponse[T] {
	if list == nil {
		list = make([]*T, 0)
	}
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
