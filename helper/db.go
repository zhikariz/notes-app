package helper

import (
	"notes-app/entity"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDb() *gorm.DB {
	dsn := os.Getenv("DSN_CONNECTION")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error Connecting To Database")
	}
	db.AutoMigrate(entity.Account{}, entity.Note{})
	return db
}
