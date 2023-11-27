package viewmodel

import (
	"github.com/sarrooo/go-clean/internal/dto"
)

// swagger:parameters registerController
type RegisterUserRequest struct {
	// in:body
	Body dto.RegisterUser `json:"body" binding:"required"`
}

// swagger:response registerController
type RegisterUserResponse struct {
	// in:body
	Body struct {
		// The access token.
		// Required: true
		Token string `json:"token"`
	} `json:"body"`
}

// swagger:parameters loginController
type LoginUserRequest struct {
	// in:body
	Body struct {
		// The email of the user.
		// Required: true
		Email string `json:"email" binding:"required,email"`

		// The password of the user.
		// Required: true
		Password string `json:"password" binding:"required,min=8,max=64"`
	} `json:"body" binding:"required"`
}

// swagger:response loginController
type LoginUserResponse struct {
	// in:body
	Body struct {
		// The access token.
		// Required: true
		Token string `json:"token"`
	} `json:"body"`
}
