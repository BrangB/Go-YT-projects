package initializers

import "github.com/brangb/go_jwt/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
