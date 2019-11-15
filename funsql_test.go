package funsql

import (
	"fmt"
	"testing"
)
import _ "github.com/LibiChai/funsql/builder/grammars"

func TestTable(t *testing.T) {
	table := Table("a","mysql")
	sql,vals,err := table.Where("c",">",0).Where("a","=","b").Select("a","b","c")
	fmt.Println(sql,vals,err)
}
