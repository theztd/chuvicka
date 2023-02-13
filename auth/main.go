package auth

import (
	"fmt"
	"log"
)

var DBHost, DBPort, DBUser, DBPassword, DBName, JWTHash string

func Auth(login, password string) (user User, err error) {
	err = db.Table("users").Where("login = ?", login).First(&user).Error
	if err != nil {
		log.Println("ERR: Invalid login", err)
	}

	if user.ValidatePassword(password) {
		return user, nil
	} else {
		/*
			Unauthorized user
		*/
		return User{}, fmt.Errorf("Unauthorized, wrong username or password")
	}
}

func Get() {

}
