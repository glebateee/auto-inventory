package server

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/glebateee/auto-inventory/internal/domain/models"
	"github.com/glebateee/auto-inventory/internal/dto/grpcserver"
	"github.com/glebateee/auto-inventory/internal/lib/sl"
	"github.com/glebateee/auto-inventory/internal/services/provider"
	aiv1 "github.com/glebateee/auto-proto/gen/go/inventory"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternal = "internal server error"
	ErrInvalid  = "invalid params"
)

type Provider interface {
	ProductPageSize(ctx context.Context, page int64, size int64) ([]models.Product, int64, error)
	ProductPageSizeCategory(ctx context.Context, offset int64, limit int64, categoryID int64) ([]models.Product, int64, error)
	Products(ctx context.Context) ([]models.Product, error)
}

type serverApi struct {
	aiv1.UnimplementedInventoryServer
	logger   *slog.Logger
	provider Provider
	validate *validator.Validate
}

func (s *serverApi) ProductList(ctx context.Context, req *aiv1.ProductListRequest) (*aiv1.ProductListResponse, error) {
	products, err := s.provider.Products(ctx)
	if err != nil {
		if errors.Is(err, provider.ErrInvalidParams) {
			return nil, status.Errorf(codes.InvalidArgument, ErrInvalid)
		}
		return nil, status.Errorf(codes.Internal, ErrInternal)
	}
	return &aiv1.ProductListResponse{
		Products: ToGRPCProductList(products),
	}, nil
}

func (s *serverApi) UpdateProduct(ctx context.Context, req *aiv1.UpdateProductRequest) (*aiv1.UpdateProductResponse, error) {
	const op = "serverApi.UpdateProduct"
	logger := s.logger.With(
		slog.String("op", op),
	)
	fields := req.GetFields()
	validateDto := grpcserver.UpdateProductDTO{
		Sku:          req.GetSku(),
		Name:         fields.GetName(),
		Description:  fields.GetDescription(),
		Category:     fields.GetCategory(),
		Manufacturer: fields.GetManufacturer(),
		Weight:       fields.GetWeight(),
		Price:        fields.GetPrice(),
		BasePrice:    fields.GetBasePrice(),
		IssueYear:    fields.GetIssueYear(),
	}
	err := s.validate.Struct(&validateDto)
	if err != nil {
		logger.Error("validation failed", sl.Err(ValidationError(err.(validator.ValidationErrors))))
		return nil, status.Error(codes.InvalidArgument, ErrInvalid)
	}

	domainFields := dtoToUpdateFields(&validateDto, req.UpdateMask)
}

func (s *serverApi) ProductPageSizeCategory(ctx context.Context, req *aiv1.ProductPageSizeCategoryRequest) (*aiv1.ProductPageSizeCategoryResponse, error) {
	const op = "serverApi.ProductPageSizeCategory"
	logger := s.logger.With(
		slog.String("op", op),
	)
	validateDto := grpcserver.ProductPageSizeCategoryDTO{
		Offset:     req.GetPage(),
		Limit:      req.GetSize(),
		CategoryID: req.GetCategoryId(),
	}
	err := s.validate.Struct(&validateDto)
	if err != nil {
		logger.Error("validation failed", sl.Err(ValidationError(err.(validator.ValidationErrors))))
		return nil, status.Error(codes.InvalidArgument, ErrInvalid)
	}
	products, total, err := s.provider.ProductPageSizeCategory(ctx, validateDto.Offset, validateDto.Limit, validateDto.CategoryID)
	if err != nil {
		if errors.Is(err, provider.ErrInvalidParams) {
			return nil, status.Error(codes.InvalidArgument, ErrInvalid)
		}
		return nil, status.Error(codes.Internal, ErrInternal)
	}
	outProducts := ToGRPCProductList(products)
	return &aiv1.ProductPageSizeCategoryResponse{
		Products:  outProducts,
		Available: total,
	}, nil

}

func (s *serverApi) ProductPageSize(ctx context.Context, req *aiv1.ProductPageSizeRequest) (*aiv1.ProductPageSizeResponse, error) {
	const op = "serverApi.ProductPageSize"
	logger := s.logger.With(
		slog.String("op", op),
	)
	validateDto := grpcserver.ProductPageSizeDTO{
		Offset: req.GetPage(),
		Limit:  req.GetSize(),
	}
	err := s.validate.Struct(&validateDto)
	if err != nil {
		logger.Error("validation failed", sl.Err(ValidationError(err.(validator.ValidationErrors))))
		return nil, status.Error(codes.InvalidArgument, ErrInvalid)
	}
	products, total, err := s.provider.ProductPageSize(ctx, validateDto.Offset, validateDto.Limit)
	if err != nil {
		if errors.Is(err, provider.ErrInvalidParams) {
			return nil, status.Error(codes.InvalidArgument, ErrInvalid)
		}
		return nil, status.Error(codes.Internal, ErrInternal)
	}
	outProducts := ToGRPCProductList(products)
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

func Register(srv *grpc.Server, logger *slog.Logger, provider Provider) {
	aiv1.RegisterInventoryServer(srv, &serverApi{
		logger:   logger,
		provider: provider,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	})
}
