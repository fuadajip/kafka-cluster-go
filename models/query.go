package models

// CustomGormPaginationQuery will return objec from given data to build custom gorm
// query builder pagination
type CustomGormPaginationQuery struct {
	Page    int
	Limit   int
	OrderBy string
	Order   string
}
