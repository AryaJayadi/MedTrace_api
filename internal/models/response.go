package models

func SuccessValueResponse(value any) BaseValueResponse {
	return BaseValueResponse{
		Success: true,
		Value:   value,
		Error:   nil,
	}
}

func ErrorValueResponse(code int, msg string) BaseValueResponse {
	return BaseValueResponse{
		Success: false,
		Value:   nil,
		Error:   &ErrorInfo{Code: code, Message: msg},
	}
}

func SuccessListResponse(list []any) BaseListResponse {
	return BaseListResponse{
		Success: true,
		List:    list,
		Error:   nil,
	}
}

func ErrorListResponse(code int, msg string) BaseListResponse {
	return BaseListResponse{
		Success: false,
		List:    nil,
		Error:   &ErrorInfo{Code: code, Message: msg},
	}
}
