package controllers

import (
	"net/http"

	"github.com/sarrooo/go-clean/internal/dto"
	"github.com/sarrooo/go-clean/internal/errcode"
	"github.com/sarrooo/go-clean/internal/models"
	"github.com/sarrooo/go-clean/internal/viewmodel"
)

var (
	sampleDtoUser = dto.RegisterUser{
		Email:     "user@gmail.com",
		Password:  "123456pass",
		FirstName: "user",
		LastName:  "user",
		Phone:     "+33606060606",
		BirthDate: "1990-01-01",
	}

	sampleModelUser = &models.User{}
)

func (suite *ControllerSuiteTest) TestRegisterController() {
	tests := controllerTestTable{
		"Success": {
			setupMock: func() {
				suite.svc.On("RegisterUser", &sampleDtoUser).Return(sampleModelUser, nil)
				suite.svc.On("GenerateToken", sampleModelUser).Return("token", nil)
				suite.svc.On("GenerateRefreshToken", sampleModelUser).Return("refreshToken", nil)
			},
			requestViewmodel: &viewmodel.RegisterUserRequest{
				Body: sampleDtoUser,
			},
			expected: controllerTestExpected{
				status: http.StatusOK,
				responseViewmodel: &viewmodel.RegisterUserResponse{
					Body: struct {
						Token string "json:\"token\""
					}{
						Token: "token",
					},
				},
			},
		},
		"Error from RegisterUser": {
			setupMock: func() {
				suite.svc.On("RegisterUser", &sampleDtoUser).Return(nil, errcode.ErrDatabase)
			},
			requestViewmodel: &viewmodel.RegisterUserRequest{
				Body: sampleDtoUser,
			},
			expected: controllerTestExpected{isError: true},
		},
		"Error from GenerateToken": {
			setupMock: func() {
				suite.svc.On("RegisterUser", &sampleDtoUser).Return(sampleModelUser, nil)
				suite.svc.On("GenerateToken", sampleModelUser).Return("", errcode.ErrGenerateToken)
			},
			requestViewmodel: &viewmodel.RegisterUserRequest{
				Body: sampleDtoUser,
			},
			expected: controllerTestExpected{isError: true},
		},
		"Error from GenerateRefreshToken": {
			setupMock: func() {
				suite.svc.On("RegisterUser", &sampleDtoUser).Return(sampleModelUser, nil)
				suite.svc.On("GenerateToken", sampleModelUser).Return("token", nil)
				suite.svc.On("GenerateRefreshToken", sampleModelUser).Return("", errcode.ErrGenerateToken)
			},
			requestViewmodel: &viewmodel.RegisterUserRequest{
				Body: sampleDtoUser,
			},
			expected: controllerTestExpected{isError: true},
		},
	}

	suite.executeTestTable(tests, registerController)
}

func (suite *ControllerSuiteTest) TestLoginController() {
	tests := controllerTestTable{
		"Success": {
			setupMock: func() {
				suite.svc.On("LoginUser", sampleDtoUser.Email, sampleDtoUser.Password).Return(sampleModelUser, nil)
				suite.svc.On("GenerateToken", sampleModelUser).Return("token", nil)
				suite.svc.On("GenerateRefreshToken", sampleModelUser).Return("refreshToken", nil)
			},
			requestViewmodel: &viewmodel.LoginUserRequest{
				Body: struct {
					Email    string "json:\"email\" binding:\"required,email\""
					Password string "json:\"password\" binding:\"required,min=8,max=64\""
				}{
					Email:    sampleDtoUser.Email,
					Password: sampleDtoUser.Password,
				},
			},
			expected: controllerTestExpected{
				status: http.StatusOK,
				responseViewmodel: &viewmodel.LoginUserResponse{
					Body: struct {
						Token string "json:\"token\""
					}{
						Token: "token",
					},
				},
			},
		},
		"Error from LoginUser": {
			setupMock: func() {
				suite.svc.On("LoginUser", sampleDtoUser.Email, sampleDtoUser.Password).Return(nil, errcode.ErrDatabase)
			},
			requestViewmodel: &viewmodel.LoginUserRequest{
				Body: struct {
					Email    string "json:\"email\" binding:\"required,email\""
					Password string "json:\"password\" binding:\"required,min=8,max=64\""
				}{
					Email:    sampleDtoUser.Email,
					Password: sampleDtoUser.Password,
				},
			},
			expected: controllerTestExpected{isError: true},
		},
		"Error from GenerateToken": {
			setupMock: func() {
				suite.svc.On("LoginUser", sampleDtoUser.Email, sampleDtoUser.Password).Return(sampleModelUser, nil)
				suite.svc.On("GenerateToken", sampleModelUser).Return("", errcode.ErrGenerateToken)
			},
			requestViewmodel: &viewmodel.LoginUserRequest{
				Body: struct {
					Email    string "json:\"email\" binding:\"required,email\""
					Password string "json:\"password\" binding:\"required,min=8,max=64\""
				}{
					Email:    sampleDtoUser.Email,
					Password: sampleDtoUser.Password,
				},
			},
			expected: controllerTestExpected{isError: true},
		},
		"Error from GenerateRefreshToken": {
			setupMock: func() {
				suite.svc.On("LoginUser", sampleDtoUser.Email, sampleDtoUser.Password).Return(sampleModelUser, nil)
				suite.svc.On("GenerateToken", sampleModelUser).Return("token", nil)
				suite.svc.On("GenerateRefreshToken", sampleModelUser).Return("", errcode.ErrGenerateToken)
			},
			requestViewmodel: &viewmodel.LoginUserRequest{
				Body: struct {
					Email    string "json:\"email\" binding:\"required,email\""
					Password string "json:\"password\" binding:\"required,min=8,max=64\""
				}{
					Email:    sampleDtoUser.Email,
					Password: sampleDtoUser.Password,
				},
			},
			expected: controllerTestExpected{isError: true},
		},
	}

	suite.executeTestTable(tests, loginController)
}
