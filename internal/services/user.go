package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sarrooo/go-clean/internal/dto"
	"golang.org/x/crypto/bcrypt"

	"github.com/sarrooo/go-clean/internal/errcode"
	"github.com/sarrooo/go-clean/internal/models"
)

// RegisterUser creates a new user in the database, and sends a welcome email
// If isEmailVerified is true, the user is created with the email verified (by example using oauth)
func (svc *Service) RegisterUser(registerUser *dto.RegisterUser) (user *models.User, err error) {
	user, err = svc.globalRepository.User.GetByEmail(registerUser.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errcode.ErrDatabase, err)
	}

	// if user already exists, return error
	if user.ID != 0 {
		return nil, fmt.Errorf("%w", errcode.ErrUserAlreadyExists)
	}

	user, err = svc.formatRegisterUser(registerUser)
	if err != nil {
		return nil, err
	}

	// create user in database
	err = svc.globalRepository.User.Create(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errcode.ErrDatabase, err)
	}

	return user, nil
}

// LoginUser checks if the user exists and if the password is correct
// If the user exists and the password is correct, it returns the user, otherwise it returns an error
func (svc *Service) LoginUser(email, password string) (user *models.User, err error) {
	// check if email already exists
	user, err = svc.globalRepository.User.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errcode.ErrDatabase, err)
	}

	// if user does not exist, return error
	if user == nil {
		return nil, fmt.Errorf("%w: %v", errcode.ErrInvalidCredentials, errors.New("user does not exist"))
	}

	// check if password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errcode.ErrInvalidCredentials, err)
	}

	return user, nil
}

func (svc *Service) formatRegisterUser(registerUser *dto.RegisterUser) (user *models.User, err error) {
	// format email
	registerUser.Email = strings.ToLower(strings.TrimSpace(registerUser.Email))

	// password hashing
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errcode.ErrExternalLib, err)
	}
	registerUser.Password = string(passwordHash)

	// parse birth date
	var parsedBirthDate time.Time
	if registerUser.BirthDate != "" {
		parsedBirthDate, err = time.Parse("2006-01-02", registerUser.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errcode.ErrInvalidParameters, err)
		}
	}

	user = &models.User{
		Email:     registerUser.Email,
		Password:  registerUser.Password,
		FirstName: registerUser.FirstName,
		LastName:  registerUser.LastName,
		Phone:     registerUser.Phone,
		BirthDate: parsedBirthDate,
	}
	return user, nil
}
