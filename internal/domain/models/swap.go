package models

import "gorm.io/gorm"

type Swap struct {
	gorm.Model
	Description string
}
