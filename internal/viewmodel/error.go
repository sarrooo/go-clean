package viewmodel

// swagger:response errorResponse
type BadRequestErrorResponse struct {
	// in:body
	Body struct {
		// The error message.
		// Required: true
		Message string `json:"message"`

		// The error context.
		Context map[string]string `json:"context,omitempty"`
	} `json:"body"`
}

// swagger:response errorResponse
type InternalServerErrorResponse struct {
	// in:body
	Body struct {
		// The error message.
		// Required: true
		Message string `json:"message"`
	} `json:"body"`
}
