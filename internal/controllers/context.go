package controllers

// Define context keys for gin.Context
const (
	// Locale
	ContextKeyLocale = "locale"

	// Response viewmodel
	ContextKeyResponseViewmodel = "response_viewmodel"

	// Error handled
	ContextKeyErrorHandled = "error_handled"

	// Request viewmodel
	ContextKeyRequestViewmodel = "request_viewmodel"

	// Status code
	ContextKeyStatusCode = "status_code"

	// Translator
	ContextKeyTranslator = "translator"

	// Error
	ContextKeyInvalidFields = "invalid_fields"
)
