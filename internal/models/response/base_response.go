package response

type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BaseValueResponse struct {
	Success bool       `json:"success"`
	Value   any        `json:"value,omitempty"` // T
	Error   *ErrorInfo `json:"error,omitempty"` // optional
}

type BaseListResponse struct {
	Success bool       `json:"success"`
	List    []any      `json:"list,omitempty"`  // []T
	Error   *ErrorInfo `json:"error,omitempty"` // optional
}
