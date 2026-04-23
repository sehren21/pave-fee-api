package dto

type OpenBillingRequest struct {
	CustomerID   string `json:"customer_id" validate:"required"`
	CurrencyCode string `json:"currency_code" validate:"required,len=3"`
	Amount       string `json:"amount" validate:"required,gt=0"`
}
