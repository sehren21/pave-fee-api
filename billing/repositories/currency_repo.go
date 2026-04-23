package repositories

import (
	"encore.app/billing/model"
	"gorm.io/gorm"
)

type ICurrencyRepo interface {
	FindOne(code string) (*model.Currency, error)
}
type CurrencyRepo struct {
	db *gorm.DB
}

func NewCurrencyRepo(db *gorm.DB) ICurrencyRepo {
	return &CurrencyRepo{
		db: db,
	}
}

func (r *CurrencyRepo) FindOne(code string) (*model.Currency, error) {
	var currency *model.Currency

	if err := r.db.Where("code = ?", code).First(&currency).Error; err != nil {
		return nil, err
	}

	return currency, nil
}
