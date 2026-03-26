package models

import (
	"github.com/gofiber/fiber/v3"
)

type Registrable interface {
	// GetMethod can be null if we are using a middleware
	GetMethod() HttpMethod
	// GetPath  returns the path of the route that the handler will be registered to
	GetPath() string
	// GetRequestModel returns the model that will be used to parse the request body.
	GetRequestModel() interface{}
	// Handle is the call function that will be registered for the mentioned path and method.
	Handle(ctx fiber.Ctx) error
}
