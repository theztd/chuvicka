package auth

import (
	"fmt"
	"log"
	"theztd/chuvicka/auth/model"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBHost, DBPort, DBUser, DBPassword, DBName string
	// UserRead                                   model.UserRead
	// User					                      model.User
)

func Connect() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Prague", DBHost, DBUser, DBPassword, DBName, DBPort)
	var err error
	model.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("FATAL: Unable to connect database (%s@%s:%s/%s)\n", DBUser, DBHost, DBPort, DBName)
	}
	return err
}

func Migrate() error {
	return model.DB.AutoMigrate(&model.Organization{}, &model.Token{}, &model.User{})
}

func Get() {

}

func Auth(login, password string) (user model.User, err error) {
	err = model.DB.Table("users").Where("login = ?", login).First(&user).Error
	if err != nil {
		log.Println("ERR: Invalid login", err)
	}

	if user.ValidatePassword(password) {
		return user, nil
	} else {
		/*
			Unauthorized user
		*/
		return model.User{}, fmt.Errorf("Unauthorized, wrong username or password")
	}
}

func Users(org model.Organization) (users []model.User, err error) {
	err = model.DB.Model(&model.User{}).Find(&users).Error
	return users, err
}

func GetOrg(name string) (org model.Organization, err error) {
	err = model.DB.Model(&model.Organization{Name: name}).Preload("Users").First(&org).Error
	return org, err
}

func JWTHashSet(str string) {
	model.JWTHash = str
}

func JWTValidate(tokenString string) bool {
	claimData := jwt.RegisteredClaims{}
	ct, err := jwt.ParseWithClaims(
		tokenString,
		&claimData,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(model.JWTHash), nil
		},
	)
	if err != nil {
		log.Println("ERR: Token verificaition error", err)
		return false
	}

	fmt.Println("DEBUG: Parsed token (podpis): ", ct.Signature)
	fmt.Println("DEBUG: Parsed claim: ", claimData)

	return ct.Valid
}

func ValidetaApiToken(token string) bool {
	var tokens []model.Token
	if len(token) < 3 {
		return false
	}
	log.Println("DEBUG: auth-token", token)
	err := model.DB.Model(&model.Token{}).Find(&tokens).Error
	if err != nil {
		log.Println("ERR: token-validate", err)
	}

	// hodne neefektivni
	for _, t := range tokens {
		if t.Validate(token) {
			log.Println("INFO: Authorized token", t.Name)
			return true
		}
	}
	return false

}
