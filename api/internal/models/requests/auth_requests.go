package requests

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,alphaspace"`
	Surname  string `json:"surname" validate:"required,alphaspace"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
