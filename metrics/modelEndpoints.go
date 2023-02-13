package metrics

import (
	"database/sql"
	"log"

	"gorm.io/gorm"
)

var db *gorm.DB

type Endpoint struct {
	gorm.Model
	Url    string       `gorm:"unique"`
	Active sql.NullBool `gorm:"default:true"`
}

// psql:listEndpoints()
func List() (results []Endpoint, err error) {
	err = db.Table("endpoints").Where("active = ?", true).Find(&results).Error
	return results, err
}

// psql:addEndpoint()
func (e *Endpoint) Add() error {
	err := db.Create(e).Error
	if err != nil {
		log.Println("ERR: Unable to add endpoint", err)
	}
	return nil
}

// psql:delEndpoint()
func (e *Endpoint) Del() error {
	return nil
}
