package model

// ErrorResponse representa um erro retornado pela API
type ErrorResponse struct {
	Message string `json:"message"`
}