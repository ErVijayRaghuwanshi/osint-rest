package common

// PingResponse is the standard response for platform ping endpoints.
type PingResponse struct {
	Message string `json:"message" example:"Pong"`
}

// BadRequestError is returned when the request is malformed or missing required parameters.
type BadRequestError struct {
	Error string `json:"error" example:"username parameter is required"`
}

// NotFoundError is returned when the requested resource is not found.
type NotFoundError struct {
	Error string `json:"error" example:"user not found"`
}

// InternalServerError is returned when an unexpected server error occurs.
type InternalServerError struct {
	Error string `json:"error" example:"internal server error"`
}
