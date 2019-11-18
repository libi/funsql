package funsql

import (
	"testing"
)

func TestTable(t *testing.T) {
	table := Table("a","mysql")
	sql,vals,err := table.
		Where("c",">",0).
		WhereIn("a",[]int{1,2,3}).
		Group("a","b").
		Having("a","=",1).
		OrHaving("c","<=",2).
		Where("c1","=",0).
		OrWhere("c1","=",0).
		WhereRaw("c1 = a1").
		WhereRaw("1 != 1").
		WhereRaw("c1 = a1 + ?",1).
		WhereNotIn("a",[]int{1,2,3}).
		WhereIn("a",[]string{"c","b",}).
		OrWhereBetween("age",8,18).
		WhereNotBetween("sex","boy","girl").
		OrWhereRaw("name in (select name from users) as user_table").
		Limit(1).
		Offset(2).

		Select("a","b","c")
	t.Logf("sql: %s binds: %+v error: %+v",sql,vals,err)
}
