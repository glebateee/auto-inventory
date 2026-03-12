package provider

import (
	"context"
	"log/slog"

	"github.com/glebateee/auto-inventory/internal/domain/models"
)

type ProductProvider interface {
	ProductPageSize(ctx context.Context, page int64, size int64) ([]models.Product, int64)
}

type ProviderService struct {
	logger   *slog.Logger
	provider ProductProvider
}

func New(
	logger *slog.Logger,
	provider ProductProvider,
) *ProviderService {
	return &ProviderService{
		logger:   logger,
		provider: provider,
	}
}

func (ps *ProviderService) ProductPageSize(ctx context.Context, page int64, size int64) ([]models.Product, int64) {
	products, total := ps.provider.ProductPageSize(ctx, page, size)
	return products, total
}
