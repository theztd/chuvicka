package model

import (
	"database/sql"
	"log"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login          string       `json:"login" gorm:"unique"`
	Email          string       `json:"email" gorm:"unique"`
	Password       string       `json:"password"`
	Active         sql.NullBool `gorm:"default:true"`
	Token          string       `gorm:"-"` // ignore field from DB
	Role           string       `gorm:"default:reader"`
	OrganizationID uint
}

// Used for reading
type UserRead struct {
	ID     uint
	Email  string
	Login  string
	Active sql.NullBool
}

func (usr *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	usr.Password = string(bytes)
	return nil
}

func (usr *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (usr *User) JWTGenerate() error {
	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chuvicka",
		ID:     usr.Email,
	})
	tokenStr, _ := _token.SignedString([]byte(JWTHash))
	usr.Token = tokenStr
	return nil
}

func (usr *User) Register() error {
	q := DB.Create(&usr)
	if q.Error != nil {
		log.Println("ERR: Unable to register user", q.Error)
		return q.Error
	}
	return nil
}
