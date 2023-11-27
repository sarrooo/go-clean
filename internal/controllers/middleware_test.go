package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sarrooo/go-clean/internal/errcode"
	"github.com/sarrooo/go-clean/internal/viewmodel"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	router = NewRouter(zap.NewNop(), nil)
)

func TestRequestViewmodelMiddleware(t *testing.T) {
	tests := map[string]struct {
		method            string
		contentType       string
		requestBody       string
		requestParams     string
		paramsViewmodel   interface{}
		expectedViewmodel interface{}
		expectedError     error
	}{
		"POST Valid Request": {
			method:          "POST",
			contentType:     "application/json",
			requestBody:     `{"field": "value"}`,
			paramsViewmodel: &viewmodel.TestBodyViewModelRequest{},
			expectedViewmodel: &viewmodel.TestBodyViewModelRequest{
				Body: struct {
					Field string `json:"field" binding:"required,min=1"`
				}{
					Field: "value",
				},
			},
			expectedError: nil,
		},
		"POST Invalid JSON": {
			method:            "POST",
			contentType:       "application/json",
			paramsViewmodel:   &viewmodel.TestBodyViewModelRequest{},
			requestBody:       `{"field": "value"`,
			expectedViewmodel: nil,
			expectedError:     errcode.ErrInvalidParameters,
		},
		"POST Error in ShouldBind": {
			method:            "POST",
			contentType:       "application/json",
			paramsViewmodel:   &viewmodel.TestBodyViewModelRequest{},
			requestBody:       `{"invalid_field": "value"}`,
			expectedViewmodel: nil,
			expectedError:     errcode.ErrInvalidParameters,
		},
		"POST Validation Error": {
			method:            "POST",
			contentType:       "application/json",
			paramsViewmodel:   &viewmodel.TestBodyViewModelRequest{},
			requestBody:       `{"field": ""}`,
			expectedViewmodel: nil,
			expectedError:     errcode.ErrInvalidParameters,
		},
		"GET Valid Request": {
			method:          "GET",
			contentType:     "",
			requestBody:     "",
			requestParams:   "?field=value",
			paramsViewmodel: &viewmodel.TestUrlParamsViewModelRequest{},
			expectedViewmodel: &viewmodel.TestUrlParamsViewModelRequest{
				Field: "value",
			},
			expectedError: nil,
		},
		"GET No Required": {
			method:            "GET",
			contentType:       "",
			requestBody:       "",
			requestParams:     "?test=value",
			paramsViewmodel:   &viewmodel.TestUrlParamsViewModelRequest{},
			expectedViewmodel: nil,
			expectedError:     errcode.ErrInvalidParameters,
		},
	}

	for key, test := range tests {
		t.Run(key, func(t *testing.T) {
			// Setup Gin context
			ctx, _ := setupGinContext(test.method, "/"+test.requestParams, test.requestBody, test.contentType)

			// Call middleware
			requestViewmodelMiddleware(test.paramsViewmodel)(ctx)

			// Check if the viewmodel was bound correctly
			boundViewmodel := ctx.Value(ContextKeyRequestViewmodel)
			assert.Equal(t, test.expectedViewmodel, boundViewmodel)

			// Check if the error was handled correctly
			if test.expectedError != nil {
				assert.Error(t, ctx.Errors.Last().Err, test.expectedError)
			} else {
				assert.Empty(t, ctx.Errors)
			}
		})
	}
}

func TestResponseViewmodelMiddleware(t *testing.T) {
	tests := []struct {
		name                   string
		responseViewmodel      interface{}
		statusCode             int
		expectedResponseBody   interface{}
		expectedResponseStatus int
	}{
		{
			name: "Response Viewmodel Present",
			responseViewmodel: &viewmodel.TestViewModelResponse{
				Body: struct {
					Field string `json:"field"`
				}{
					Field: "value",
				},
			},
			statusCode:             http.StatusOK,
			expectedResponseBody:   map[string]interface{}{"field": "value"},
			expectedResponseStatus: http.StatusOK,
		},
		{
			name:                   "Response Viewmodel Present - No Body Field",
			responseViewmodel:      &viewmodel.TestNoBodyViewModelResponse{},
			statusCode:             http.StatusOK,
			expectedResponseBody:   nil, // No JSON response expected
			expectedResponseStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, recorder := setupGinContext(http.MethodGet, "/", "", "")

			// Set status code in context
			ctx.Set(ContextKeyStatusCode, test.statusCode)

			// Set response viewmodel in context
			ctx.Set(ContextKeyResponseViewmodel, test.responseViewmodel)

			router.responseViewmodelMiddleware()(ctx)

			assert.Equal(t, test.expectedResponseStatus, recorder.Code)

			if test.expectedResponseBody != nil {
				expectedJSON, _ := json.Marshal(test.expectedResponseBody)
				assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
			}
		})
	}
}

