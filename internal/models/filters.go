package models

import (
	"strings"

	"github.com/geekilx/restaurantAPI/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

func (f Filters) sortColumn() string {
	for _, safeSort := range f.SortSafeList {
		if f.Sort == safeSort {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	return "id"
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func ValidateFilters(v *validator.Validator, f Filters) {
	println(f.Page > 10_000_000 || f.Page < 1)
	println(f.PageSize > 100 || f.Page < 1)
	v.Check(f.PageSize > 100 || f.PageSize < 1, "page_size", "page size must be between 1 and 100")
	v.Check(f.Page > 10_000_000 || f.Page < 1, "page", "page must be between 1 and 10,000,000")
	v.Check(!validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

func CalculateMetadata(totalRecord, page, pageSize int) Metadata {

	return Metadata{
		PageSize:     pageSize,
		FirstPage:    1,
		CurrentPage:  page,
		LastPage:     (totalRecord + pageSize - 1) / pageSize,
		TotalRecords: totalRecord,
	}

}
