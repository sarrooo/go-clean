package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sarrooo/go-clean/internal/services"
	"github.com/sarrooo/go-clean/internal/viewmodel"
)

func registerArtistesRoutes(group *gin.RouterGroup, svc services.ServiceInterface) {
	group.GET("/:id", requestViewmodelMiddleware(&viewmodel.GetArtistRequest{}), getArtistController(svc))
}

// swagger:route GET /artists/{id} artistes getArtistController
//
// Endpoint for getting artist.
//
// responses:
//
//	200: getArtistController
//	400: errorResponse
func getArtistController(svc services.ServiceInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := ctx.MustGet(ContextKeyRequestViewmodel).(*viewmodel.GetArtistRequest)
		response := &viewmodel.GetArtistResponse{}

		_ = request

		ctx.Set(ContextKeyStatusCode, 200)
		ctx.Set(ContextKeyResponseViewmodel, response)
	}
}
