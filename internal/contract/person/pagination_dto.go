package contract

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data" swaggertype:"array,object"`          // List of items in current page
	Page       int         `json:"page" example:"1"`                         // Current page number
	PageSize   int         `json:"page_size" example:"10"`                   // Number of items per page
	TotalItems int64       `json:"total_items" example:"100"`                // Total number of items in database
	TotalPages int         `json:"total_pages" example:"10"`                 // Total number of pages available
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error" example:"validation_error"`                      // Error code
	Message string `json:"message" example:"Invalid CPF"`                         // Descriptive error message
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	ID      int    `json:"id" example:"1"`                                        // ID of created resource
	Message string `json:"message" example:"Person created successfully"`          // Success message
}
