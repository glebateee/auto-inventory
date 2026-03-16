package grpcserver

type ProductPageSizeDTO struct {
	Offset int64 `validate:"required,gte=0,lte=100"`
	Limit  int64 `validate:"required,gt=0,lte=100"`
}

type ProductPageSizeCategoryDTO struct {
	Offset     int64 `validate:"required,gte=0,lte=100"`
	Limit      int64 `validate:"required,gt=0,lte=100"`
	CategoryID int64 `validate:"required,gt=0,lte=20"`
}

type UpdateProductDTO struct {
	Sku          string `validate:"required,min=3,max=50"`
	Name         string `validate:"omitempty,min=2,max=50"`
	Description  string `validate:"omitempty,max=1000"`
	Category     string `validate:"omitempty,min=1,max=255"`
	Manufacturer string `validate:"omitempty,min=1,max=255"`
	Weight       int64  `validate:"omitempty,min=0"`
	Price        int64  `validate:"omitempty,min=0"`
	BasePrice    int64  `validate:"omitempty,min=0"`
	IssueYear    int64  `validate:"omitempty,min=1900,max=2100"`
}
