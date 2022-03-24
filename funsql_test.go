package funsql

import (
	"reflect"
	"testing"
)

func TestTable(t *testing.T) {
	table := Table("a", "mysql")
	builder := table.Join("b", "a.id = b.a_id").
		Join("c", "c.id = b.a_id").
		Where("c", ">", 0).
		WhereIn("a", []int{1, 2, 3})
	builder = builder.
		Join("b", "a.id = b.a_id").
		Join("c", "c.id = b.a_id").
		Where("c", ">", 0).
		WhereIn("a", []int{1, 2, 3}).
		Group("a", "b").
		Having("a", "=", 1).
		OrHaving("c", "<=", 2).
		Where("c1", "=", 0).
		OrWhere("c1", "=", 0).
		WhereRaw("c1 = a1").
		WhereRaw("1 != 1").
		WhereRaw("c1 = a1 + ?", 1).
		WhereNotIn("a", []int{1, 2, 3}).
		WhereIn("a", []string{"c", "b"}).
		OrWhereBetween("age", 8, 18).
		WhereNotBetween("sex", "boy", "girl").
		OrWhereRaw("name in (select name from users) as user_table").
		OrderBy("age").
		OrderByDesc("name").
		Limit(1).
		Offset(2)
	sql, vals, err := builder.Select("a", "b", "c")
	t.Logf("sql: %s binds: %#v error: %+v", sql, vals, err)
	if sql != "select a, b, c from a join b on a.id = b.a_id join c on c.id = b.a_id join b on a.id = b.a_id join c on c.id = b.a_id where  c > ? and a  in (?, ?, ?) and c > ? and a  in (?, ?, ?) and c1 = ? or c1 = ? and c1 = a1 and 1 != 1 and c1 = a1 + ? and a not in (?, ?, ?) and a  in (?, ?) or age  between ? and ? and sex not between ? and ? or name in (select name from users) as user_table  group by a,b having  a = ? or c <= ?  order by age asc, name desc limit 1 offset 2 " {
		t.Fatal("sql not match")
	}
	if !reflect.DeepEqual(vals, []interface{}{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 1, 2, 3, "c", "b", 8, 18, "boy", "girl", 1, 2}) {
		t.Fatal("vals not match")
	}

	sql, vals, err = builder.SelectRaw("count(*),a")
	t.Logf("sql: %s binds: %#v error: %+v", sql, vals, err)
	if sql != "select count(*),a  from a join b on a.id = b.a_id join c on c.id = b.a_id join b on a.id = b.a_id join c on c.id = b.a_id where  c > ? and a  in (?, ?, ?) and c > ? and a  in (?, ?, ?) and c1 = ? or c1 = ? and c1 = a1 and 1 != 1 and c1 = a1 + ? and a not in (?, ?, ?) and a  in (?, ?) or age  between ? and ? and sex not between ? and ? or name in (select name from users) as user_table  group by a,b having  a = ? or c <= ?  order by age asc, name desc limit 1 offset 2 " {
		t.Fatal("sql not match")
	}
	if !reflect.DeepEqual(vals, []interface{}{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 1, 2, 3, "c", "b", 8, 18, "boy", "girl", 1, 2}) {
		t.Fatal("vals not match")
	}

}
