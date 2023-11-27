package services

import (
	"github.com/sarrooo/go-clean/internal/dto"
	"github.com/sarrooo/go-clean/internal/models"
	"github.com/sarrooo/go-clean/internal/repositories"
	"go.uber.org/zap"
)

// TODO: Can we find a way to split this object into different object to avoid having 1000+ methods on the same space
type ServiceInterface interface {
	RegisterUser(registerUser *dto.RegisterUser) (user *models.User, err error)
	LoginUser(email, password string) (user *models.User, err error)

	GenerateToken(user *models.User) (tokenString string, err error)
}

type Service struct {
	logger           *zap.Logger
	globalRepository *repositories.GlobalRepository
}

func New(
	logger *zap.Logger,
	globalRepository *repositories.GlobalRepository,
) *Service {
	service := &Service{
		logger,
		globalRepository,
	}
	return service
}
