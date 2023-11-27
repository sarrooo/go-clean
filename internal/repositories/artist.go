package repositories

import (
	"github.com/sarrooo/go-clean/internal/models"
	"gorm.io/gorm"
)

type ArtistRepositoryInterface interface {
	GetByID(id uint) (*models.Artist, error)
	Create(artist *models.Artist) error
	Update(artist *models.Artist) error
	Delete(id uint) error
}

type ArtistRepository struct {
	DB *gorm.DB
}

func (rpt *ArtistRepository) GetByID(id uint) (*models.Artist, error) {
	var artist models.Artist
	err := rpt.DB.Where("id = ?", id).First(&artist).Error
	if err != nil {
		return nil, err
	}
	return &artist, nil
}

func (rpt *ArtistRepository) Create(artist *models.Artist) error {
	return rpt.DB.Create(artist).Error
}

func (rpt *ArtistRepository) Update(artist *models.Artist) error {
	return rpt.DB.UpdateColumns(artist).Error
}

func (rpt *ArtistRepository) Delete(id uint) error {
	return rpt.DB.Delete(&models.Artist{}, id).Error
}
