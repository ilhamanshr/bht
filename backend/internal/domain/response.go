package domain

// Response is the standard success response wrapper
type Response struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is the standard error response wrapper
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"error message"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   message,
	}
}
