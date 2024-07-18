package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	RoleID 		  string `gorm:size:36;index"`
	FirstName     string `gorm:"size:100;not null"`
	LastName      string `gorm:"size:100;not null"`
	Email         string `gorm:"size:100;not null;uniqueIndex"`
	Password      string `gorm:"size:255;not null"`
	RememberToken string `gorm:"size:255;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (u *User) FindByEmail(db *gorm.DB,email string) (*User, error) {
	var user User
    err := db.Debug().Model(User{}).Where("LOWER(email) = ?", strings.ToLower(email)).
	    First(&user).Error
	if err!= nil {
        return nil, err
    }

	return &user, nil


}

func (u *User) FindByID(db *gorm.DB,userID string) (*User, error) {
	var user User
    err := db.Debug().Model(User{}).Where("id = ?", userID).
	    First(&user).Error
	if err!= nil {
        return nil, err
    }

	return &user, nil


}

func (u *User) CreateUser(db *gorm.DB, param *User) (*User, error) {
	user := &User{
		ID:        param.ID,
		FirstName: param.FirstName,
		LastName:  param.LastName,
		Email:     param.Email,
		Password:  param.Password,
		RoleID:		 "null",
	}

	err := db.Debug().Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) GetRoleIDByUserID(db *gorm.DB, userID string) (string, error) {
	var user User
	err := db.Debug().Model(User{}).Select("role_id").Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return "", err
	}

	return user.RoleID, nil
}

// func IsAdmin(db *gorm.DB, userID string) (bool, error) {
// 	roleModel := models.Role{}
//     hasRole, err := roleModel.HasRole(s.DB, userID)
// }
