package repositories

import (
	"encore.app/billing/dto"
	"encore.app/billing/helpers"
	"encore.app/billing/model"
	"gorm.io/gorm"
)

type IBillRepo interface {
	Create(bill *model.Bill) (*model.Bill, error)
	FindOne(id string) (*model.Bill, error)
	Update(bill *model.Bill) (*model.Bill, error)
	List(req dto.ListBillingRequest) (helpers.PaginationResponse[model.Bill], error)
}
type BillRepo struct {
	db *gorm.DB
}

func NewBillRepo(db *gorm.DB) IBillRepo {
	return &BillRepo{
		db: db,
	}
}

func (r *BillRepo) Create(bill *model.Bill) (*model.Bill, error) {
	if err := r.db.Create(bill).Error; err != nil {
		return nil, err
	}
	return bill, nil
}

func (r *BillRepo) FindOne(id string) (*model.Bill, error) {
	var bill *model.Bill

	if err := r.db.Model(&model.Bill{}).
		Select("bills.*, (COALESCE(SUM(bill_items.amount), 0) + bills.amount) as amount").
		Joins("LEFT JOIN bill_items ON bills.id = bill_items.bill_id AND bill_items.deleted_at IS NULL").
		Where("bills.id = ?", id).
		Group("bills.id").
		First(&bill).Error; err != nil {
		return nil, err
	}

	return bill, nil
}

func (r *BillRepo) Update(bill *model.Bill) (*model.Bill, error) {
	if err := r.db.Save(bill).Error; err != nil {
		return nil, err
	}
	return bill, nil
}

func (r *BillRepo) List(req dto.ListBillingRequest) (helpers.PaginationResponse[model.Bill], error) {
	var bills []model.Bill
	var total int64

	query := r.db.Model(&model.Bill{}).
		Select("bills.*, (COALESCE(SUM(bill_items.amount), 0) + bills.amount) as amount").
		Joins("LEFT JOIN bill_items ON bills.id = bill_items.bill_id AND bill_items.deleted_at IS NULL").
		Group("bills.id")

	if req.CustomerID != "" {
		query = query.Where("bills.customer_id = ?", req.CustomerID)
	}

	// Count first - need to count distinct bills
	countQuery := r.db.Model(&model.Bill{})
	if req.CustomerID != "" {
		countQuery = countQuery.Where("customer_id = ?", req.CustomerID)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return helpers.PaginationResponse[model.Bill]{}, err
	}

	// Fetch paginated data
	if err := query.
		Limit(req.GetLimit()).
		Offset(req.GetOffset()).
		Find(&bills).Error; err != nil {
		return helpers.PaginationResponse[model.Bill]{}, err
	}

	return helpers.NewPaginationResponse(bills, total, req.Page, req.Limit), nil
}
