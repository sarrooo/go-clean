package services

import (
	"reflect"
	"testing"

	"github.com/sarrooo/go-clean/internal/repositories"
	"github.com/sarrooo/go-clean/mocks"
)

type GlobalRepositoryMocks struct {
	User   *mocks.UserRepositoryInterface
	Artist *mocks.ArtistRepositoryInterface

	// Add new repository here
}

// Create new GlobalRepository with all mocks
func newGlobalRepositoryTesting() *repositories.GlobalRepository {
	gr := &repositories.GlobalRepository{
		User:   &mocks.UserRepositoryInterface{},
		Artist: &mocks.ArtistRepositoryInterface{},

		// Add new repository here
	}
	return gr
}

// Cast GlobalRepository to GlobalRepositoryMocks
// This function allow to access to all mock expectations
func castMockGlobalRepository(gr *repositories.GlobalRepository) *GlobalRepositoryMocks {
	return &GlobalRepositoryMocks{
		User:   gr.User.(*mocks.UserRepositoryInterface),
		Artist: gr.Artist.(*mocks.ArtistRepositoryInterface),

		// Add new repository here
	}
}

// Clear all mock expectations and calls
func (gr *GlobalRepositoryMocks) ResetMockCalls() {
	v := reflect.ValueOf(gr).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanInterface() {
			mockObj := field.Interface()
			mockValue := reflect.ValueOf(mockObj)
			expectedCallsField := mockValue.Elem().FieldByName("ExpectedCalls")
			if expectedCallsField.IsValid() && expectedCallsField.CanSet() {
				expectedCallsField.Set(reflect.Zero(expectedCallsField.Type()))
			}
		}
	}
}

// Assert all mock expectations were met
func (gr *GlobalRepositoryMocks) AssertMockCalls(t *testing.T) {
	v := reflect.ValueOf(gr).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanInterface() {
			mockObj := field.Interface()
			mockValue := reflect.ValueOf(mockObj)
			mockValue.MethodByName("AssertExpectations").Call([]reflect.Value{reflect.ValueOf(t)})
		}
	}
}
