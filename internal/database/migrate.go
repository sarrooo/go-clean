package database

import (
	"github.com/sarrooo/go-clean/internal/models"
	"gorm.io/gorm"
)

func migrate(DB *gorm.DB) error {
	if err := DB.AutoMigrate(all()...); err != nil {
		return err
	}
	return nil
}

func all() []interface{} {
	out := []interface{}{}
	for _, v := range allMap() {
		out = append(out, v)
	}
	return out
}

func allMap() map[string]interface{} {
	return map[string]interface{}{
		"User":      models.User{},
		"Artist":    models.Artist{},
		"Album":     models.Album{},
		"UserAlbum": models.UserAlbum{},
	}
}
