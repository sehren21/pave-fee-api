package helpers

type PaginationRequest struct {
	Limit  int    `json:"limit" validate:"required"`
	Page   int    `json:"page" validate:"required"`
	Search string `json:"search"`
}

func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *PaginationRequest) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *PaginationRequest) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}
func (p *PaginationRequest) GetSearch() string {
	return p.Search
}

type PaginationResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

func NewPaginationResponse[T any](data []T, total int64, page, pageSize int) PaginationResponse[T] {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return PaginationResponse[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
