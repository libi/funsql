package builder

import "reflect"

// var operators = []string {
// "=", "<", ">", "<=", ">=", "<>", "!=", "<=>",
// "like", "like binary", "not like", "ilike",
// "&", "|", "^", "<<", ">>",
// "rlike", "regexp", "not regexp",
// "~", "~*", "!~", "!~*", "similar to",
// "not similar to", "not ilike", "~~*", "!~~*",
// }

// FunBuilder
type FunBuilder struct {
	grammar Grammar

	table     string
	bindings  map[string][]interface{}
	columns   []string
	columnRaw string
	joins     []Join
	wheres    []Where
	groups    []string
	havings   []Having
	orders    []Order
	limit     int
	offset    int
}

type Join struct {
	Type  string
	Table string
	Query string
}

type Where struct {
	Type     string // basic raw
	Sql      string
	Column   string
	Operator string
	Value    interface{}
	IsAnd    bool // and OR or ?
}
type Having Where
type Order struct {
	Column string
	Sort   string
}

// New a FunBuilder
func New(talbeName string) *FunBuilder {
	funBuilder := &FunBuilder{
		table: talbeName,
	}
	return funBuilder
}

func (f *FunBuilder) SetGrammar(grammarName string) {
	if g, ok := grammars[grammarName]; ok {
		f.grammar = g
	} else {
		panic("unknown grammar")
	}
}
func (f *FunBuilder) GetTable() string {
	return f.table
}
func (f *FunBuilder) GetColumns() []string {
	return f.columns
}
func (f *FunBuilder) GetColumnsRaw() string {
	return f.columnRaw
}
func (f *FunBuilder) GetJoins() []Join {
	return f.joins
}
func (f *FunBuilder) GetWheres() []Where {
	return f.wheres
}
func (f *FunBuilder) GetBindings() map[string][]interface{} {
	return f.bindings
}
func (f *FunBuilder) GetGroups() []string {
	return f.groups
}
func (f *FunBuilder) GetHavings() []Having {
	return f.havings
}
func (f *FunBuilder) GetOrders() []Order {
	return f.orders
}
func (f *FunBuilder) GetLimit() int {
	return f.limit
}
func (f *FunBuilder) GetOffset() int {
	return f.offset
}

// Select
func (f *FunBuilder) Select(colums ...string) (sql string, val []interface{}, err error) {
	f.columns = colums
	return f.grammar.CompileSelect(f)
}

// SelectRaw
func (f *FunBuilder) SelectRaw(selectRaw string) (sql string, val []interface{}, err error) {
	f.columnRaw = selectRaw
	return f.grammar.CompileSelect(f)
}

// Update
func (f *FunBuilder) Update(value map[string]interface{}) (sql string, val []interface{}, err error) {
	return f.grammar.CompileUpdate(f, value)
}

// Insert
func (f *FunBuilder) Insert(value interface{}) (sql string, val []interface{}, err error) {
	return f.grammar.CompileInsert(f, value)
}

// Delete
func (f *FunBuilder) Delete() (sql string, val []interface{}, err error) {
	return f.grammar.CompileDelete(f)
}

func (f *FunBuilder) Join(tableName string, on string) *FunBuilder {
	f.addJoin(Join{
		Type:  "join",
		Table: tableName,
		Query: on,
	})
	return f
}

func (f *FunBuilder) Where(column string, operator string, value interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:     "basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		IsAnd:    true,
	})
	f.addBindings(value, "where")
	return f
}
func (f *FunBuilder) OrWhere(column string, operator string, value interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:     "basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		IsAnd:    false,
	})
	return f
}
func (f *FunBuilder) WhereIn(column string, values interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:   "in",
		Column: column,
		Value:  values,
		IsAnd:  true,
	})
	f.addBindings(values, "where")
	return f
}
func (f *FunBuilder) WhereNotIn(column string, values interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:   "notIn",
		Column: column,
		Value:  values,
		IsAnd:  true,
	})
	f.addBindings(values, "where")
	return f
}

func (f *FunBuilder) OrWhereIn(column string, values interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:   "in",
		Column: column,
		Value:  values,
		IsAnd:  false,
	})
	f.addBindings(values, "where")
	return f
}

func (f *FunBuilder) OrWhereNotIn(column string, values interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:   "notIn",
		Column: column,
		Value:  values,
		IsAnd:  false,
	})
	f.addBindings(values, "where")
	return f
}
func (f *FunBuilder) OrWhereNotBetween(column string, value ...interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:   "notBetween",
		Column: column,
		Value:  value,
		IsAnd:  false,
	})
	f.addBindings(value, "where")
	return f
}
func (f *FunBuilder) OrWhereBetween(column string, value ...interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:   "between",
		Column: column,
		Value:  value,
		IsAnd:  false,
	})
	f.addBindings(value, "where")
	return f
}
func (f *FunBuilder) WhereNotBetween(column string, value ...interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:   "notBetween",
		Column: column,
		Value:  value,
		IsAnd:  true,
	})
	f.addBindings(value, "where")
	return f
}
func (f *FunBuilder) WhereBetween(column string, value ...interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:   "between",
		Column: column,
		Value:  value,
		IsAnd:  true,
	})
	f.addBindings(value, "where")
	return f
}

