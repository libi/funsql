package grammars

import (
	"errors"
	"fmt"
	. "github.com/LibiChai/funsql/builder"
	"strings"
)

// var selectComponents = []string{ "aggregate", "columns", "from", "joins", "wheres", "groups",
// "havings", "orders", "limit", "offset", "unions", "lock",}

type baseGrammar struct {
	placeholder string

}

func (b *baseGrammar)CompileSelect(builder *FunBuilder)(sql string,val []interface{},err error){

	sqls,err := b.compileSelectComponents(builder)
	if(err != nil){
		return "",nil,err
	}

	return strings.Join(sqls," "),builder.GetBindings()["where"],nil
}
func (b *baseGrammar)CompileUpdate(builder *FunBuilder)(sql string,val []interface{},err error){
	fmt.Println("update")
	return "",nil,nil
}
func (b *baseGrammar)CompileDelete(builder *FunBuilder)(sql string,val []interface{},err error){
	fmt.Println("update")
	return "",nil,nil
}
func (b *baseGrammar)CompileInsert(builder *FunBuilder)(sql string,val []interface{},err error){
	fmt.Println("update")
	return "",nil,nil
}


func (b *baseGrammar)compileSelectComponents(builder *FunBuilder) (sqls []string,err error){
	sqls = make([]string,0)
	columns,err := b.compileColumns(builder)
	if(err != nil){
		return nil,err
	}
	sqls = append(sqls,columns)

	table,err := b.compileTable(builder)
	if(err != nil){
		return nil,err
	}
	sqls = append(sqls,table)

	where,err := b.compileWheres(builder)
	if(err != nil){
		return nil,err
	}
	sqls = append(sqls,where)

	return sqls,nil
}


func (b *baseGrammar)compileTable(builder *FunBuilder) (sql string,err error){
	if(builder.GetTable() == ""){
		return "",NoTableNameErr
	}
	return "from "+builder.GetTable(),nil
}

func (b *baseGrammar)compileColumns(builder *FunBuilder) (sql string,err error){
	selectSql := "select "
	if(builder.GetColumns() == nil){
		return selectSql + "* ",nil
	}
	columns := strings.Join(builder.GetColumns(),", ")
	return selectSql + columns,nil
}

func (b *baseGrammar)compileWheres(builder *FunBuilder) (sql string,err error){
	if(builder.GetWheres() == nil){
		return "",nil
	}
	sql = ""
	for _,where := range builder.GetWheres(){
		if(where.IsAnd){
			sql += "and "
		}else{
			sql += "or "
		}

		var tmpsql string
		switch where.Type {
		case "raw":
			sql += b.whereRaw(where)
		case "basic":
			sql += b.whereBasic(where)
		case "in":
			tmpsql,err = b.whereIn(where,"")
			sql += tmpsql
		case "notIn":
			tmpsql,err = b.whereIn(where,"not")
			sql += tmpsql
		case "null":
			sql += b.whereNull(where,"")
		case "notNull":
			sql += b.whereNull(where,"not")
		case "between":
			tmpsql,err = b.whereBetween(where,"")
			sql += tmpsql
		case "notBetween":
			tmpsql,err = b.whereBetween(where,"not")
			sql += tmpsql
		}
		if(err != nil){
			return "",err
		}
		sql += " "
	}

	return "where " +b.removeLeading(sql),nil
}

func (b *baseGrammar)removeLeading(sql string) string{
	sql = strings.TrimLeft(sql,"and")
	sql = strings.TrimLeft(sql,"or")
	return sql
}

func (b *baseGrammar)whereRaw(where Where) string{
	return where.Sql
}
func (b *baseGrammar)whereBasic(where Where) string{
	return where.Column+" "+where.Operator+ " "+b.parameter(where.Value)
}
func (b *baseGrammar)whereIn(where Where,not string) (sql string,err error){
	values ,ok := where.Value.([]interface{})
	if(!ok){
		//todo err define
		return "",errors.New("wherein value wrong")
	}
	return where.Column +" "+not+" in ("+b.parameterize(values)+")" ,nil
}
func (b *baseGrammar)whereNull(where Where,not string)string{
	return where.Column + " is "+not+" null"
}
func (b *baseGrammar)whereBetween(where Where,not string)(sql string,err error){
	values ,ok := where.Value.([]interface{})
	if(!ok){
		//todo err define
		return "",errors.New("wherein value wrong")
	}
	if(len(values) < 2){
		//todo err define
		return "",errors.New("wherein value wrong")
	}
	min := b.parameter(values[0])
	max := b.parameter(values[1])
	return where.Column+ " "+not+" between "+min+" and "+max,nil
}

func (b *baseGrammar)parameter(value interface{}) string{
	return b.placeholder
}

func (b *baseGrammar)parameterize(value []interface{}) string{
	res := make([]string,len(value))
	for _,v := range value {
		res = append(res,b.parameter(v))
	}
	return strings.Join(res,", ")
}