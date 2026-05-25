package models

import "gorm.io/gorm"

type Swap struct {
	gorm.Model
	Description string `json:"description"`
	OrderSwap   string `json:"order_swap"`
}
