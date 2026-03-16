package server

import (
	"github.com/glebateee/auto-inventory/internal/domain/models"
	aiv1 "github.com/glebateee/auto-proto/gen/go/inventory"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToGRPCProductList(products []models.Product) []*aiv1.Product {
	outProducts := make([]*aiv1.Product, 0, len(products))
	for _, p := range products {
		outProducts = append(outProducts, &aiv1.Product{
			Id:           p.Id,
			Sku:          p.Sku,
			Name:         p.Name,
			Description:  p.Description,
			Category:     p.Category,
			Manufacturer: p.Manufacturer,
			Weight:       p.Weight,
			Price:        p.Price,
			BasePrice:    p.BasePrice,
			IssueYear:    int64(p.IssueYear),
			CreatedAt:    timestamppb.New(p.CreatedAt),
			UpdatedAt:    timestamppb.New(p.UpdatedAt),
		})
	}
	return outProducts
}
