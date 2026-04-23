package billing

import (
	"context"
	"testing"

	"encore.app/billing/dto"
	"encore.app/billing/helpers"
	"encore.app/billing/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockBillRepo struct {
	mock.Mock
}

func (m *MockBillRepo) List(req dto.ListBillingRequest) (helpers.PaginationResponse[model.Bill], error) {
	args := m.Called(req)
	return args.Get(0).(helpers.PaginationResponse[model.Bill]), args.Error(1)
}

func (m *MockBillRepo) Create(bill *model.Bill) (*model.Bill, error) {
	args := m.Called(bill)
	return args.Get(0).(*model.Bill), args.Error(1)
}

func (m *MockBillRepo) FindOne(id string) (*model.Bill, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Bill), args.Error(1)
}

func (m *MockBillRepo) Update(bill *model.Bill) (*model.Bill, error) {
	args := m.Called(bill)
	return args.Get(0).(*model.Bill), args.Error(1)
}

type MockCurrencyRepo struct {
	mock.Mock
}

func (m *MockCurrencyRepo) FindOne(code string) (*model.Currency, error) {
	args := m.Called(code)
	return args.Get(0).(*model.Currency), args.Error(1)
}

func TestService_OpenBilling(t *testing.T) {
	mockBillRepo := new(MockBillRepo)
	mockCurrencyRepo := new(MockCurrencyRepo)

	service := &Service{
		BillRepo:     mockBillRepo,
		CurrencyRepo: mockCurrencyRepo,
	}

	req := dto.OpenBillingRequest{
		CustomerID:   "test-customer",
		CurrencyCode: "USD",
		Amount:       "100.5",
	}

	currency := &model.Currency{
		BaseModel: helpers.BaseModel{ID: uuid.New()},
		Code:      "USD",
		Name:      "US Dollar",
	}

	amount, _ := decimal.NewFromString(req.Amount)

	bill := &model.Bill{
		BaseModel:  helpers.BaseModel{ID: uuid.New()},
		CustomerID: req.CustomerID,
		CurrencyId: currency.ID,
		Status:     model.BillStatusOpen,
		Amount:     amount,
	}

	mockCurrencyRepo.On("FindOne", "USD").Return(currency, nil)
	mockBillRepo.On("Create", mock.AnythingOfType("*model.Bill")).Return(bill, nil)

	result, err := service.OpenBilling(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-customer", result.Data.CustomerID)

	mockCurrencyRepo.AssertExpectations(t)
	mockBillRepo.AssertExpectations(t)
}
