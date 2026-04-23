package repositories

import (
	"encore.app/billing/model"
	"gorm.io/gorm"
)

type IBillItemRepo interface {
	Create(billItem *model.BillItem) (*model.BillItem, error)
}

type BillItemRepo struct {
	db *gorm.DB
}

func NewBillItemRepo(db *gorm.DB) IBillItemRepo {
	return &BillItemRepo{
		db: db,
	}
}

func (r *BillItemRepo) Create(billItem *model.BillItem) (*model.BillItem, error) {
	if err := r.db.Create(billItem).Error; err != nil {
		return nil, err
	}
	return billItem, nil
}
