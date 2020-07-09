package scaner

import (
	"database/sql"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func Scan(rows *sql.Rows, value interface{}) error {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return errors.New("value is nil")
	}
	if v.Kind() != reflect.Ptr {
		return errors.Errorf("value non-pointer %T", value)
	}
	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return errors.New("value not is slice")
	}

	typ := v.Type()

	// slice ele type
	structType := indirectType(typ.Elem())

	columns, err := rows.Columns()
	if (err != nil) {
		return err
	}

	if structType.Kind() == reflect.Struct {
		fieldIndexs, err := mapColumnFields(columns, structType)
		if (err != nil) {
			return err
		}
		v1 := make([]reflect.Value, 0)
		for rows.Next() {
			item := reflect.New(structType)
			params := make([]interface{}, len(fieldIndexs))
			for _, index := range fieldIndexs {
				params[index] = item.Elem().Field(index).Addr().Interface()
			}
			err := rows.Scan(params...)
			if err != nil {
				return err
			}
			v1 = append(v1, item.Elem())
		}
		v.Set(reflect.Append(v, v1...))

	} else {
		v1 := make([]reflect.Value, 0)
		for rows.Next() {
			item := reflect.New(v.Type().Elem())
			err := rows.Scan(item.Interface())
			if err != nil {
				return err
			}
			v1 = append(v1, item.Elem())
		}
		v.Set(reflect.Append(v, v1...))
	}
	return errors.New("unsupport data type")
}

func mapColumnFields(columns []string, typ reflect.Type) (fieldIndexs []int, err error) {
	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i).Tag.Get("fs")
		if (fieldName == "") {
			// todo 驼峰转蛇形
			fieldName = strings.ToLower(typ.Field(i).Name)
		}

		for j, column := range columns {
			if (fieldName == column) {
				fieldIndexs = append(fieldIndexs, j)
				break
			}
		}
	}
	if (len(fieldIndexs) != len(columns)) {
		err = errors.New("find some unknown field")
		return
	}
	return
}

func indirectType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func ScanRow(row *sql.Row, res ...interface{}) error {
	return nil
}
