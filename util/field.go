package util

import "reflect"

func GetFieldName(field reflect.StructField) string {
	fieldName := field.Tag.Get("fs")
	if fieldName == "" {
		fieldName = Underscore(field.Name)
	}
	return fieldName
}
