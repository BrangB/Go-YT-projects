package config

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("user.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database connected successfully")
		fmt.Println(DB)
	}

}
