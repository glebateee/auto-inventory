package server

import (
	"context"
	"errors"
	"strings"

	"github.com/glebateee/auto-inventory/internal/domain/models"
	"github.com/glebateee/auto-inventory/internal/services/provider"
	aiv1 "github.com/glebateee/auto-proto/gen/go/inventory"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Provider interface {
	ProductPageSize(ctx context.Context, page int64, size int64) ([]models.Product, int64, error)
}

type serverApi struct {
	aiv1.UnimplementedInventoryServer
	provider Provider
	validate *validator.Validate
}

func (s *serverApi) ProductPageSize(ctx context.Context, req *aiv1.ProductPageSizeRequest) (*aiv1.ProductPageSizeResponse, error) {
	products, total, err := s.provider.ProductPageSize(ctx, req.GetPage(), req.GetSize())
	if err != nil {
		if errors.Is(err, provider.ErrInvalidParams) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid params")
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
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
		})
	}
	return &aiv1.ProductPageSizeResponse{
		Products:  outProducts,
		Available: total,
	}, nil
}

func (s *serverApi) Health(ctx context.Context, req *aiv1.HealthRequest) (*aiv1.HealthResponse, error) {
	switch strings.ToLower(req.GetStatus()) {
	case "ro":
		return nil, status.Errorf(codes.Internal, "this is internal status")
	case "re":
		return nil, status.Errorf(codes.InvalidArgument, "this is invalid argument status")
	}
	return &aiv1.HealthResponse{
		Status: strings.ToUpper(req.GetStatus()),
	}, nil
}

func Register(srv *grpc.Server, provider Provider) {
	aiv1.RegisterInventoryServer(srv, &serverApi{
		provider: provider,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	})
}
