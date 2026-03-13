package server

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationError(verrs validator.ValidationErrors) error {
	var errs []string
	for _, err := range verrs {
		switch err.ActualTag() {
		case "required":
			errs = append(errs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "gte":
			errs = append(errs, fmt.Sprintf("field %s should be >= %s, got: %v", err.Field(), err.Param(), err.Value()))
		case "gt":
			errs = append(errs, fmt.Sprintf("field %s should be > %s, got: %v", err.Field(), err.Param(), err.Value()))
		case "lte":
			errs = append(errs, fmt.Sprintf("field %s should be <= %s, got: %v", err.Field(), err.Param(), err.Value()))
		case "lt":
			errs = append(errs, fmt.Sprintf("field %s should be < %s, got: %v", err.Field(), err.Param(), err.Value()))
		default:
			errs = append(errs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return errors.New(strings.Join(errs, ", "))
}
