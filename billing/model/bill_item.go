package model

import (
	"encore.app/billing/helpers"
	"encore.dev/types/uuid"
	"github.com/shopspring/decimal"
)

type BillItem struct {
	helpers.BaseModel
	BillID      uuid.UUID       `json:"bill_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Amount      decimal.Decimal `json:"amount" gorm:"type:numeric;default:0;not null"`
	Description string          `json:"description" gorm:"type:text;not null"`
}
