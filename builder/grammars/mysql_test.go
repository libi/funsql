package grammars

import (
	"github.com/LibiChai/funsql/builder"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	m := newMysqlGrammar()
	t.Logf("%s grammar's placeholder is %s", m.grammarName, m.placeholder)
}

func TestBaseGrammar_CompileSelect(t *testing.T) {
	b := builder.New("table_name")
	b.SetGrammar("mysql")
	b.Where("c1", "=", 0)
	b.OrWhere("c1", "=", 0)
	b.WhereRaw("c1 = a1")
	b.WhereRaw("1 != 1")
	b.WhereRaw("c1 = a1 + ?", 1)
	b.WhereNotIn("a", []int{1, 2, 3})
	b.WhereIn("a", []string{"c", "b"})
	b.OrWhereBetween("age", 8, 18)
	b.WhereNotBetween("sex", "boy", "girl")
	b.OrWhereRaw("name in (select name from users) as user_table")
	s1, v1, err1 := b.Select()
	t.Logf("mysql grammar compile result sql: %s ,binds %+v, err %+v", s1, v1, err1)

	s2, v2, err2 := b.Select("name", "age", "level")

	t.Logf("mysql grammar compile result sql: %s ,binds %+v, err %+v", s2, v2, err2)

}

func TestBaseGrammar_CompileUpdate(t *testing.T) {
	b := builder.New("table_name")
	b.SetGrammar("mysql")
	b.Where("c1", "=", 0)
	b.OrWhere("c1", "=", 0)
	b.WhereRaw("c1 = a1")
	b.WhereRaw("1 != 1")
	b.WhereRaw("c1 = a1 + ?", 1)
	b.WhereNotIn("a", []int{1, 2, 3})
	b.WhereIn("a", []string{"c", "b"})
	b.OrWhereBetween("age", 8, 18)
	b.WhereNotBetween("sex", "boy", "girl")
	b.OrWhereRaw("name in (select name from users) as user_table")
	s1, v1, err1 := b.Update(map[string]interface{}{"c1": 1, "sex": "girl"})
	t.Logf("mysql grammar compile result sql: %s ,binds %+v, err %+v", s1, v1, err1)
}

type TmpStruct struct {
	Name string
	Age  int
}

func TestBaseGrammar_CompileInsert(t *testing.T) {
	b := builder.New("table_name")
	b.SetGrammar("mysql")
	s1, v1, err1 := b.Insert(map[string]interface{}{"c1": 1, "sex": "girl"})
	t.Logf("mysql grammar compile result sql: %s ,binds %+v, err %+v", s1, v1, err1)

	tmp := TmpStruct{
		Name: "t1",
		Age:  20,
	}
	s2, v2, err2 := b.Insert(tmp)
	t.Logf("mysql grammar compile result sql: %s ,binds %+v, err %+v", s2, v2, err2)
	s2, v2, err2 = b.Insert(&tmp)
	t.Logf("mysql grammar compile result sql: %s ,binds %+v, err %+v", s2, v2, err2)

}

func TestBaseGrammar_CompileDelete(t *testing.T) {
	b := builder.New("table_name")
	b.SetGrammar("mysql")
	b.Where("c1", "=", 0)
	b.OrWhere("c1", "=", 0)
	b.WhereRaw("c1 = a1")
	b.WhereRaw("1 != 1")
	b.WhereRaw("c1 = a1 + ?", 1)
	b.WhereNotIn("a", []int{1, 2, 3})
	b.WhereIn("a", []string{"c", "b"})
	b.OrWhereBetween("age", 8, 18)
	b.WhereNotBetween("sex", "boy", "girl")
	b.OrWhereRaw("name in (select name from users) as user_table")
	s1, v1, err1 := b.Delete()
	t.Logf("mysql grammar compile result sql: %s ,binds %+v, err %+v", s1, v1, err1)
	if s1 != "delete from table_name where  c1 = ? or c1 = ? and c1 = a1 and 1 != 1 and c1 = a1 + ? and a not in (?, ?, ?) and a  in (?, ?) or age  between ? and ? and sex not between ? and ? or name in (select name from users) as user_table " {
		log.Fatal("expected sql ", s1)
	}

	if len(v1) != 11 {
		log.Fatal("expected val", v1)
	}
}
