package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Description   string        `json:"description"`
	ConsumerId    uint          `json:"consumer_id"`
	TypeID        uint          `json:"type_id"`
	Type          TypeTag       `json:"type" gorm:"foreignKey:TypeID;references:ID"`
	SwapID        uint          `json:"swap_id"`
	Swap          Swap          `json:"swap" gorm:"foreignKey:SwapID;references:ID"`
	OperationID   uint          `json:"operation_id"`
	OperationType TypeOperation `json:"operation_type" gorm:"foreignKey:OperationID;references:ID"`
	Offset        int           `json:"offset"`
	IdClp         uint          `json:"id_clp"`
	Clp           CLP           `json:"CLP" gorm:"foreignKey:IdClp;references:ID"`
}
