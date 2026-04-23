package dto

type AddBillItemRequest struct {
	BillID      string `json:"bill_id" validate:"required,uuid4"`
	Amount      string `json:"amount" validate:"required,gt=0"`
	Description string `json:"description" validate:"required,min=1,max=255"`
}
