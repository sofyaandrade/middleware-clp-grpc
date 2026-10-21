package models

import "gorm.io/gorm"

type UserProfile struct {
	gorm.Model
	Description string `json:"description"`
}
