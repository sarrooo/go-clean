package controllers

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/sarrooo/go-clean/internal/services"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

type Router struct {
	engine              *gin.Engine
	logger              *zap.Logger
	languageMatcher     language.Matcher
	universalTranslator *ut.UniversalTranslator
}

func NewRouter(logger *zap.Logger, svc services.ServiceInterface) *Router {
	router := &Router{}

	router.logger = logger

	en := en.New()
	fr := fr.New()
	router.universalTranslator = ut.New(en, en, fr)

	router.languageMatcher = language.NewMatcher([]language.Tag{
		language.English,
		language.French,
	})

	router.engine = gin.Default()
	config(router.engine)

	/* Middleware */
	router.engine.Use(router.corsMiddleware())
	router.engine.Use(router.handleLanguageMiddleware())
	router.engine.Use(router.responseViewmodelMiddleware())
	router.engine.Use(router.errorHandlerMiddleware())

	router.registerRoutes(svc)

	return router
}

func (rtr *Router) Run(addr string) error {
	return rtr.engine.Run(addr)
}

func (rtr *Router) registerRoutes(svc services.ServiceInterface) {
	/* Auth */
	auth := rtr.engine.Group("/auth")
	registerAuthRoutes(auth, svc)

	/* Albums */
	albums := rtr.engine.Group("/albums")
	registerArtistesRoutes(albums, svc)
}

func config(router *gin.Engine) {
	router.RedirectTrailingSlash = false

	// Custom validator
	// This is used to be able to get the json tag name instead of the struct field name
	// to send it to the client, and the client can use it to display the error message
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return fld.Name
			}
			return name
		})
	}
}
