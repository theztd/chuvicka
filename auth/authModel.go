package auth

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Login    string       `json:"login" gorm:"unique"`
	Email    string       `json:"email" gorm:"unique"`
	Password string       `json:"password"`
	Active   sql.NullBool `gorm:"default:true"`
	Token    string       `gorm:"-"` // ignore field from DB
	ApiToken string       `gorm:"api_token"`
}

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

func (usr *User) Register() error {
	q := db.Create(&usr)
	if q.Error != nil {
		log.Println("ERR: Unable to register user", q.Error)
		return q.Error
	}
	return nil
}

func Connect() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Prague", DBHost, DBUser, DBPassword, DBName, DBPort)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("FATAL: Unable to connect database (%s@%s:%s/%s)\n", DBUser, DBHost, DBPort, DBName)
	}
	return err
}

func Migrate() error {
	return db.AutoMigrate(&User{})
}
