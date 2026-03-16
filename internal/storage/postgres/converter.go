package postgres

import (
	"github.com/glebateee/auto-inventory/internal/domain/models"
	sqlc "github.com/glebateee/auto-inventory/internal/storage/postgres/sqlc/gen"
)

func FromSqlcProductList(sqlcProducts []sqlc.ProductPageSizeRow) []models.Product {
	products := make([]models.Product, 0, len(sqlcProducts))
	for _, p := range sqlcProducts {
		products = append(products, models.Product{
			Id:           int64(p.ID),
			Sku:          p.Sku,
			Name:         p.Name,
			Description:  p.Description.String,
			Category:     p.CategoryName,
			Manufacturer: p.ManufacturerName,
			Weight:       int64(p.Weight),
			Price:        int64(p.Price),
			BasePrice:    int64(p.Baseprice),
			IssueYear:    p.Issueyear,
			CreatedAt:    p.CreatedAt.Time,
			UpdatedAt:    p.UpdatedAt.Time,
		})
	}
	return products
}

func FromSqlcProductListCat(sqlcProducts []sqlc.ProductPageSizeCategoryRow) []models.Product {
	products := make([]models.Product, 0, len(sqlcProducts))
	for _, p := range sqlcProducts {
		products = append(products, models.Product{
			Id:           int64(p.ID),
			Sku:          p.Sku,
			Name:         p.Name,
			Description:  p.Description.String,
			Category:     p.CategoryName,
			Manufacturer: p.ManufacturerName,
			Weight:       int64(p.Weight),
			Price:        int64(p.Price),
			BasePrice:    int64(p.Baseprice),
			IssueYear:    p.Issueyear,
			CreatedAt:    p.CreatedAt.Time,
			UpdatedAt:    p.UpdatedAt.Time,
		})
	}
	return products
}
