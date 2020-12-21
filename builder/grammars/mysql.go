package grammars

import "github.com/libi/funsql/builder"

func init() {
	builder.RegisterGrammar("mysql", newMysqlGrammar())
}

type MysqlGrammar struct {
	baseGrammar
}

func newMysqlGrammar() *MysqlGrammar {
	mysqlGrammar := new(MysqlGrammar)
	mysqlGrammar.init()
	mysqlGrammar.placeholder = "?"
	mysqlGrammar.grammarName = "mysql"

	return mysqlGrammar
}
