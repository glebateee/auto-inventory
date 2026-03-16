package provider

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/glebateee/auto-inventory/internal/domain/models"
	"github.com/glebateee/auto-inventory/internal/lib/sl"
	"github.com/glebateee/auto-inventory/internal/storage"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var (
	ErrInvalidParams = errors.New("invalid parameters for request")
)

type ProductProvider interface {
	ProductPageSize(ctx context.Context, page int64, size int64) ([]models.Product, int64, error)
	ProductPageSizeCategory(ctx context.Context, offset int64, limit int64, categoryID int64) ([]models.Product, int64, error)
	Products(ctx context.Context) ([]models.Product, error)
	UpdateProduct(ctx context.Context, sku string, fields *models.UpdateProductFields, mask *fieldmaskpb.FieldMask) (*models.Product, error)
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

func (ps *ProviderService) UpdateProduct(
	ctx context.Context,
	sku string,
	fields *models.UpdateProductFields,
	mask *fieldmaskpb.FieldMask,
) (*models.Product, error) {
	const op = "services.provider.UpdateProduct"
	logger := ps.logger.With(
		slog.String("op", op),
		slog.String("sku", sku),
	)
	logger.Info("updating product")
	product, err := ps.provider.UpdateProduct(ctx, sku, fields, mask)
	if err != nil {
		logger.Error("failed to update product", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("product updated successfully", slog.Int64("id", product.Id))
	return product, nil
}

func (ps *ProviderService) Products(ctx context.Context) ([]models.Product, error) {
	const op = "services.provider.Products"
	logger := ps.logger.With(
		slog.String("op", op),
	)
	products, err := ps.provider.Products(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			logger.Warn("no rows found in products table")
			return []models.Product{}, fmt.Errorf("%s: %w", op, ErrInvalidParams)
		}
		return []models.Product{}, fmt.Errorf("%s: %w", op, err)
	}
	return products, nil
}

func (ps *ProviderService) ProductPageSizeCategory(ctx context.Context, offset int64, limit int64, categoryID int64) ([]models.Product, int64, error) {
	const op = "services.provider.ProductPageSize"
	logger := ps.logger.With(
		slog.String("op", op),
		slog.Int64("page", offset),
		slog.Int64("size", limit),
		slog.Int64("category", categoryID),
	)
	logger.Info("processing request")
	products, total, err := ps.provider.ProductPageSizeCategory(ctx, offset, limit, categoryID)
	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			logger.Warn("no rows found with provided parameters", slog.Int64("records", total))
			return nil, total, fmt.Errorf("%s: %w", op, ErrInvalidParams)
		}
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	logger.Info("request processed successfully", slog.Int64("records", total))
	return products, total, nil
}

func (ps *ProviderService) ProductPageSize(ctx context.Context, page int64, size int64) ([]models.Product, int64, error) {
	const op = "services.provider.ProductPageSize"
	logger := ps.logger.With(
		slog.String("op", op),
		slog.Int64("page", page),
		slog.Int64("size", size),
	)
	logger.Info("processing request")
	products, total, err := ps.provider.ProductPageSize(ctx, page, size)
	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			logger.Warn("no rows found with provided parameters", slog.Int64("records", total))
			return nil, total, fmt.Errorf("%s: %w", op, ErrInvalidParams)
		}
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	logger.Info("request processed successfully", slog.Int64("records", total))
	return products, total, nil
}
