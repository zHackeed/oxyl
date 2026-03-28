package models

import (
	"github.com/gofiber/fiber/v3"
)

type Registrable interface {
	// GetMethod can be null if we are using a middleware
	GetMethod() HttpMethod
	// GetPath  returns the path of the route that the handler will be registered to
	GetPath() string
	// RequestRequirements returns the requirements for the request, if any.
	RequestRequirements() *RequestRequirements
	// Handle is the call function that will be registered for the mentioned path and method.
	Handle(ctx fiber.Ctx) error
}
