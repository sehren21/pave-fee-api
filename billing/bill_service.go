package billing

import (
	"context"
	"errors"
	"time"

	"encore.app/billing/dto"
	"encore.app/billing/helpers"
	"encore.app/billing/model"
	"encore.app/billing/workflows"
	"encore.dev/rlog"
	"github.com/shopspring/decimal"
	"go.temporal.io/sdk/client"
)

func (s *Service) Open(ctx context.Context, req dto.OpenBillingRequest) (*model.Bill, error) {
	if err := s.validator.Struct(req); err != nil {
		rlog.Error("validation failed", "error", err)
		return nil, errors.New("invalid request")
	}

	currency, err := s.CurrencyRepo.FindOne(req.CurrencyCode)
	if err != nil {
		rlog.Error("currency not found", "code", req.CurrencyCode, "error", err)
		return nil, errors.New("invalid currency code")
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		rlog.Error("invalid amount", "amount", req.Amount, "error", err)
		return nil, errors.New("invalid amount")
	}

	bill := model.Bill{
		CustomerID:  req.CustomerID,
		CurrencyId:  currency.ID,
		Status:      model.BillStatusOpen,
		Amount:      amount,
		PeriodStart: time.Now(),
	}

	res, err := s.BillRepo.Create(&bill)
	if err != nil {
		return nil, err
	}

	// init workflow for open billing
	workflowOptions := client.StartWorkflowOptions{
		ID:        "bill-" + res.ID.String(),
		TaskQueue: "billing-task-queue",
	}

	_, err = s.client.ExecuteWorkflow(
		ctx,
		workflowOptions,
		workflows.BillingWorkflow,
		bill.ID.String(),
	)
	return res, err
}

func (s *Service) Detail(_ context.Context, billID string) (*model.Bill, error) {
	bill, err := s.BillRepo.FindOne(billID)
	if err != nil {
		return nil, err
	}

	return bill, nil
}

func (s *Service) List(_ context.Context, req dto.ListBillingRequest) (*helpers.PaginationResponse[model.Bill], error) {
	bills, err := s.BillRepo.List(req)
	if err != nil {
		return nil, err
	}

	return &bills, nil
}
