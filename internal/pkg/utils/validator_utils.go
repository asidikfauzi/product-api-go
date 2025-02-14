package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func FormatValidationError(input interface{}) map[string][]string {
	validate := validator.New()
	out := make(map[string][]string)

	if err := validate.Struct(input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			t := reflect.TypeOf(input)

			for _, fe := range validationErrors {
				fieldName := getJSONFieldName(t, fe.StructField())
				if fieldName != "" {
					out[fieldName] = append(out[fieldName], validationMessage(fe))
				}
			}
		}
	}

	return out
}

func getJSONFieldName(t reflect.Type, structField string) string {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == structField {
			return field.Tag.Get("json")
		}
	}
	return ""
}

func validationMessage(fe validator.FieldError) string {
	fieldName := FormatFieldName(fe.Field())

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s cannot be empty", fieldName)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fieldName)
	case "min":
		return fmt.Sprintf("%s must have at least %s characters", fieldName, fe.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", fieldName, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of the following values: %s", fieldName, fe.Param())
	case "eqfield":
		return fmt.Sprintf("%s must be the same as %s", fieldName, FormatFieldName(fe.Param()))
	default:
		return fe.Error()
	}
}
