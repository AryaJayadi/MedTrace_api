package response

type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BaseValueResponse[T any] struct {
	Success bool       `json:"success"`
	Value   *T         `json:"value,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

type BaseListResponse[T any] struct {
	Success bool       `json:"success"`
	List    []*T       `json:"list,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}