func (f *FunBuilder) WhereRaw(sql string, value ...interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:  "raw",
		Sql:   sql,
		Value: value,
		IsAnd: true,
	})
	f.addBindings(value, "where")
	return f
}
func (f *FunBuilder) OrWhereRaw(sql string, value ...interface{}) *FunBuilder {
	f.addWhere(Where{
		Type:  "raw",
		Sql:   sql,
		Value: value,
		IsAnd: false,
	})
	f.addBindings(value, "where")
	return f
}

func (f *FunBuilder) Group(column ...string) *FunBuilder {
	if f.groups == nil {
		f.groups = make([]string, 0)
	}
	f.groups = append(f.groups, column...)
	return f
}
func (f *FunBuilder) Having(column string, operator string, value interface{}) *FunBuilder {
	f.addHaving(Having{
		Type:     "basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		IsAnd:    true,
	})
	f.addBindings(value, "having")
	return f
}
func (f *FunBuilder) OrHaving(column string, operator string, value interface{}) *FunBuilder {
	f.addHaving(Having{
		Type:     "basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		IsAnd:    false,
	})
	f.addBindings(value, "having")
	return f
}
func (f *FunBuilder) HavingNotBetween(column string, value ...interface{}) *FunBuilder {
	f.addHaving(Having{
		Type:   "notBetween",
		Column: column,
		Value:  value,
		IsAnd:  true,
	})
	f.addBindings(value, "having")
	return f
}
func (f *FunBuilder) HavingBetween(column string, value ...interface{}) *FunBuilder {
	f.addHaving(Having{
		Type:   "between",
		Column: column,
		Value:  value,
		IsAnd:  true,
	})
	f.addBindings(value, "having")
	return f
}
func (f *FunBuilder) OrHavingNotBetween(column string, values ...interface{}) *FunBuilder {
	f.addHaving(Having{
		Type:   "notBetween",
		Column: column,
		Value:  values,
		IsAnd:  false,
	})
	f.addBindings(values, "having")
	return f
}
func (f *FunBuilder) OrHavingBetween(column string, values ...interface{}) *FunBuilder {
	f.addHaving(Having{
		Type:   "between",
		Column: column,
		Value:  values,
		IsAnd:  false,
	})
	f.addBindings(values, "having")
	return f
}

func (f *FunBuilder) HavingRaw(sql string, value ...interface{}) *FunBuilder {
	f.addHaving(Having{
		Type:  "raw",
		Sql:   sql,
		Value: value,
		IsAnd: true,
	})
	f.addBindings(value, "having")
	return f
}
func (f *FunBuilder) OrHavingRaw(sql string, value ...interface{}) *FunBuilder {
	f.addHaving(Having{
		Type:  "raw",
		Sql:   sql,
		Value: value,
		IsAnd: false,
	})
	f.addBindings(value, "having")
	return f
}
func (f *FunBuilder) OrderBy(column string) *FunBuilder {
	if f.orders == nil {
		f.orders = make([]Order, 0)
	}
	f.orders = append(f.orders, Order{
		Column: column,
		Sort:   "asc",
	})
	return f
}
func (f *FunBuilder) OrderByDesc(column string) *FunBuilder {
	if f.orders == nil {
		f.orders = make([]Order, 0)
	}
	f.orders = append(f.orders, Order{
		Column: column,
		Sort:   "desc",
	})
	return f
}
func (f *FunBuilder) Limit(limit int) *FunBuilder {
	f.limit = limit
	return f
}
func (f *FunBuilder) Offset(offset int) *FunBuilder {
	f.offset = offset
	return f
}

func (f *FunBuilder) addJoin(j Join) {
	if f.joins == nil {
		f.joins = make([]Join, 0)
	}
	f.joins = append(f.joins, j)
}

func (f *FunBuilder) addWhere(w Where) {
	if f.wheres == nil {
		f.wheres = make([]Where, 0)
	}
	f.wheres = append(f.wheres, w)
}
func (f *FunBuilder) addHaving(w Having) {
	if f.havings == nil {
		f.havings = make([]Having, 0)
	}
	f.havings = append(f.havings, w)
}

func (f *FunBuilder) addBindings(value interface{}, qtype string) {
	if value == nil {
		return
	}
	if f.bindings == nil {
		f.bindings = make(map[string][]interface{})
	}
	if f.bindings[qtype] == nil {
		f.bindings[qtype] = make([]interface{}, 0)
	}
	v := reflect.ValueOf(value)
	if v.Type().Kind() == reflect.Slice {
		length := v.Len()
		for i := 0; i < length; i++ {
			f.bindings[qtype] = append(f.bindings[qtype], v.Index(i).Interface())
		}
	} else {
		f.bindings[qtype] = append(f.bindings[qtype], value)
	}
}
