package model

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
	Users       []User
	Tokens      []Token
	Description string
}

func (org *Organization) Save() error {
	q := DB.Create(&org)
	if q.Error != nil {
		log.Println("ERR: Unable to save organization", q.Error)
		return q.Error
	}
	return nil
}

func (org *Organization) Pretty() (err error) {
	fmt.Println("======  Organization  ======")
	fmt.Println("Name:", org.Name)
	fmt.Println("Description: |\n  ", org.Description)
	fmt.Println("Members:")
	for _, usr := range org.Users {
		fmt.Println("  - email:", usr.Email)
		fmt.Println("    id:", usr.ID)
		fmt.Println("    role:", usr.Role)
		fmt.Println("    active:", usr.Active.Bool)
	}
	fmt.Println("Tokens:")
	for _, t := range org.Tokens {
		fmt.Println("  - name:", t.Name)
		fmt.Println("    created:", t.CreatedAt)
		fmt.Println("    expire:", t.Expire)
	}
	fmt.Println("")

	return err
}
