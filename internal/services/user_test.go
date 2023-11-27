package services

import (
	"errors"

	"github.com/sarrooo/go-clean/internal/dto"
	"github.com/sarrooo/go-clean/internal/errcode"
	"github.com/sarrooo/go-clean/internal/models"
	"github.com/stretchr/testify/mock"
)

var (
	sampleDtoUser = &dto.RegisterUser{
		Email:    "email@email.com",
		Password: "password",
	}
	sampleModelUser = &models.User{
		Model: models.Model{
			ID: 1,
		},
		Email:    sampleDtoUser.Email,
		Password: "$2a$10$0BbNi1Zwt/5OecX2gvmbQOz9B2JSgretNowsyzBxiwwj8EAW1aKqq",
	}
)

func (suite *ServiceSuiteTest) TestRegisterUser() {
	type parametersType struct {
		registerUser *dto.RegisterUser
	}

	type expectedType struct {
		user *models.User
		err  error
	}

	tests := map[string]struct {
		setupMock  func()
		parameters parametersType
		expected   expectedType
	}{
		"Success": {
			setupMock: func() {
				suite.globalRepositoryMock.User.On("GetByEmail", sampleDtoUser.Email).Return(&models.User{}, nil)
				suite.globalRepositoryMock.User.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
			},
			parameters: parametersType{
				registerUser: sampleDtoUser,
			},
			expected: expectedType{
				user: sampleModelUser,
				err:  nil,
			},
		},
		"Error in GetByEmail": {
			setupMock: func() {
				suite.globalRepositoryMock.User.On("GetByEmail", sampleDtoUser.Email).Return(nil, errcode.ErrDatabase)
			},
			parameters: parametersType{
				registerUser: sampleDtoUser,
			},
			expected: expectedType{
				user: nil,
				err:  errcode.ErrDatabase,
			},
		},
		"User already exist": {
			setupMock: func() {
				suite.globalRepositoryMock.User.On("GetByEmail", sampleDtoUser.Email).Return(sampleModelUser, nil)
			},
			parameters: parametersType{
				registerUser: sampleDtoUser,
			},
			expected: expectedType{
				user: nil,
				err:  errcode.ErrUserAlreadyExists,
			},
		},
		"Error in Create": {
			setupMock: func() {
				suite.globalRepositoryMock.User.On("GetByEmail", sampleDtoUser.Email).Return(&models.User{}, nil)
				suite.globalRepositoryMock.User.On("Create", mock.AnythingOfType("*models.User")).Return(errcode.ErrDatabase)
			},
			parameters: parametersType{
				registerUser: sampleDtoUser,
			},
			expected: expectedType{
				user: nil,
				err:  errcode.ErrDatabase,
			},
		},
		"Wrong Birth Date format": {
			setupMock: func() {
				suite.globalRepositoryMock.User.On("GetByEmail", sampleDtoUser.Email).Return(&models.User{}, nil)
			},
			parameters: parametersType{
				registerUser: &dto.RegisterUser{
					Email:     sampleDtoUser.Email,
					Password:  sampleDtoUser.Password,
					BirthDate: "wrong date format",
				},
			},
			expected: expectedType{
				user: nil,
				err:  errcode.ErrInvalidParameters,
			},
		},
	}

	for testName, test := range tests {
		suite.Run(testName, func() {
			test.setupMock()

			user, err := suite.svc.RegisterUser(test.parameters.registerUser)

			if test.expected.err != nil {
				suite.Assert().Error(err, "Error should have occurred")
				suite.Assert().True(errors.Is(err, test.expected.err), "Error type should match")
				suite.Assert().Nil(user, "User should be nil")
			} else {
				suite.Assert().NoError(err, "No error should have occurred")
			}
		})
	}
}

func (suite *ServiceSuiteTest) TestLoginUser() {
	type parametersType struct {
		email, password string
	}

	type expectedType struct {
		user *models.User
		err  error
	}

	tests := map[string]struct {
		setupMock  func()
		parameters parametersType
		expected   expectedType
	}{
		"Success": {
			setupMock: func() {
				suite.globalRepositoryMock.User.On("GetByEmail", sampleDtoUser.Email).Return(sampleModelUser, nil)
			},
			parameters: parametersType{
				email:    sampleDtoUser.Email,
				password: sampleDtoUser.Password,
			},
			expected: expectedType{
				user: sampleModelUser,
				err:  nil,
			},
		},
		"Error in GetByEmail": {
			setupMock: func() {
				suite.globalRepositoryMock.User.On("GetByEmail", sampleDtoUser.Email).Return(nil, errcode.ErrDatabase)
			},
			parameters: parametersType{
				email:    sampleDtoUser.Email,
				password: sampleDtoUser.Password,
			},
			expected: expectedType{
				user: nil,
				err:  errcode.ErrDatabase,
			},
		},
		"User not exist": {
			setupMock: func() {
				suite.globalRepositoryMock.User.On("GetByEmail", sampleDtoUser.Email).Return(nil, nil)
			},
			parameters: parametersType{
				email:    sampleDtoUser.Email,
				password: sampleDtoUser.Password,
			},
			expected: expectedType{
				user: nil,
				err:  errcode.ErrInvalidCredentials,
			},
		},
		"Wrong Password": {
			setupMock: func() {
				suite.globalRepositoryMock.User.On("GetByEmail", sampleDtoUser.Email).Return(sampleModelUser, nil)
			},
			parameters: parametersType{
				email:    sampleDtoUser.Email,
				password: "wrong password",
			},
			expected: expectedType{
				user: nil,
				err:  errcode.ErrInvalidCredentials,
			},
		},
	}

	for testName, test := range tests {
		suite.Run(testName, func() {
			test.setupMock()

			user, err := suite.svc.LoginUser(test.parameters.email, test.parameters.password)

			if test.expected.err != nil {
				suite.Assert().Error(err, "Error should have occurred")
				suite.Assert().True(errors.Is(err, test.expected.err), "Error type should match")
				suite.Assert().Nil(user, "User should be nil")
			} else {
				suite.Assert().NoError(err, "No error should have occurred")
			}
		})
	}
}
