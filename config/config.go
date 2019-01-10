package config

import (
	"github.com/jinzhu/gorm"
	"os"
)

var DB *gorm.DB

func Open() error {
	var err error
	DB, err = gorm.Open(
		"postgres",
		"host="+os.Getenv("DB_HOST")+" user="+os.Getenv("DB_USERNAME")+
			" dbname="+os.Getenv("DB_DATABASE")+" sslmode=disable password="+
			os.Getenv("DB_PASSWORD"))
	return err
}

func Close() error {
	return DB.Close()
}
