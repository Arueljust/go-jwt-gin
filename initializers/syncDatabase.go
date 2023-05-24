package initializers

import "github.com/Arueljust/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
