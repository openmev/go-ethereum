package teller

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDbConnection() (*gorm.DB, error) {
	dsn := "host=localhost user=jonah1005 password=jonah1005 dbname=tellerLog"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}
