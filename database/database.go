package database

import (
	"TestVote/model"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open() *gorm.DB {
	var db *gorm.DB
	var err error

	// Buka Koneksi
	db, err = gorm.Open(
		mysql.Open(
			fmt.Sprintf(
				"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
				os.Getenv("DB_USER"),
				os.Getenv("DB_PASS"),
				os.Getenv("DB_HOST"),
				os.Getenv("DB_NAME"),
			),
		),
		&gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	// Model
	if err = db.AutoMigrate(
		&model.Users{},
		&model.CalonKepala{},
		&model.CalonSenat{}); err != nil {
		log.Fatal(err.Error())
	}

	return db
}
