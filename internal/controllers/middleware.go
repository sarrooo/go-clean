package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	fr_translations "github.com/go-playground/validator/v10/translations/fr"
	"github.com/sarrooo/go-clean/internal/errcode"
	"github.com/sarrooo/go-clean/internal/viewmodel"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// Bind request view model and pass it to the next handler
// It use the gin method to bind, please check `ShouldBind` documentation
// WARNING : It not handle URI parameters
func requestViewmodelMiddleware(requestViewmodel interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Read the body and rewrite it with body key
		// It's because we use go-swagger annotation
		// and the body annotation must be in Body struct
		requestBody, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.Error(fmt.Errorf("%w: %v", errcode.ErrInvalidParameters, err))
			ctx.Abort()
			return
		}
		transformedBody := []byte(`{"body":` + string(requestBody) + `}`)
		ctx.Request.Body = io.NopCloser(bytes.NewReader(transformedBody))

		// Create a new instance of requestViewmodel and bind it
		requestViewmodelInstance := reflect.New(reflect.TypeOf(requestViewmodel).Elem()).Interface()

		// Bind URI tagged fields
		// These will be validated by the validator in Gin binding methods
		err = bindURITaggedFields(ctx, requestViewmodelInstance)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		// Use Gin binding methods
		// It choose the binding method according to the method and "Content-Type" header
		if err := ctx.ShouldBind(requestViewmodelInstance); err != nil {
			// Rewrite the body with the original body
			ctx.Request.Body = io.NopCloser(bytes.NewReader(requestBody))

			// Check if the error is a validation error
			var verr validator.ValidationErrors
			if errors.As(err, &verr) {
				// Get the translator from the context and translate the errors messages
				translatedErrors := make(map[string]string, len(verr))
				trans, exists := ctx.Get(ContextKeyTranslator)
				if exists {
					trans := trans.(ut.Translator)
					translatedErrors = verr.Translate(trans)
				}

				// Create a map to hold the failed fields
				failedFields := make(map[string]string, len(verr))
				for _, fieldError := range verr {
					// Get the translated error message
					// If there is no translation, use the default error message
					validationMessage := translatedErrors[fieldError.Namespace()]
					if validationMessage == "" {
						validationMessage = fieldError.Error()
					}
					failedFields[fieldError.Field()] = validationMessage
				}
				// Set the failed fields in the context, it will be used by the error handler middleware
				ctx.Set(ContextKeyInvalidFields, failedFields)
			}
			err = fmt.Errorf("%w: %v", errcode.ErrInvalidParameters, err)
			ctx.Error(err)
			ctx.Abort()
			return
		}
		// Rewrite the body with the original body
		ctx.Request.Body = io.NopCloser(bytes.NewReader(requestBody))

		// Set binding result in Gin context
		ctx.Set(ContextKeyRequestViewmodel, requestViewmodelInstance)
		ctx.Next()
	}
}

// Bind URI tagged fields
// It use the Gin params to get the URI parameters
func bindURITaggedFields(ctx *gin.Context, data interface{}) error {
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := typ.Field(i).Tag.Get("uri")

		if tag == "" {
			continue
		}

		paramValue, exists := ctx.Params.Get(tag)
		if !exists {
			continue
		}

		// You might need additional type conversion logic based on the actual field type
		switch field.Kind() {
		case reflect.String:
			field.SetString(paramValue)
		case reflect.Uint:
			// Assuming the parameter is an uinteger
			paramInt, err := strconv.ParseUint(paramValue, 10, 32)
			if err != nil {
				return errcode.ErrInvalidParameters
			}
			field.SetUint(paramInt)
		case reflect.Int:
			// Assuming the parameter is an integer
			paramInt, err := strconv.Atoi(paramValue)
			if err != nil {
				return errcode.ErrInvalidParameters
			}
			field.SetInt(int64(paramInt))
		default:
			continue
		}
	}
	return nil
}

// This middleware get the response view model from the Gin context and send it
func (rtr *Router) responseViewmodelMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		// Check if there is a responseViewmodel in the context
		if responseViewmodel, exist := ctx.Get(ContextKeyResponseViewmodel); exist {
			statusCode := ctx.GetInt(ContextKeyStatusCode)

			value := reflect.ValueOf(responseViewmodel)
			if value.Kind() == reflect.Ptr && !value.IsNil() {
				// Check if responseViewmodel has a Body field
				bodyField := value.Elem().FieldByName("Body")
				if bodyField.IsValid() {
					ctx.JSON(statusCode, bodyField.Interface())
					return
				}
				// If responseViewmodel has no Body field, send just the status code
				ctx.Status(statusCode)
			}
		}
		ctx.Status(http.StatusBadRequest)
	}
}

// It the centralized error handling middleware
// When an error occured in handler, just set the error with `c.Error(...)` and return
// The middleware will handle error, log it, answer to the client
func (rtr *Router) errorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		err := ctx.Errors.Last()
		if err != nil {
			body, errReading := io.ReadAll(ctx.Request.Body)
			if errReading != nil {
				body = []byte("Error reading body")
			}

			rtr.logger.Error("middleware error",
				zap.Error(err.Err),
				zap.String("path", ctx.Request.URL.Path),
				zap.String("method", ctx.Request.Method),
				zap.String("ip", ctx.ClientIP()),
				zap.String("body", string(body)))

			// check if the error contain a GoCleanError
			// if yes, return the first GoCleanError to the client
			// if no, return a generic error
			// This strategy limit the informations given to the client
			var GoCleanError errcode.GoCleanError
			if errors.As(err.Err, &GoCleanError) {
				// Define the response viewmodel
				response := &viewmodel.BadRequestErrorResponse{}
				response.Body.Message = GoCleanError.Error()

				// Check if the error is an invalid parameters error
				if errors.Is(GoCleanError, errcode.ErrInvalidParameters) {
					// Extract the failed fields from the error message
					failedFields := ctx.GetStringMapString(ContextKeyInvalidFields)
					if len(failedFields) != 0 {
						response.Body.Context = failedFields
					}
				}
				ctx.Set(ContextKeyStatusCode, http.StatusBadRequest)
				ctx.Set(ContextKeyResponseViewmodel, response)
				return
			}

			// If the error is not a GoCleanError, return a generic error
			response := &viewmodel.InternalServerErrorResponse{}
			response.Body.Message = "internal error"
			ctx.Set(ContextKeyStatusCode, http.StatusInternalServerError)
			ctx.Set(ContextKeyResponseViewmodel, response)
			return
		}
	}
}

// Set the CORS rules
func (rtr *Router) corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().
			Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		ctx.Writer.Header().
			Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Accept, Authorization, Two-Factor-Code, Recaptcha, Lang, Country, Session-Id, Api-Key")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	}
}

func (rtr *Router) handleLanguageMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		locale := ctx.GetHeader("Accept-Language")

		tag, _ := language.MatchStrings(rtr.languageMatcher, locale)
		base, _ := tag.Base()
		locale = base.String()

		trans, _ := rtr.universalTranslator.GetTranslator(locale)

		// Register validator translations for the current locale
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			switch locale {
			case "fr":
				fr_translations.RegisterDefaultTranslations(v, trans)
			default:
				en_translations.RegisterDefaultTranslations(v, trans)
			}
		}

		// Set the language in the context
		ctx.Set(ContextKeyLocale, locale)
		ctx.Set(ContextKeyTranslator, trans)
		ctx.Next()
	}
}