func TestErrorHandlerMiddleware(t *testing.T) {
	tests := map[string]struct {
		err             error
		failedFields    map[string]string
		expectedStatus  int
		expectedMessage string
		expectedContext map[string]string
	}{
		"GoCleanError with InvalidParameters": {
			err:             errcode.ErrInvalidParameters,
			failedFields:    map[string]string{"field1": "error message 1", "field2": "error message 2"},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "invalid parameters",
			expectedContext: map[string]string{
				"field1": "error message 1",
				"field2": "error message 2",
			},
		},
		"GoCleanError without InvalidParameters": {
			err:             errcode.ErrNotFound,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "not found",
			expectedContext: nil,
		},
		"Non-GoCleanError": {
			err:             errors.New("generic error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "internal error",
			expectedContext: nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Setup Gin context
			ctx, _ := setupGinContext(http.MethodGet, "/", "", "")

			// Set the error in context
			ctx.Error(test.err)
			ctx.Set(ContextKeyInvalidFields, test.failedFields)

			// Call middleware
			router.errorHandlerMiddleware()(ctx)

			// Check if the response viewmodel and status code match the expected values
			responseViewmodel := ctx.Value(ContextKeyResponseViewmodel)
			assert.NotNil(t, responseViewmodel)

			if response, ok := responseViewmodel.(*viewmodel.BadRequestErrorResponse); ok {
				assert.Equal(t, test.expectedStatus, ctx.GetInt(ContextKeyStatusCode))
				assert.Equal(t, test.expectedMessage, response.Body.Message)
				assert.Equal(t, test.expectedContext, response.Body.Context)
			} else if response, ok := responseViewmodel.(*viewmodel.InternalServerErrorResponse); ok {
				assert.Equal(t, test.expectedStatus, ctx.GetInt(ContextKeyStatusCode))
				assert.Equal(t, test.expectedMessage, response.Body.Message)
			} else {
				t.Errorf("Unexpected response viewmodel type")
			}
		})
	}
}

func TestCorsMiddleware(t *testing.T) {
	tests := map[string]struct {
		method       string
		expectedCode int
	}{
		"OPTIONS Request": {
			method:       "OPTIONS",
			expectedCode: http.StatusOK,
		},
		"GET Request": {
			method:       "GET",
			expectedCode: http.StatusOK,
		},
		"POST Request": {
			method:       "POST",
			expectedCode: http.StatusOK,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Setup Gin context
			ctx, recorder := setupGinContext(test.method, "/", "", "")

			// Call middleware
			router.corsMiddleware()(ctx)

			// Check if the response status code matches the expected value
			assert.Equal(t, test.expectedCode, recorder.Code)
		})
	}
}

func TestHandleLanguageMiddleware(t *testing.T) {
	tests := map[string]struct {
		acceptLanguage string
		expectedLocale string
	}{
		"No Accept-Language Header": {
			acceptLanguage: "",
			expectedLocale: "en",
		},
		"Accept-Language Header Set": {
			acceptLanguage: "fr",
			expectedLocale: "fr",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Setup Gin context
			ctx, _ := setupGinContext("GET", "/", "", "")
			ctx.Request.Header.Set("Accept-Language", test.acceptLanguage)

			// Call middleware
			router.handleLanguageMiddleware()(ctx)

			// Check if the locale in the context matches the expected value
			assert.Equal(t, test.expectedLocale, ctx.GetString(ContextKeyLocale))
		})
	}
}

func setupGinContext(method, url, requestBody string, contentType string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = httptest.NewRequest(method, url, strings.NewReader(requestBody))
	ctx.Request.Header.Set("Content-Type", contentType)

	return ctx, recorder
}
