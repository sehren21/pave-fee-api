package dto

type CloseBillingRequest struct {
	BillID string `json:"bill_id" validate:"required"`
}
