package model

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	Name           string
	Token          string
	Expire         time.Time
	Description    string
	OrganizationID uint
}

func (t *Token) Hash(token string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), 14)
	if err != nil {
		return err
	}
	t.Token = string(bytes)
	return nil
}

func (t *Token) Validate(token string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(t.Token), []byte(token))
	if err != nil {
		return false
	}
	return true
}

func (t *Token) Save() error {
	q := DB.Create(&t)
	if q.Error != nil {
		log.Println("ERR: Unable to save organization", q.Error)
		return q.Error
	}
	return nil
}
