package services

import (
	"fmt"

	"github.com/sarrooo/go-clean/internal/errcode"
	"github.com/sarrooo/go-clean/internal/models"
)

func (svc *Service) CreateArtist(artist *models.Artist) (err error) {
	err = svc.globalRepository.Artist.Create(artist)
	if err != nil {
		return fmt.Errorf("%w: %v", errcode.ErrDatabase, err)
	}

	return nil
}

func (svc *Service) GetArtist(id uint) (artist *models.Artist, err error) {
	artist, err = svc.globalRepository.Artist.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errcode.ErrDatabase, err)
	}

	if artist.ID == 0 {
		return nil, fmt.Errorf("%w: %v", errcode.ErrNotFound, err)
	}

	return artist, nil
}

func (svc *Service) DeleteArtist(id uint) (err error) {
	err = svc.globalRepository.Artist.Delete(id)
	if err != nil {
		return fmt.Errorf("%w: %v", errcode.ErrDatabase, err)
	}

	return nil
}
