package metrics

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

var DB *gorm.DB

type Request struct {
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Data    string            `json:"data"`
}

type Response struct {
	Code  uint   `json:"code"`
	Match string `json:"match"`
}

func (r *Response) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {

		return errors.New(fmt.Sprint("[ERROR] Scan error", value))
	}

	return json.Unmarshal(bytes, r)
}

func (r Response) Value() (driver.Value, error) {
	log.Println("[DEBUG] Value")
	return json.Marshal(r)
}

type Endpoint struct {
	gorm.Model
	Url string `gorm:"unique"`
	//   Request  map[string]interface{} `gorm:"type:jsonb;default:'{\"method\":\"get\", \"Headers\":{},\"Data\":\"\"}'"`
	Request  Response     `gorm:"type:jsonb;default:'[]'"`
	Response Response     `gorm:"type:jsonb;default:'[]'"`
	Active   sql.NullBool `gorm:"default:true"`
}

// psql:listEndpoints()
func List() (results []Endpoint, err error) {
	err = DB.Table("endpoints").Where("active = ?", true).Find(&results).Error
	return results, err
}

// psql:addEndpoint()
func (e *Endpoint) Add() error {
	err := DB.Create(e).Error
	if err != nil {
		log.Println("ERR: Unable to add endpoint", err)
	}
	return nil
}

// psql:delEndpoint()
func (e *Endpoint) Del() error {
	DB.Where("url LIKE ?", e.Url).Delete(e)
	return nil
}
