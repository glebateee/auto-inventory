package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/glebateee/auto-inventory/internal/domain/models"
	"github.com/glebateee/auto-inventory/internal/dto/grpcserver"
	"github.com/glebateee/auto-inventory/internal/dto/grpcserver/converter"
	"github.com/glebateee/auto-inventory/internal/lib/sl"
	"github.com/glebateee/auto-inventory/internal/services/provider"
	aiv1 "github.com/glebateee/auto-proto/gen/go/inventory"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ErrInternal = "internal server error"
	ErrInvalid  = "invalid params"
)

type Provider interface {
	ProductPageSize(ctx context.Context, page int64, size int64) ([]models.Product, int64, error)
	ProductPageSizeCategory(ctx context.Context, offset int64, limit int64, categoryID int64) ([]models.Product, int64, error)
	Products(ctx context.Context) ([]models.Product, error)
	//UpdateProduct(ctx context.Context, sku string, fields *models.UpdateProductFields, mask *fieldmaskpb.FieldMask) (*models.Product, error)
	DeleteProductSku(ctx context.Context, sku string) error
}

type serverApi struct {
	aiv1.UnimplementedInventoryServer
	logger   *slog.Logger
	provider Provider
	validate *validator.Validate
}

func (s *serverApi) DeleteProduct(ctx context.Context, req *aiv1.DeleteProductRequest) (*aiv1.DeleteProductResponse, error) {
	const op = "serverApi.DeleteProduct"
	logger := s.logger.With(
		slog.String("op", op),
	)
	dto := grpcserver.DeleteProductDTO{
		Sku: req.GetSku(),
	}
	if err := s.validate.Struct(&dto); err != nil {
		logger.Error("validation error", sl.Err(ValidationError(err.(validator.ValidationErrors))))
		return nil, status.Errorf(codes.InvalidArgument, ErrInvalid)
	}
	if err := s.provider.DeleteProductSku(ctx, dto.Sku); err != nil {
		if errors.Is(err, provider.ErrInvalidParams) {
			return nil, status.Error(codes.InvalidArgument, ErrInvalid)
		}
		return nil, status.Error(codes.Internal, ErrInternal)
	}
	return &aiv1.DeleteProductResponse{}, nil
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

func validateUpdateDTO(model *models.UpdateProductFields, grpcMsg *aiv1.UpdateProductFields) (map[string]reflect.Value, error) {
	grpcMsgType := reflect.TypeOf(grpcMsg).Elem()
	fields, err := converter.UpdateDTOTagsMap(model)
	if err != nil {
		return nil, err
	}
	fmt.Println(fields)
	notFoundFields := make([]string, 0, 3)
	for i := range grpcMsgType.NumField() {
		field := grpcMsgType.Field(i)
		if !field.IsExported() {
			continue
		}
		tag, ok := field.Tag.Lookup("json")
		if !ok {
			return nil, fmt.Errorf("json tag not found for grpc field %s", field.Name)
		}
		jsonName := strings.Split(tag, ",")[0]
		if _, ok := fields[jsonName]; !ok {
			notFoundFields = append(notFoundFields, jsonName)
		}
	}
	if len(notFoundFields) > 0 {
		return nil, fmt.Errorf("update tag not found in dto for fields: %s", strings.Join(notFoundFields, ", "))
	}
	return fields, nil
}

func (s *serverApi) UpdateProduct(ctx context.Context, req *aiv1.UpdateProductRequest) (*aiv1.UpdateProductResponse, error) {
	return nil, status.Error(codes.Internal, "unimplemented method")
	// const op = "serverApi.UpdateProduct"
	// logger := s.logger.With(
	// 	slog.String("op", op),
	// )
	// if req.Fields == nil {
	// 	return nil, status.Error(codes.InvalidArgument, "fields must be provided")
	// }
	// if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
	// 	return nil, status.Error(codes.InvalidArgument, "update_mask must be provided and non-empty")
	// }
	// model := &models.UpdateProductFields{}

	// fields, err := validateUpdateDTO(model, req.GetFields())
	// if err != nil {
	// 	logger.Error("validation dto failed", sl.Err(err))
	// 	return nil, status.Error(codes.Internal, ErrInternal)
	// }
	// fmt.Println(fields)
	// dto := grpcserver.UpdateProductDTO{}
	// converter.DtoToUpdateFields(model, req.GetUpdateMask())
	// return &aiv1.UpdateProductResponse{}, nil
	// validateDto := grpcserver.UpdateProductDTO{
	// 	Sku:          req.GetSku(),
	// 	Name:         fields.GetName(),
	// 	Description:  fields.GetDescription(),
	// 	Category:     fields.GetCategory(),
	// 	Manufacturer: fields.GetManufacturer(),
	// 	Weight:       fields.GetWeight(),
	// 	Price:        fields.GetPrice(),
	// 	BasePrice:    fields.GetBasePrice(),
	// 	IssueYear:    fields.GetIssueYear(),
	// }
	// err := s.validate.Struct(&validateDto)
	// if err != nil {
	// 	logger.Error("validation failed", sl.Err(ValidationError(err.(validator.ValidationErrors))))
	// 	return nil, status.Error(codes.InvalidArgument, ErrInvalid)
	// }
	// domainFields, err := converter.DtoToUpdateFields(&validateDto, req.UpdateMask)
	// if err != nil {
	// 	logger.Error("invalid dto setup", sl.Err(err))
	// 	return nil, status.Error(codes.Internal, ErrInternal)
	// }
	// updatedProduct, err := s.provider.UpdateProduct(ctx, validateDto.Sku, domainFields, req.UpdateMask)

	// if err != nil {
	// 	if errors.Is(err, provider.ErrInvalidParams) {
	// 		return nil, status.Error(codes.InvalidArgument, err.Error())
	// 	}
	// 	logger.Error("failed to update product", sl.Err(err))
	// 	return nil, status.Error(codes.Internal, ErrInternal)
	// }

	// // Convert domain product to gRPC
	// grpcProduct := ToGRPCProduct(updatedProduct) // implement this

	// return &aiv1.UpdateProductResponse{
	// 	Product: grpcProduct,
	// }, nil
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
	if err := s.validate.Struct(&validateDto); err != nil {
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
