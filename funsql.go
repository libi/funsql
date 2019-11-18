package funsql

import "github.com/LibiChai/funsql/builder"
import _ "github.com/LibiChai/funsql/builder/grammars"

const defaultGrammar = "mysql"

func Table(tableName string, grammarName ...string) *builder.FunBuilder {
	b := builder.New(tableName)
	if len(grammarName) > 0 {
		b.SetGrammar(grammarName[0])
	} else {
		b.SetGrammar(defaultGrammar)
	}
	return b
}
