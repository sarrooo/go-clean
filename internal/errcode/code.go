package errcode

import (
	"errors"
	"fmt"
)

// GoCleanError is the error type used in the project
// It is a wrapper of the standard error type
// There is a goal:
// - The error handle middleware can log the error, and answer just a message to the client
type GoCleanError struct {
	error
	int
}

var (
	//// generic errors (100-199)
	ErrUndefined      = newErrcode("undefined error", 100)
	ErrNotImplemented = newErrcode("not implemented", 101)

	//// database errors (200-299)
	ErrDatabase        = newErrcode("database error", 200)
	ErrDatabaseMigrate = newErrcode("database migrate error", 201)
	ErrDropProduction  = newErrcode("production database cannot be dropped", 202)

	//// controllers errors (300-399)
	ErrInvalidParameters   = newErrcode("invalid parameters", 300)
	ErrNotFound            = newErrcode("not found", 301)
	ErrUnknown             = newErrcode("unknown", 302)
	ErrConfigurationFailed = newErrcode("configuration failed", 303)

	//// auth errors (400-499)
	ErrUnauthorized = newErrcode("unauthorized", 400)
	ErrForbidden    = newErrcode("forbidden", 401)

	//// business logic errors (500-599)
	ErrExternalLib        = newErrcode("external librairie", 500)
	ErrTemplatingEmail    = newErrcode("templating email error", 501)
	ErrSendingEmail       = newErrcode("email error", 501)
	ErrEmail              = newErrcode("email error", 501)
	ErrUserAlreadyExists  = newErrcode("user already exists", 502)
	ErrGenerateToken      = newErrcode("error creating token", 503)
	ErrInvalidToken       = newErrcode("invalid token", 504)
	ErrTokenExpirated     = newErrcode("invalid token", 505)
	ErrInvalidCredentials = newErrcode("invalid credentials", 506)
)

func newErrcode(message string, code int) GoCleanError {
	return GoCleanError{errors.New(message), code}
}

func Wrap(err *error, format string, args ...any) {
	if *err != nil {
		*err = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *err)
	}
}
