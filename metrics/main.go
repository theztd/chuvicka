package metrics

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBHost, DBPort, DBUser, DBPassword, DBName string

func Migrate() {
	DB.AutoMigrate(&Endpoint{})
}

func Connect() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Prague", DBHost, DBUser, DBPassword, DBName, DBPort)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("FATAL: Unable to connect database (%s@%s:%s/%s)\n", DBUser, DBHost, DBPort, DBName)
	}
	return err
}
