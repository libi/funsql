package grammars

import (
	"fmt"
	. "github.com/LibiChai/funsql/builder"
	"github.com/LibiChai/funsql/util"
	"reflect"
	"strconv"
	"strings"
)

type baseGrammar struct {
	placeholder     string
	grammarName     string
	selectComponent []func(builder *FunBuilder) (string, error)
}

func (b *baseGrammar) init() {
	b.selectComponent = []func(builder *FunBuilder) (string, error){
		b.compileColumns,
		b.compileTable,
		b.complieJoins,
		b.compileWheres,
		b.compileGroups,
		b.compileHavings,
		b.compileOrders,
		b.compileLimit,
		b.compileOffset,
	}
}
func (b *baseGrammar) CompileSelect(builder *FunBuilder) (sql string, val []interface{}, err error) {

	sqls, err := b.compileSelectComponents(builder)
	if err != nil {
		return "", nil, err
	}
	for _, sqlItem := range sqls {
		if sqlItem == "" {
			continue
		}
		sql += sqlItem + " "
	}

	bindings := builder.GetBindings()
	if bindings["where"] != nil {
		val = append(val, bindings["where"]...)
	}
	if bindings["having"] != nil {
		val = append(val, bindings["having"]...)
	}
	return
}
func (b *baseGrammar) CompileUpdate(builder *FunBuilder, value map[string]interface{}) (sql string, val []interface{}, err error) {
	joinSql := ""
	if len(builder.GetJoins()) > 0 {
		joinSql += " "
		joinSql, err = b.complieJoins(builder)
		if err != nil {
			return "", nil, err
		}
	}

	bindings := builder.GetBindings()

	whereSql, err := b.compileWheres(builder)
	updateSql := ""
	isFirst := true
	for column, bindVal := range value {
		if !isFirst {
			updateSql += ","
		}
		updateSql += column + "=" + b.placeholder
		isFirst = false
		val = append(val, bindVal)
	}
	if bindings["where"] != nil {
		val = append(val, bindings["where"]...)
	}
	return fmt.Sprintf("update %s%s set %s %s", builder.GetTable(), joinSql, updateSql, whereSql), val, nil
}
func (b *baseGrammar) CompileInsert(builder *FunBuilder, value interface{}) (sql string, val []interface{}, err error) {
	v1 := reflect.ValueOf(value)
	if v1.Kind() == reflect.Ptr {
		v1 = v1.Elem()
	}
	columns := make([]string, 0)
	placeholders := make([]string, 0)

	switch v1.Kind() {
	case reflect.Struct:
		for i := 0; i < v1.NumField(); i++ {
			fieldName := util.GetFieldName(v1.Type().Field(i))
			if fieldName == "-" {
				continue
			}
			columns = append(columns, fieldName)
			placeholders = append(placeholders, b.placeholder)
			val = append(val, v1.Field(i).Interface())
		}
	case reflect.Map:
		keys := v1.MapKeys()
		for _, k := range keys {
			columns = append(columns, k.String())
			placeholders = append(placeholders, b.placeholder)
			val = append(val, v1.MapIndex(k).Interface())
		}
	}

	return fmt.Sprintf("insert into %s (%s) values (%s)", builder.GetTable(), strings.Join(columns, ","), strings.Join(placeholders, ",")), val, nil
}
func (b *baseGrammar) CompileDelete(builder *FunBuilder) (sql string, val []interface{}, err error) {
	fmt.Println("delete")
	wheres, err := b.compileWheres(builder)
	if err != nil {
		return "", nil, err
	}
	bindings := builder.GetBindings()
	if bindings["where"] != nil {
		val = append(val, bindings["where"]...)
	}
	if bindings["having"] != nil {
		val = append(val, bindings["having"]...)
	}

	return fmt.Sprintf("delete from %s %s", builder.GetTable(), wheres), val, nil
}

func (b *baseGrammar) compileSelectComponents(builder *FunBuilder) (sqls []string, err error) {
	sqls = make([]string, 0)
	for _, component := range b.selectComponent {
		sql, err := component(builder)
		if err != nil {
			return nil, err
		}
		sqls = append(sqls, sql)
	}
	return sqls, nil
}

func (b *baseGrammar) compileTable(builder *FunBuilder) (sql string, err error) {
	if builder.GetTable() == "" {
		return "", NoTableNameErr
	}
	return "from " + builder.GetTable(), nil
}

func (b *baseGrammar) compileColumns(builder *FunBuilder) (sql string, err error) {
	selectSql := "select "
	if builder.GetColumns() == nil {
		return selectSql + "* ", nil
	}
	columns := strings.Join(builder.GetColumns(), ", ")
	return selectSql + columns, nil
}

func (b *baseGrammar) complieJoins(builder *FunBuilder) (sql string, err error) {
	if builder.GetJoins() == nil {
		return "", nil
	}
	sql = ""
	for i, join := range builder.GetJoins() {
		switch join.Type {
		case "join":
			sql += "join "
		case "leftJoin":
			sql += "left join "
		case "rightJoin":
			sql += "right join "
		default:
			sql += "join "

		}

		sql += join.Table
		sql += " on "
		sql += join.Query
		if i != len(builder.GetJoins())-1 {
			sql += " "
		}
	}
	return
}

