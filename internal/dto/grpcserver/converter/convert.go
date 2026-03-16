package converter

import (
	"fmt"
	"reflect"

	"github.com/glebateee/auto-inventory/internal/domain/models"
	"github.com/glebateee/auto-inventory/internal/dto/grpcserver"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func UpdateDTOTagsMap(model *models.UpdateProductFields) (map[string]reflect.Value, error) {
	modelType := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()

	fields := make(map[string]reflect.Value, modelType.NumField())
	for i := range modelType.NumField() {
		key, ok := modelType.Field(i).Tag.Lookup("update")
		if !ok {
			return nil, fmt.Errorf("update tag not set for field %s", modelType.Field(i).Name)
		}
		fields[key] = modelValue.Field(i)
	}
	return fields, nil
}

func DtoToUpdateFields(dto *grpcserver.UpdateProductDTO, mask *fieldmaskpb.FieldMask) (*models.UpdateProductFields, error) {
	return nil, nil
	// fields, err := UpdateDTOTagsMap()
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(fields)
	// for _, key := range mask.GetPaths() {
	// 	switch key {
	// 	case "name":
	// 		model.Name = &dto.Name
	// 	case "description":
	// 		model.Description = &dto.Description
	// 	case "category":
	// 		model.Category = &dto.Category
	// 	case "manufacturer":
	// 		model.Manufacturer = &dto.Manufacturer
	// 	case "weight":
	// 		model.Weight = &dto.Weight
	// 	case "price":
	// 		model.Price = &dto.Price
	// 	case "base_price":
	// 		model.BasePrice = &dto.BasePrice
	// 	case "issue_year":
	// 		// IssueYear in DTO is int64, domain expects *int16
	// 		val := int16(dto.IssueYear)
	// 		model.IssueYear = &val
	// 	}
	// }
	// return model, nil
}
