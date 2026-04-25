package models

import "gorm.io/gorm"

type TypeTag struct {
	gorm.Model
	Description string
}
