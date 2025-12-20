package models

import (
	"github.com/mohod24/go-project-management/models/types"
) 

type ListPosition struct {
	InteralID int64 `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID  int64 `json:"public_id" db:"public_id" gorm:"column:public_id"`
	BoardID   int64 `json:"board_internal_id" db:"board_internal_id" gorm:"column:board_internal_id"`

	//ListOrder
	ListOrder types.UUIDArray `json:"list_order"`
}