func (b *baseGrammar) compileWheres(builder *FunBuilder) (sql string, err error) {
	if builder.GetWheres() == nil {
		return "", nil
	}
	sql = ""
	for _, where := range builder.GetWheres() {
		if where.IsAnd {
			sql += "and "
		} else {
			sql += "or "
		}

		var tmpsql string
		switch where.Type {
		case "raw":
			sql += b.whereRaw(where)
		case "basic":
			sql += b.whereBasic(where)
		case "in":
			tmpsql, err = b.whereIn(where, "")
			sql += tmpsql
		case "notIn":
			tmpsql, err = b.whereIn(where, "not")
			sql += tmpsql
		case "null":
			sql += b.whereNull(where, "")
		case "notNull":
			sql += b.whereNull(where, "not")
		case "between":
			tmpsql, err = b.whereBetween(where, "")
			sql += tmpsql
		case "notBetween":
			tmpsql, err = b.whereBetween(where, "not")
			sql += tmpsql
		}
		if err != nil {
			return "", err
		}
		sql += " "
	}

	return "where " + b.removeLeading(sql), nil
}

func (b *baseGrammar) compileGroups(builder *FunBuilder) (sql string, err error) {
	if builder.GetGroups() == nil {
		return "", nil
	}
	return "group by " + b.columnize(builder.GetGroups()), nil
}

func (b *baseGrammar) compileHavings(builder *FunBuilder) (sql string, err error) {
	if builder.GetHavings() == nil {
		return "", nil
	}
	sql = ""
	var tmpsql string
	for _, having := range builder.GetHavings() {
		if having.IsAnd {
			sql += "and "
		} else {
			sql += "or "
		}
		switch having.Type {
		case "basic":
			sql += b.havingBasic(having)

		case "raw":
			sql += having.Sql
		case "between":
			tmpsql, err = b.havingBetween(having, "")
			sql += tmpsql
		case "notBetween":
			tmpsql, err = b.havingBetween(having, "not")
			sql += tmpsql
		}
		if err != nil {
			return "", err
		}
		sql += " "
	}
	return "having " + b.removeLeading(sql), nil
}

func (b *baseGrammar) compileOrders(builder *FunBuilder) (sql string, err error) {
	if builder.GetOrders() == nil {
		return "", nil
	}
	sql += "order by "
	sqlArr := make([]string, 0)
	for _, order := range builder.GetOrders() {
		sqlArr = append(sqlArr, order.Column+" "+order.Sort)
	}
	return sql + strings.Join(sqlArr, ", "), nil
}
func (b *baseGrammar) compileLimit(builder *FunBuilder) (sql string, err error) {
	if builder.GetLimit() == 0 {
		return "", nil
	}
	return "limit " + strconv.Itoa(builder.GetLimit()), nil
}
func (b *baseGrammar) compileOffset(builder *FunBuilder) (sql string, err error) {
	if builder.GetOffset() == 0 {
		return "", nil
	}
	return "offset " + strconv.Itoa(builder.GetOffset()), nil
}

func (b *baseGrammar) removeLeading(sql string) string {
	sql = strings.TrimLeft(sql, "and")
	sql = strings.TrimLeft(sql, "or")
	return sql
}

func (b *baseGrammar) whereRaw(where Where) string {
	return where.Sql
}
func (b *baseGrammar) whereBasic(where Where) string {
	return where.Column + " " + where.Operator + " " + b.parameter(where.Value)
}
func (b *baseGrammar) whereIn(where Where, not string) (sql string, err error) {
	if !b.checkIsSlice(where.Value) {
		return "", WhereParamErr
	}
	return where.Column + " " + not + " in (" + b.parameterize(where.Value) + ")", nil
}
func (b *baseGrammar) whereNull(where Where, not string) string {
	return where.Column + " is " + not + " null"
}
func (b *baseGrammar) whereBetween(where Where, not string) (sql string, err error) {
	values, ok := where.Value.([]interface{})
	if !ok {
		return "", WhereBetweenParamErr
	}
	if len(values) < 2 {
		return "", WhereBetweenParamErr
	}
	min := b.parameter(values[0])
	max := b.parameter(values[1])
	return where.Column + " " + not + " between " + min + " and " + max, nil
}
func (b *baseGrammar) havingBasic(having Having) string {
	return having.Column + " " + having.Operator + " " + b.parameter(having.Value)
}
func (b *baseGrammar) havingBetween(having Having, not string) (sql string, err error) {
	return b.whereBetween(Where(having), not)
}
func (b *baseGrammar) parameter(value interface{}) string {
	return b.placeholder
}

func (b *baseGrammar) parameterize(value interface{}) string {
	v := reflect.ValueOf(value)
	res := make([]string, 0)
	for i := 0; i < v.Len(); i++ {
		res = append(res, b.parameter(v.Index(i).Interface()))
	}
	return strings.Join(res, ", ")
}

func (b *baseGrammar) columnize(columns []string) string {
	return strings.Join(columns, ",")
}

func (b *baseGrammar) checkIsSlice(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	return v.Type().Kind() == reflect.Slice
}
