package repositories

import (
	"os"
	"testing"

	"github.com/sarrooo/go-clean/internal/database"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	testDB = database.TestingSqliteDB()

	database.CreateTestingEntities(testDB)

	os.Exit(m.Run())
}
