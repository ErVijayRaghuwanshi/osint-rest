package common

import "github.com/gin-gonic/gin"

// SetCacheHeader sets the X-Cache response header to HIT or MISS.
func SetCacheHeader(c *gin.Context, hit bool) {
	if hit {
		c.Header("X-Cache", "HIT")
	} else {
		c.Header("X-Cache", "MISS")
	}
}

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
