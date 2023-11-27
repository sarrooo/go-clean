package repositories

import (
	"gorm.io/gorm"
)

type GlobalRepository struct {
	User   UserRepositoryInterface
	Artist ArtistRepositoryInterface

	// Add new repository here
}

func NewGlobalRepository(DB *gorm.DB) *GlobalRepository {
	gr := &GlobalRepository{
		User:   &UserRepository{DB: DB},
		Artist: &ArtistRepository{DB: DB},

		// Add new repository here
	}
	return gr
}
