package scaner

import (
	"database/sql"
	"fmt"
	"github.com/libi/funsql/util"
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
	fieldIndexs = make([]int, len(columns))
	if err != nil {
		return
	}
	for n, column := range columns {
		fieldIndex, ok := searchFieldIndex(typ, column)
		if !ok {
			err = errors.New(fmt.Sprintf("not found field: %s", column))
			return
		}
		fieldIndexs[n] = fieldIndex
	}

	return
}
func searchFieldIndex(typ reflect.Type, columnName string) (index int, ok bool) {
	for i := 0; i < typ.NumField(); i++ {

		fieldName := util.GetFieldName(typ.Field(i))
		if fieldName == "-" {
			continue
		}
		if fieldName == columnName {
			return i, true
		}

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
		params[i] = value.Elem().Field(index).Addr().Interface()
	}
	return errors.Wrap(rows.Scan(params...), "scan struce error")
}
