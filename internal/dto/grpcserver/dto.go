package grpcserver

type ProductPageSizeDTO struct {
	Offset int64 `validate:"required,gte=0,lte=100"`
	Limit  int64 `validate:"required,gt=0,lte=100"`
}
