package models

type UpdateProductFields struct {
	Name         *string `update:"name"`
	Description  *string `update:"description"`
	Category     *string `update:"category"`
	Manufacturer *string `update:"manufacturer"`
	Weight       *int64  `update:"weight"`
	Price        *int64  `update:"price"`
	BasePrice    *int64  `update:"base_price"`
	IssueYear    *int16  `update:"issue_year"`
}
