package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	Name        string `gorm:"size:100;not null;index"`
	Description string `gorm:"size:255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

func (r *Role) HasRole(db *gorm.DB, userID string) (bool, error) {
    var count int64
    err := db.Debug().Model(&Role{}).Where("id = ?", userID).Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
