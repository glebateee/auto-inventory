package postgres

import (
	"github.com/glebateee/auto-inventory/internal/domain/models"
	sqlc "github.com/glebateee/auto-inventory/internal/storage/postgres/sqlc/gen"
)

func mapSlice[T any](rows []T, mapper func(T) models.Product) []models.Product {
	products := make([]models.Product, 0, len(rows))
	for _, r := range rows {
		products = append(products, mapper(r))
	}
	return products
}

func FromSqlcProducts(sqlcProducts []sqlc.ProductsRow) []models.Product {
	return mapSlice(sqlcProducts, func(p sqlc.ProductsRow) models.Product {
		return models.Product{
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
		}
	})
}

func FromSqlcProductList(sqlcProducts []sqlc.ProductPageSizeRow) []models.Product {
	return mapSlice(sqlcProducts, func(p sqlc.ProductPageSizeRow) models.Product {
		return models.Product{
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
		}
	})
}

func FromSqlcProductListCat(sqlcProducts []sqlc.ProductPageSizeCategoryRow) []models.Product {
	return mapSlice(sqlcProducts, func(p sqlc.ProductPageSizeCategoryRow) models.Product {
		return models.Product{
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
		}
	})
}
