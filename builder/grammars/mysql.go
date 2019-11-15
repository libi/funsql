package grammars

import "github.com/LibiChai/funsql/builder"

func init(){
	builder.RegisterGrammar("mysql",newMysqlGrammar())
}
type MysqlGrammar struct {
	baseGrammar
}

func newMysqlGrammar() *MysqlGrammar{
	mysqlGrammar := new(MysqlGrammar)
	mysqlGrammar.placeholder = "?"
	return mysqlGrammar
}
