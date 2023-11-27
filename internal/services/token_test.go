package services

import (
	"errors"

	"github.com/sarrooo/go-clean/internal/models"
)

func (suite *ServiceSuiteTest) TestGenerateToken() {
	type parametersType struct {
		user *models.User
	}

	type expectedType struct {
		err error
	}

	// Define sample user
	sampleUser := &models.User{
		Email: "test@example.com",
	}

	tests := map[string]struct {
		parameters parametersType
		expected   expectedType
	}{
		"Success": {
			parameters: parametersType{
				user: sampleUser,
			},
			expected: expectedType{
				err: nil,
			},
		},
	}

	for testName, test := range tests {
		suite.Run(testName, func() {
			tokenString, err := suite.svc.GenerateToken(test.parameters.user)

			if test.expected.err != nil {
				suite.Assert().Error(err, "Error should have occurred")
				suite.Assert().True(errors.Is(err, test.expected.err), "Error type should match")
				suite.Assert().Empty(tokenString, "Token string should be empty")
			} else {
				suite.Assert().NoError(err, "No error should have occurred")
				suite.Assert().NotEmpty(tokenString, "Token string should not be empty")
			}
		})
	}
}
