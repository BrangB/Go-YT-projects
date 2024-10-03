package config

import (
	"fmt"

	"github.com/brangb/go_jwt_sqlite/models"
)

func SyncDatabase() {
	err := DB.AutoMigrate(&models.User{}, &models.Room{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Database migration is successfully done!!!")
	}
}
