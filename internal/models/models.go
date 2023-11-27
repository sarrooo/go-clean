package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type User struct {
	Model
	FirstName string
	LastName  string
	BirthDate time.Time
	Phone     string
	Email     string `gorm:"unique"`
	Password  string

	// Relations
	UserAlbums []*UserAlbum
}

type Artist struct {
	Model
	Name string

	// Relations
	Albums []*Album
}

type Album struct {
	Model
	Name     string `gorm:"uniqueIndex:album_idx"`
	ArtistID uint   `gorm:"uniqueIndex:album_idx"`
	Artist   *Artist
}

type UserAlbum struct {
	Model
	UserID  uint
	User    *User
	AlbumID uint
	Album   *Album
}

type Translation struct {
	Model
	EntityType string `gorm:"uniqueIndex:translation_idx"`
	EntityID   uint   `gorm:"uniqueIndex:translation_idx"`
	Field      string `gorm:"uniqueIndex:translation_idx"`
	Language   string `gorm:"uniqueIndex:translation_idx"`
	Value      string
}
