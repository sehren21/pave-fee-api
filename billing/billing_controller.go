package billing

import (
	"context"

	"encore.app/billing/dto"
	"encore.app/billing/helpers"
	"encore.app/billing/model"
)

//encore:api public method=POST path=/billing/open
func (s *Service) OpenBilling(ctx context.Context, body dto.OpenBillingRequest) (*dto.BaseResponse[*model.Bill], error) {
	bill, err := s.Open(ctx, body)

	if err != nil {
		return nil, err
	}

	return &dto.BaseResponse[*model.Bill]{
		Status:  true,
		Message: "success",
		Data:    bill,
	}, nil
}

//encore:api public method=GET path=/billing
func (s *Service) ListBilling(ctx context.Context, param dto.ListBillingRequest) (*helpers.PaginationResponse[model.Bill], error) {
	return s.List(ctx, param)
}

//encore:api public method=GET path=/billing/:id
func (s *Service) DetailBilling(ctx context.Context, id string) (*dto.BaseResponse[*model.Bill], error) {
	bill, err := s.Detail(ctx, id)

	if err != nil {
		return nil, err
	}

	return &dto.BaseResponse[*model.Bill]{
		Status:  true,
		Message: "success",
		Data:    bill,
	}, nil
}

//encore:api public method=POST path=/billing/add-item
func (s *Service) AddBillingItem(ctx context.Context, body dto.AddBillItemRequest) (*dto.BaseResponse[string], error) {
	err := s.client.SignalWorkflow(
		ctx,
		"bill-"+body.BillID,
		"",
		"add-bill-item",
		body,
	)

	if err != nil {
		return nil, err
	}

	return &dto.BaseResponse[string]{
		Message: "success",
		Status:  true,
	}, nil
}

//encore:api public method=POST path=/billing/close
func (s *Service) CloseBilling(ctx context.Context, body dto.CloseBillingRequest) (*dto.BaseResponse[string], error) {
	err := s.client.SignalWorkflow(
		ctx,
		"bill-"+body.BillID,
		"",
		"close-bill",
		body,
	)

	if err != nil {
		return nil, err
	}

	return &dto.BaseResponse[string]{
		Message: "success",
		Status:  true,
	}, nil
}
