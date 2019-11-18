package builder

import (
	"testing"
)


var b = New("table_name")
func TestFunBuilder_Where(t *testing.T) {
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

	t.Logf("wheres: %+v",b.GetWheres())
}
func TestFunBuilder_Group(t *testing.T) {
	b.Group("a","b")
	t.Logf("groups: %+v",b.GetGroups())
}
func TestFunBuilder_Having(t *testing.T) {
	b.Having("a","like","%some%")
	b.Having("b","=",1)
	b.HavingBetween("")
	b.HavingRaw("1 != 1")
	b.HavingRaw("c1 = a1 + ?",1)
	b.OrHavingBetween("age",8,18)
	b.HavingNotBetween("sex",[]string{"boy","girl"})
	b.OrHavingRaw("name in (select name from users) as user_table")
	t.Logf("havings: %+v",b.GetHavings())
}

func TestFunBuilder_OrderBy(t *testing.T) {
	b.OrderBy("a")
	b.OrderByDesc("c")
	t.Logf("orders: %+v",b.GetOrders())
}

func TestFunBuilder_Limit(t *testing.T) {
	b.Limit(1)
	t.Logf("limit: %d",b.GetLimit())
}

func TestFunBuilder_Offset(t *testing.T) {
	b.Offset(2)
	t.Logf("offset: %d",b.GetOffset())
}
