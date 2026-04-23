package dto

import "encore.app/billing/helpers"

type ListBillingRequest struct {
	helpers.PaginationRequest
	CustomerID string `json:"customer_id" validate:"required"`
}
