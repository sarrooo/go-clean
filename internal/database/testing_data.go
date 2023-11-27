package database

import (
	"time"

	"github.com/sarrooo/go-clean/internal/models"
)

var (
	// Models
	m1 = models.Model{
		ID:        1,
		CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	m2 = models.Model{
		ID:        2,
		CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Users
	DummyUsers = []models.User{
		{
			Model:     m1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "jhon.doe@mail.com",
		},
	}

	DummyArtists = []models.Artist{
		{
			Model: m1,
			Name:  "Eminem",
		},
		{
			Model: m2,
			Name:  "Drake",
		},
	}

	DummyAlbums = []models.Album{
		{
			Model:    m1,
			Name:     "Kamikaze",
			ArtistID: DummyArtists[0].ID,
		},
		{
			Model:    m2,
			Name:     "Anti",
			ArtistID: DummyArtists[1].ID,
		},
	}

	DummyUserAlbums = []models.UserAlbum{
		{
			Model:   m1,
			UserID:  DummyUsers[0].ID,
			AlbumID: DummyAlbums[0].ID,
		},
	}
)
