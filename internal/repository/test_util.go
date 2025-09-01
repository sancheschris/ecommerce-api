package repository

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB(models ...interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(models...)
	return db
}