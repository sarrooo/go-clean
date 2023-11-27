package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sarrooo/go-clean/internal/errcode"
	"github.com/sarrooo/go-clean/internal/models"
	"github.com/spf13/viper"
)

func (svc *Service) GenerateToken(user *models.User) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Minute * 24 * 30).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err = token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("%w: %v", errcode.ErrGenerateToken, err)
	}

	return tokenString, nil
}
