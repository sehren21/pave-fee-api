package model

import (
	"time"

	"encore.app/billing/helpers"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BillStatus string

const (
	BillStatusOpen   BillStatus = "OPEN"
	BillStatusClosed BillStatus = "CLOSED"
)

type Bill struct {
	helpers.BaseModel
	CustomerID  string          `json:"customer_id" gorm:"type:text;not null"`
	CurrencyId  uuid.UUID       `json:"currency_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Status      BillStatus      `json:"status" gorm:"type:bill_status;default:'OPEN'"`
	Amount      decimal.Decimal `json:"amount" gorm:"type:decimal(20,6);default:0.0"`
	PeriodStart time.Time       `json:"period_start" gorm:"type:timestamp;default:now_timestamp();default:null"`
	PeriodEnd   *time.Time      `json:"period_end" gorm:"type:timestamp;default:null"`
	ClosedAt    *time.Time      `json:"closed_at" gorm:"type:timestamp;default:null"`
}
