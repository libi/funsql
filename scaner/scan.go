package db

import (
	"database/sql"
	"github.com/LibiChai/funsql/util"
	"github.com/pkg/errors"
	"reflect"
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

	if structType.Kind() == reflect.Struct {
		fieldIndexs, err := mapColumnFields(rows, structType)
		if err != nil {
			return err
		}
		v1 := make([]reflect.Value, 0)
		for rows.Next() {
			item := reflect.New(structType)
			if err := scanStruct(rows, fieldIndexs, item); err != nil {
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
	return nil
}

func mapColumnFields(rows *sql.Rows, typ reflect.Type) (fieldIndexs []int, err error) {
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	for i := 0; i < typ.NumField(); i++ {

		fieldName := util.GetFieldName(typ.Field(i))
		if fieldName == "-" {
			continue
		}

		for j, column := range columns {
			if fieldName == column {
				fieldIndexs = append(fieldIndexs, j)
				break
			}
		}
	}
	if len(fieldIndexs) != len(columns) {
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

func ScanRow(rows *sql.Rows, value interface{}) error {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return errors.New("value is nil")
	}
	if v.Kind() != reflect.Ptr {
		return errors.Errorf("value non-pointer %T", value)
	}
	if v.Elem().Kind() != reflect.Struct {
		return errors.New("value not is struct")
	}

	fieldIndexs, err := mapColumnFields(rows, v.Elem().Type())
	if err != nil {
		return err
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}
	return scanStruct(rows, fieldIndexs, v)
}

func scanStruct(rows *sql.Rows, fieldIndexs []int, value reflect.Value) error {
	params := make([]interface{}, len(fieldIndexs))
	for i, index := range fieldIndexs {
		params[index] = value.Elem().Field(i).Addr().Interface()
	}
	return rows.Scan(params...)
}
