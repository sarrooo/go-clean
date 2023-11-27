package services

import (
	"fmt"
	"testing"

	"github.com/sarrooo/go-clean/internal/logger"
	"github.com/stretchr/testify/suite"
)

type ServiceSuiteTest struct {
	suite.Suite
	svc *Service

	// Mocks
	globalRepositoryMock *GlobalRepositoryMocks
}

func (suite *ServiceSuiteTest) SetupSuite() {
	logger, err := logger.New()
	if err != nil {
		panic(fmt.Errorf("error create logger: %w", err))
	}

	globalRpt := newGlobalRepositoryTesting()
	suite.globalRepositoryMock = castMockGlobalRepository(globalRpt)

	suite.svc = New(logger, globalRpt)
}

func (suite *ServiceSuiteTest) SetupSubTest() {
	suite.globalRepositoryMock.ResetMockCalls()
}

func (suite *ServiceSuiteTest) TearDownSubTest() {
	suite.globalRepositoryMock.AssertMockCalls(suite.T())
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuiteTest))
}
