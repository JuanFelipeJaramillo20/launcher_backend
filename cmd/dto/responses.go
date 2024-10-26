package dto

// CommonSuccess represents a successful response
// swagger:response CommonSuccess
type CommonSuccess struct {
	// Status of the response
	// required: true
	Status string `json:"status"`

	// Message detailing the success
	// required: true
	Message string `json:"message"`
}

// CommonError represents an error response
// swagger:response CommonError
type CommonError struct {
	// Status of the error
	// required: true
	Status string `json:"status"`

	// Error message
	// required: true
	Message string `json:"message"`
}
