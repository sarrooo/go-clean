package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sarrooo/go-clean/internal/services"
	"github.com/sarrooo/go-clean/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ControllerSuiteTest struct {
	suite.Suite
	ctx *gin.Context
	svc *mocks.ServiceInterface
}

type controllerTestTable = map[string]controllerTest

type controllerTest struct {
	setupMock        func()
	requestViewmodel interface{}
	expected         controllerTestExpected
}

type controllerTestExpected struct {
	responseViewmodel interface{}
	status            int
	isError           bool
}

func (suite *ControllerSuiteTest) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// Setup Mocks
	suite.svc = &mocks.ServiceInterface{}
}

func (suite *ControllerSuiteTest) SetupSubTest() {
	// Before each sub test, reset the context
	suite.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())

	// Reset mocks calls
	suite.svc.ExpectedCalls = nil
}

func (suite *ControllerSuiteTest) executeTestTable(tests controllerTestTable, controller func(svc services.ServiceInterface) gin.HandlerFunc) {
	for key, test := range tests {
		suite.Run(key, func() {

			// Setup mocks calls
			test.setupMock()

			// Set request viewmodel
			suite.ctx.Set(ContextKeyRequestViewmodel, test.requestViewmodel)

			// Call controller
			controller(suite.svc)(suite.ctx)

			// Check mocks
			suite.svc.AssertExpectations(suite.T())

			if test.expected.isError {
				if (suite.ctx.Errors.Last() == nil) || (suite.ctx.Errors.Last().Err == nil) {
					assert.Fail(suite.T(), "Error expected")
				}
				return
			}

			// Check response
			responseViewmodel := suite.ctx.MustGet(ContextKeyResponseViewmodel)
			assert.Equal(suite.T(), test.expected.status, suite.ctx.GetInt(ContextKeyStatusCode))
			assert.Equal(suite.T(), test.expected.responseViewmodel, responseViewmodel)
		})
	}
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerSuiteTest))
}
