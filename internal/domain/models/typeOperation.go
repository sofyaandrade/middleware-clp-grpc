package models

import "gorm.io/gorm"

type TypeOperation struct {
	gorm.Model
	Description string
}
