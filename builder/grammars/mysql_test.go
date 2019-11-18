package grammars

import (
	"github.com/LibiChai/funsql/builder"
	"testing"
)

func TestNew(t *testing.T) {
	m := newMysqlGrammar()
	t.Logf("%s grammar's placeholder is %s",m.grammarName,m.placeholder)
}

func TestBaseGrammar_CompileSelect(t *testing.T) {
	b := builder.New("table_name")
	b.SetGrammar("mysql")
	b.Where("c1","=",0)
	b.OrWhere("c1","=",0)
	b.WhereRaw("c1 = a1")
	b.WhereRaw("1 != 1")
	b.WhereRaw("c1 = a1 + ?",1)
	b.WhereNotIn("a",[]int{1,2,3})
	b.WhereIn("a",[]string{"c","b",})
	b.OrWhereBetween("age",8,18)
	b.WhereNotBetween("sex","boy","girl")
	b.OrWhereRaw("name in (select name from users) as user_table")
	s1,v1,err1 :=b.Select()
	t.Logf("mysql grammar compile result sql: %s ,binds %+v, err %+v",s1,v1,err1)

	s2,v2,err2 := b.Select("name","age","level")

	t.Logf("mysql grammar compile result sql: %s ,binds %+v, err %+v",s2,v2,err2)

}
