package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sarrooo/go-clean/internal/models"
	"github.com/sarrooo/go-clean/internal/services"
	"github.com/sarrooo/go-clean/internal/viewmodel"
)

func registerArtistesRoutes(group *gin.RouterGroup, svc services.ServiceInterface) {
	group.POST("/", requestViewmodelMiddleware(&viewmodel.CreateArtistResponse{}), createArtistController(svc))
	group.GET("/:id", requestViewmodelMiddleware(&viewmodel.GetArtistRequest{}), getArtistController(svc))
	group.DELETE("/:id", requestViewmodelMiddleware(&viewmodel.DeleteArtistRequest{}), deleteArtistController(svc))
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

		artist, err := svc.GetArtist(request.ID)
		if err != nil {
			ctx.Error(err)
			return
		}

		response.Body.ID = artist.ID
		response.Body.Name = artist.Name

		ctx.Set(ContextKeyStatusCode, 200)
		ctx.Set(ContextKeyResponseViewmodel, response)
	}
}

// swagger:route POST /artists artistes createArtistController
//
// Endpoint for creating artist.
//
// responses:
//
//	200: createArtistController
//	400: errorResponse
func createArtistController(svc services.ServiceInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := ctx.MustGet(ContextKeyRequestViewmodel).(*viewmodel.CreateArtistRequest)
		response := &viewmodel.CreateArtistResponse{}

		artist := &models.Artist{
			Name: request.Name,
		}
		err := svc.CreateArtist(artist)
		if err != nil {
			ctx.Error(err)
			return
		}

		response.Body.ID = artist.ID

		ctx.Set(ContextKeyStatusCode, 200)
		ctx.Set(ContextKeyResponseViewmodel, response)
	}
}

// swagger:route DELETE /artists/{id} artistes deleteArtistController
//
// Endpoint for deleting artist.
//
// responses:
//
//	200: deleteArtistController
//	400: errorResponse
func deleteArtistController(svc services.ServiceInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := ctx.MustGet(ContextKeyRequestViewmodel).(*viewmodel.DeleteArtistRequest)
		response := &viewmodel.DeleteArtistResponse{}

		err := svc.DeleteArtist(request.ID)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.Set(ContextKeyStatusCode, 200)
		ctx.Set(ContextKeyResponseViewmodel, response)
	}
}
