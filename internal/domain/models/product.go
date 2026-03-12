package models

type Product struct {
	Id           int64
	Sku          string
	Name         string
	Description  string
	Category     string
	Manufacturer string
	Weight       int64
	Price        int64
	BasePrice    int64
	IssueYear    int64
}
