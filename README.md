funsql 
==========
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/LibiChai/funsql.svg?branch=master)](https://travis-ci.org/LibiChai/funsql)
[![Go Report Card](https://goreportcard.com/badge/github.com/libi/funsql)](https://goreportcard.com/report/github.com/libi/funsql)


## 简介
**funsql**是一个轻量级的sqlbuilder/scan包，支持函数链式调用且函数名与sql语法相同，可以非常简易的使用go语法进行sql拼装，
目前已经支持大部分的sql操作语句。

funsql拼装后的结果包含sql语句和绑定参数两部分，可以非常方便的将这两个参数传入标准sql包进行执行。

增加了scan支持，可将标准库sql包返回的rows直接绑定到自定义的结构体切片内。

## 快速开始
```go
sqlStr,args,err := funsql.Table("users").Where("age",">",10).Select()
row,err := sql.Query(sqlStr,args...)
//输出sql 与 args 可以直接作为sql标准库Query()的入参。
//sql: select * from users where age > ?
//binds: []int{10,}
//err: nil

```
## 入口函数 Table()
为了简化语句，每次sql拼装都必须首先执行入口函数。入口函数会返回*FunBuilder实例作为后续链式操作de基础对象。
入口函数必须传入当前sql的表名。可选传入语法解析器，默认为mysql（目前只实现了mysql）。

需要使用时直接调用即可。

```go
funsql.Table("users")

//需要指定语法解析器时
funsql.Table("users","mysql")
```

## 结果函数 
funsql 使用链式调用进行sql拼装，但是最终都必须以 **结果函数** 作为链式的结果才能返回拼装好的sql语句与绑定参数。

结果函数包括([x]为todo项目):
- Select 
- SelectRaw [x]
- Update 
- Delete 
- Insert

```go
funsql.Table("users").Select("name","age")
```

## Select 
在调用Select不传入操作时默认为 * 。当然大部分时候并不希望返回数据库表的所有列，此时可以指定需要的字段即可。
字段名必须为字符串类型。

查询用户表的所有用户的用户名
```go
funsql.Table("users").Select("name")
```

## Join
构造Join查询只需要传入2个参数 分别为表名 ，条件sql

```go
funsql.Table("users").Join("order","users.id = order.uid").Select()
```

## Where
### 简单语句

需要构造简单where时需要传入3个参数，分别为字段名，数据库支持的运算符，需要对比的值。

例如查询一个年龄大于10岁的用户
```go
funsql.Table("users").Where("age",">",10).Select()
```

多个条件语句可以直接进行链式调用。例如查询一个年龄大于10岁，性别为男的用户
```go
funsql.Table("users").Where("age",">",10).Where("sex","=","man").Select()
```

### Or
默认的多个where之间默认为and，如果需要使用or时使用or语句。在Where函数之前增加Or即可。

查询一个年龄为10岁或者20岁的的用户
```go
funsql.Table("users").Where("age","=",10).OrWhere("age","=",20).Select()
```

### WhereRaw / OrWhereRaw
有时后sql拼装无法满足需求，需要原生sql条件时可以使用WhereRaw / OrWhereRaw

复杂嵌套查询
```go
funsql.Table("users").WhereRaw("age/2 = 10 and name in (select created_by from books) as t")
```

WhereRaw同时支持自动参数传入，此处需要注意传入的额外参数数量需要与原生sql内的标识符数量一致
```go
funsql.Table("users").WhereRaw("age/2 = 10 and name in (select created_by from books where id = ?) as t",1)
```



### 更多where

#### WhereIn

查询字段值包含在指定的数组内
```go
funsql.Table("users").WhereIn("age",[]int{10,20,28}).Select()
```

#### WhereNotIn
查询字段值不包含在指定的数组内
```go
funsql.Table("users").WhereNotIn("age",[]int{10,20,12}).Select()
```

#### WhereBetween
查询字段值在某一范围内
```go
funsql.Table("users").WhereBetween("age",10,20)..Select()
```

#### WhereNotBetween
查询字段值不在某一范围内
```go
funsql.Table("users").WhereNotBetween("age","=",10).OrWhere("age","=",20).Select()
```


## Group By / Having
当需要对结果进行分组时可以使用GroupBy函数进行分组，使用Having（语法与Where类似）对分组结果进行查询。

```go
funsql.Table("users").WhereNotBetween("age","=",10).OrWhere("age","=",20).
GroupBy("sex","age").Having("age",">",10).Select()
```

## Order By 
需要对结果进行排序时使用OrderBy和OrderByDesc进行升序或者降序排序，可使用多个排序规则

例如先使用年龄升序排序再使用性别降序排序

```go
funsql.Table("users").WhereNotBetween("age","=",10).OrWhere("age","=",20).
OrderBy("age").OrderByDesc("sex").Select()
```



## Limit / Offset
需要限制结果数量或者返回指定offset的数据时，可以使用Limit 和 Offset，可以单独或者组合使用.

```go
funsql.Table("users").WhereNotBetween("age","=",10).OrWhere("age","=",20).
GroupBy("sex","age").Having("age",">",10).Limit(5).Offset(3).Select()
```

## Scan使用
Scan支持二维结构体切片与一维切片，适用于查询多条多字段数据和多条单个字段数据。结构体可选使用tag fs作为字段标识，不存在tag fs时使用小写形式的字段名。

### 多条多字段使用
```go
type Order struct {
	ID      int64  `fs:"id"`
	OrderNo string `fs:"order_no"`
}
orders := make([]Order, 0)
err = Scan(rows, &orders)

```

### 多条单个字段

```go
orderIDs := make([]int64, 0)
err = Scan(rows, &orderIDs)

```




