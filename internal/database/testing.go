package database

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestingSqliteDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = migrate(db) // Assuming the migrate function exists in your package
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func GetFailedTx(db *gorm.DB) *gorm.DB {
	failedTx := db.Begin()
	failedTx.Rollback()
	return failedTx
}

// CreateTestingEntities creates dummy entities for testing purposes
func CreateTestingEntities(db *gorm.DB) {
	dummyEntities := []interface{}{
		DummyUsers,
		DummyArtists,
		DummyAlbums,
		DummyUserAlbums,
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := createEntities(tx, dummyEntities...); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Failed to create entities %v", err)
	}
}

// createEntities creates entities if the table is empty
func createEntities(tx *gorm.DB, entities ...interface{}) error {
	for _, entity := range entities {
		// Check if the table is empty
		res := tx.Model(entity).First(&struct{}{})
		if res.Error == nil || res.Error != gorm.ErrRecordNotFound {
			fmt.Println("Table not empty, skipping...")
			continue
		}

		// Create entity if the table is empty
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(entity).Error; err != nil {
			if gorm.ErrDuplicatedKey == err {
				fmt.Printf("%T already exists\n", entity)
				continue
			}
			return err
		}
	}
	return nil
}
