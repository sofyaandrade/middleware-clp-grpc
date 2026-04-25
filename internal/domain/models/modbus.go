package models

import "gorm.io/gorm"

type CLP struct {
	gorm.Model
	Description string  `json:"description"`
	Ip          string  `json:"ip"`
	TypeClpId   uint    `json:"type_clp_id"`
	TypeClp     TypeClp `json:"type_clp" gorm:"foreignKey:TypeClpId;references:ID"`
	Port        int     `json:"port"`
	IdPlc       int     `json:"id_plc"`
	Tags        []Tag   `json:"tags" gorm:"foreignKey:IdClp;references:ID"`
}
