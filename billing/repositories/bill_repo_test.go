package repositories

import (
	"testing"
	"time"

	"encore.app/billing/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBillRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	repo := &BillRepo{db: gormDB}

	bill := &model.Bill{
		CustomerID: "test-customer",
		CurrencyId: uuid.New(),
		Status:     model.BillStatusOpen,
		Amount:     decimal.NewFromFloat(100.50),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "bills"`).WillReturnRows(sqlmock.NewRows([]string{"id", "currency_id", "period_start", "period_end", "closed_at"}).AddRow(uuid.New(), uuid.New(), time.Now(), nil, nil))
	mock.ExpectCommit()

	result, err := repo.Create(bill)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
