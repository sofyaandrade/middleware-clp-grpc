package models

import "gorm.io/gorm"

type TypeClp struct {
	gorm.Model
	Description string
}
