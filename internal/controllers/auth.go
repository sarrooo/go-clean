package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sarrooo/go-clean/internal/services"
	"github.com/sarrooo/go-clean/internal/viewmodel"
)

func registerAuthRoutes(group *gin.RouterGroup, svc services.ServiceInterface) {
	group.POST("/register", requestViewmodelMiddleware(&viewmodel.RegisterUserRequest{}), registerController(svc))
	group.POST("/login", requestViewmodelMiddleware(&viewmodel.LoginUserRequest{}), loginController(svc))
}

// swagger:route POST /auth/register auth registerController
//
// Endpoint for user registration.
//
// responses:
//
//	200: registerController
//	400: errorResponse
func registerController(svc services.ServiceInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := ctx.MustGet(ContextKeyRequestViewmodel).(*viewmodel.RegisterUserRequest)
		response := &viewmodel.RegisterUserResponse{}

		user, err := svc.RegisterUser(&request.Body)
		if err != nil {
			ctx.Error(err)
			return
		}

		token, err := svc.GenerateToken(user)
		if err != nil {
			ctx.Error(err)
			return
		}

		response.Body.Token = token

		ctx.Set(ContextKeyStatusCode, http.StatusOK)
		ctx.Set(ContextKeyResponseViewmodel, response)
	}
}

// swagger:route POST /auth/login auth loginController
//
// Endpoint for user login.
//
// responses:
//
//	200: loginController
//	400: errorResponse
func loginController(svc services.ServiceInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := ctx.MustGet(ContextKeyRequestViewmodel).(*viewmodel.LoginUserRequest)
		response := &viewmodel.LoginUserResponse{}

		user, err := svc.LoginUser(request.Body.Email, request.Body.Password)
		if err != nil {
			ctx.Error(err)
			return
		}

		token, err := svc.GenerateToken(user)
		if err != nil {
			ctx.Error(err)
			return
		}

		response.Body.Token = token

		ctx.Set(ContextKeyStatusCode, http.StatusOK)
		ctx.Set(ContextKeyResponseViewmodel, response)
	}
}
