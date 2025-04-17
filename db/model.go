// db/model.go
package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("servers.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	if err := DB.AutoMigrate(&Group{}, &Server{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}
