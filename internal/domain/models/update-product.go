package models

type UpdateProductFields struct {
	Name         *string
	Description  *string
	Category     *string
	Manufacturer *string
	Weight       *int64
	Price        *int64
	BasePrice    *int64
	IssueYear    *int16
}
