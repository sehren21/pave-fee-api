package activities

import (
	"context"
	"time"

	"encore.app/billing/dto"
	"encore.app/billing/model"
	"encore.app/billing/repositories"
	"encore.dev/rlog"
	"encore.dev/types/uuid"
	"github.com/shopspring/decimal"
)

type BillActivity struct {
	BillRepo     repositories.IBillRepo
	BillItemRepo repositories.IBillItemRepo
}

func NewBillActivity(
	billRepo repositories.IBillRepo,
	billItemRepo repositories.IBillItemRepo,
) *BillActivity {
	return &BillActivity{
		BillRepo:     billRepo,
		BillItemRepo: billItemRepo,
	}
}

func (a *BillActivity) AddBillItemActivity(_ context.Context, item dto.AddBillItemRequest) (*model.BillItem, error) {
	rlog.Info("Adding bill item", "billID", item.BillID)
	billUUID, err := uuid.FromString(item.BillID)
	if err != nil {
		return nil, err
	}
	_, err = a.BillRepo.FindOne(billUUID.String())
	if err != nil {
		return nil, err
	}

	amount, err := decimal.NewFromString(item.Amount)
	if err != nil {
		return nil, err
	}

	billItem := &model.BillItem{
		BillID:      billUUID,
		Amount:      amount,
		Description: item.Description,
	}

	return a.BillItemRepo.Create(billItem)
}

func (a *BillActivity) CloseBillActivity(_ context.Context, billID string, total decimal.Decimal) (*model.Bill, error) {
	bill, err := a.BillRepo.FindOne(billID)
	if err != nil {
		return nil, err
	}
	bill.Amount = total
	bill.Status = model.BillStatusClosed
	now := time.Now()
	bill.ClosedAt = &now
	return a.BillRepo.Update(bill)
}
