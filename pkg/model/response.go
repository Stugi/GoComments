package model

type ErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"` // RequestID будет опциональным
}
