package funsql

import "github.com/LibiChai/funsql/builder"

const defaultGrammar = "mysql"

func Table(tableName string,grammarName ...string) *builder.FunBuilder {
	if(len(grammarName) > 0){
		return builder.New(tableName,grammarName[0])
	}
	return builder.New(tableName,defaultGrammar)
}

