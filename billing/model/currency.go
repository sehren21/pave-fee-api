package model

import "encore.app/billing/helpers"

type Currency struct {
	helpers.BaseModel
	Code string `json:"code" gorm:"type:text;not null"`
	Name string `json:"name" gorm:"type:text;not null"`
}
