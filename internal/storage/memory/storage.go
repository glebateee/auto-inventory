package memory

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/glebateee/auto-inventory/internal/domain/models"
)

type Storage struct {
	products []models.Product
}

func New() *Storage {
	st := &Storage{
		products: make([]models.Product, 10),
	}
	st.Seed()
	return st
}

func (s *Storage) Seed() {
	for i := range s.products {
		s.products[i] = models.Product{
			Id:           int64(i),
			Sku:          gofakeit.Unit(),
			Name:         gofakeit.CarTransmissionType(),
			Description:  gofakeit.ProductDescription(),
			Category:     gofakeit.CarModel(),
			Manufacturer: gofakeit.MinecraftAnimal(),
			Weight:       int64(gofakeit.UintRange(10, 20)),
			Price:        int64(gofakeit.UintRange(100, 200)),
			BasePrice:    int64(gofakeit.UintRange(100, 200)),
			IssueYear:    int16(gofakeit.Year()),
		}
	}
}

func (s *Storage) ProductPageSize(ctx context.Context, page int64, size int64) ([]models.Product, int64) {
	start := (page - 1) * size
	if page > 0 && int64(len(s.products)) > start {
		end := min(start+size, int64(len(s.products)))
		return s.products[start:end], int64(len(s.products))
	}
	return nil, int64(len(s.products))
}
