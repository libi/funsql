package builder

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

	table string
	bindings map[string][]interface{}
	columns []string
	wheres []Where
	groups []string
	havings []Where
	orders []string
	limit int
}


type Where struct {
	Type string		// basic raw
	Sql string
	Column string
	Operator string
	Value interface{}
	IsAnd bool       // and OR or ?
}

// New a FunBuilder
func New(talbeName string,grammarName string) *FunBuilder{
	funBuilder := &FunBuilder{
		table:talbeName,
	}
	funBuilder.setGrammar(grammarName)
	return funBuilder
}

func (f *FunBuilder)setGrammar(grammarName string){
	if g,ok :=grammars[grammarName];ok{
		f.grammar = g
	}else{
		panic("unknown grammar")
	}
}
func (f *FunBuilder)GetTable() string{
	return f.table
}
func (f *FunBuilder)GetColumns() []string{
	return f.columns
}
func (f *FunBuilder)GetWheres() []Where{
	return f.wheres
}
func (f *FunBuilder)GetBindings() map[string][]interface{}{
	return f.bindings
}

// Select
func (f *FunBuilder)Select(colums ...string)(sql string,val []interface{},err error) {
	f.columns = colums
	return f.grammar.CompileSelect(f)
}

// Update
func (f *FunBuilder)Update(){

}

// Delete

func (f *FunBuilder)Where(column string,operator string,value interface{}) *FunBuilder{
	f.addWhere(Where{
		Type:"basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		IsAnd:    true,
	})
	f.addBindings(value,"where")
	return f
}

func (f *FunBuilder)WhereRaw(sql string,value []interface{}) *FunBuilder{
	f.addWhere(Where{
		Type:"raw",
		Sql:sql,
		Value:value,
	})
	f.addBindings(value,"where")
	return f
}

func (f *FunBuilder)OrWhere(column string,operator string,value interface{}) *FunBuilder{
	f.addWhere(Where{
		Type:"basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		IsAnd:    false,
	})
	return f
}

func (f *FunBuilder)addWhere(w Where){
	if(f.wheres == nil){
		f.wheres = make([]Where,0)
	}
	f.wheres = append(f.wheres,w)
}

func (f *FunBuilder)addBindings(value interface{},qtype string){
	if(value == nil){
		return
	}
	if(f.bindings == nil){
		f.bindings = make(map[string][]interface{})
	}
	if(f.bindings[qtype] == nil){
		f.bindings[qtype] = make([]interface{},0)
	}
	values,ok := value.([]interface{})
	if(ok){
		for _,item := range values{
			f.bindings[qtype] = append(f.bindings[qtype],item)
		}
	}else{
		f.bindings[qtype] = append(f.bindings[qtype],value)
	}
}


