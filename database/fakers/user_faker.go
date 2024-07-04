package fakers

import (
	"AntiqueGo/app/models"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserFaker(db *gorm.DB) *models.User {
	return &models.User{
		ID: 			uuid.New().String(),
		FirstName: 		faker.FirstName(),
		LastName:       faker.LastName(),
        Email:        	faker.Email(),
        Password:      "$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		RememberToken: "",
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
		DeletedAt:     gorm.DeletedAt{},
	}
}