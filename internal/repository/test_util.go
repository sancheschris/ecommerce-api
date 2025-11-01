package repository

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB(models ...interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(models...)
	if err != nil {
		panic(err)
	}
	return db
}
