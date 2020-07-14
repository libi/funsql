package builder

var grammars = make(map[string]Grammar)

type Grammar interface {
	CompileSelect(builder *FunBuilder) (sql string, val []interface{}, err error)
	CompileUpdate(builder *FunBuilder, value map[string]interface{}) (sql string, val []interface{}, err error)
	CompileInsert(builder *FunBuilder, values interface{}) (sql string, val []interface{}, err error)
	CompileDelete(builder *FunBuilder) (sql string, val []interface{}, err error)
}

func RegisterGrammar(grammarName string, grammar Grammar) {
	grammars[grammarName] = grammar
}
