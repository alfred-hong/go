快速了解 MySQL 数据库
MySQL 是目前主流关系型的数据库，它的胞胎兄弟 MariaDB (MySQL 的一个分支)，除此之外使用最多的就是 Oracle 和 PostgreSQL 数据库。

SQL 语言类型:

> DDL : 主要是操作数据库
> DML : 主要进行表的增删改查
> DCL : 主要进行用户和权限操作

MySQL 至此插件式的存储引擎，其常见存储引擎MyISAM 和 InnoDB：
MyISAM 特点:

> 查询速度快
> 只支持表锁
> 不支持事务

InnoDB 特点:

> 整体操作速度快
> 支持表锁和行锁
> 支持事务

事务的特点即我们常说的ACID:

> A(Atomicity)- 原子性 (多个语句要么全成功，要么即失败，将不会更改数据库的数据)
> C(Consistence) - 一致性 (在每次提交或回滚之后以及正在进行的事务处理期间,数据库始终保持一致状态,要么全部旧值要么全部新值)
> I(Isolation) - 隔离性 (事务之间的相互隔离的)
> D(Durability) - 持久性 (事务操作的结果是不会丢失的)

1.MySQL驱动下载

描述: Go语言中的database/sql包提供了保证SQL或类SQL数据库的泛用接口, 并不提供具体的数据库驱动, 所以使用database/sql包时必须注入（至少）一个数据库驱动。

Go语言中我们常用的数据库操作, 基本上都有完整的第三方实现，例如本节的MySQL驱动(https://github.com/go-sql-driver/mysql)

```sh
# 下载mysql驱动依赖, 第三方的依赖默认保存在 `$GOPATH/src` (注意是在项目目录里)
➜ go get -u github.com/go-sql-driver/mysql
go: downloading github.com/go-sql-driver/mysql v1.6.0
➜ weiyigeek.top go get github.com/go-sql-driver/mysql

# 项目地址
➜ weiyigeek.top pwd
/home/weiyigeek/app/program/project/go/src/weiyigeek.top

# 第三方包地址
➜ go-sql-driv pwd
/home/weiyigeek/app/program/project/go/pkg/mod/github.com/go-sql-driv
```


2.MySQL驱动格式

描述: 使用MySQL驱动格式函数原型如下所示：

> func Open(driverName, dataSourceName string) (*DB, error) : Open方法是打开一个dirverName指定的数据库，dataSourceName指定数据源，一般至少包括数据库文件名和其它连接必要的信息。
> func (db *DB) SetMaxOpenConns(n int) : SetMaxOpenConns方法是设置与数据库建立连接的最大数目。
> 如果n大于0且小于最大闲置连接数，会将最大闲置连接数减小到匹配最大开启连接数的限制。
> 如果n<=0，不会限制最大开启连接数，默认为0（无限制）。
> func (db *DB) SetMaxIdleConns(n int) : SetMaxIdleConns方法是设置连接池中的最大闲置连接数。
> 如果n大于最大开启连接数，则新的最大闲置连接数会减小到匹配最大开启连接数的限制。
> 如果n<=0，不会保留闲置连接。

基础示例

```go
import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

)

func main() {
	// 数据库DSN(Data Source Name)连接数据源
	dsn := "root:WwW.weiyigeek.top@tcp(10.20.172.248:3306)/test?charset=utf8&parseTime=True"

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("DSN : %s Format failed, Error: %v \n", dsn, err)
		panic(err)
	}
	// 此行代码要写在上面err判断的下面（注意点）。
	defer db.Close()
	
	// 判断连接的数据库
	err = db.Ping()
	if err != nil {
		fmt.Printf("Connection %s Failed, Error: %v \n", dsn, err)
		return
	}
	
	fmt.Println("数据库连接成功！")

}
```

Tips: 为什么上面代码中的defer db.Close()语句不应该写在if err != nil的前面呢？



3.MySQL初始化连接

描述: 上面的例子可以看到Open函数可能只是验证其参数格式是否正确，实际上并不创建与数据库的连接，此时我们如果要检查数据源的名称是否真实有效，应该调用Ping方法。

下述代码中sql.DB是表示连接的数据库对象（结构体实例），它保存了连接数据库相关的所有信息。它内部维护着一个具有零到多个底层连接的连接池，它可以安全地被多个goroutine同时使用。

MySQL 用户密码更改:

```go
-- MySQL 5.7.x & MySQL 8.x
ALTER USER `root`@`%` IDENTIFIED BY 'weiyigeek.top';
```

初始化示例:

```go
// Go 语言利用 MySQL Driver 连接 MySQL 示例
package main
import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

)
// 定义一个全局对象db
var db *sql.DB
// 定义一个初始化数据库的函数
func initDB() (err error) {
	// DSN(Data Source Name) - 数据库连接数据源
	// MySQL 5.7.X 与 MySQL 8.x 都是支持的
	dsn := "root:weiyigeek.top@tcp(10.20.172.248:3306)/test?charset=utf8&parseTime=True"
	// 注册第三方mysql驱动到sql中，此处并不会校验账号密码是否正确，此处赋值给全局变量db。
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("DSN : %s Format failed\n %v \n", dsn, err)
		return err
	}
	// 尝试与数据库建立连接（校验DSN是否正确）
	err = db.Ping()
	if err != nil {
		fmt.Printf("Connection %s Failed,\n%v \n", dsn, err)
		return err
	}
	// 设置与数据库建立连接的最大数目
	db.SetMaxOpenConns(1024)
	// 设置连接池中的最大闲置连接数，0 表示不会保留闲置。
	db.SetMaxIdleConns(0)
	fmt.Println("数据库初始化连接成功!")
	return nil
}

func main() {
	// 调用输出化数据库的函数
	err := initDB()
	defer db.Close()

	if err != nil {
		fmt.Println("Database Init failed!")
		return
	}

}
```

执行结果:

```sh
# 连接成功时

数据库初始化连接成功!

# 连接失败时

Connection root:www.weiyigeek.top@tcp(10.20.172.248:3306)/test?charset=utf8&parseTime=True Failed,
Error 1045: Access denied for user 'root'@'10.20.172.108' (using password: YES)
Database Init failed!
```


4.MySQL的CRUD操作

库表准备

我们首先需要在MySQL(8.x)数据库中创建一个名为test数据库和一个user表,SQL语句如下所示:

```mysql
-- 建库建表
CREATE DATABASE test;
USE test;
CREATE TABLE `user` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(20) DEFAULT '',
  `age` INT(11) DEFAULT '0',
  PRIMARY KEY(`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- 测试数据插入
INSERT INTO `test`.`user`(`uid`, `name`, `age`) VALUES (1, 'WeiyiGeek', 20);
INSERT INTO `test`.`user`(`uid`, `name`, `age`) VALUES (2, 'Elastic', 18);
INSERT INTO `test`.`user`(`uid`, `name`, `age`) VALUES (3, 'Logstash', 20);
INSERT INTO `test`.`user`(`uid`, `name`, `age`) VALUES (4, 'Beats', 10);
INSERT INTO `test`.`user`(`uid`, `name`, `age`) VALUES (5, 'Kibana', 19);
INSERT INTO `test`.`user`(`uid`, `name`, `age`) VALUES (6, 'C', 25);
INSERT INTO `test`.`user`(`uid`, `name`, `age`) VALUES (7, 'C++', 25);
INSERT INTO `test`.`user`(`uid`, `name`, `age`) VALUES (8, 'Python', 26);
```

示例结构体声明:

```go
type user struct {
	id   int
	age  int
	name string
}
```


单行查询

函数原型: func (db *DB) QueryRow(query string, args ...interface{}) *Row
函数说明: 单行查询db.QueryRow()执行一次查询，并期望返回最多一行结果（即Row）。
Tips: QueryRow总是返回非nil的值，直到返回值的Scan方法被调用时，才会返回被延迟的错误。（如：未找到结果）

简单示例:

```go
// 查询单条数据示例
func queryRowDemo() {
  var u user
	sqlStr := "select id, name, age from user where id=?"
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放 [注意点]
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
}
```


多行查询

函数原型: func (db \*DB) Query(query string, args ...interface{}) (\*Rows, error)

函数说明: 多行查询db.Query()执行一次查询，返回多行结果(即 Rows), 一般用于执行select命令, 参数args表示 query中的占位参数(空接口)。

简单示例:

```go
// 查询多条数据示例
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接 [否则将一直占有连接池资源导致后续无法正常连接]
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}

}
```


插入/更新/删除数据

函数原型: func (db *DB) Exec(query string, args ...interface{}) (Result, error)
函数说明: Exec执行一次命令（包括查询、删除、更新、插入等），返回的Result是对已执行的SQL命令的总结。参数args表示query中的占位参数。

具体插入数据示例代码如下：

```go
// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "王五", 38)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
  // 新插入数据的id
	theID, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}
```


具体更新数据示例代码如下：

```go
// 更新数据
func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 3)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}
```


具体删除数据的示例代码如下：

```go
// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 3)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}
```


综合实践

下述代码简单实现利用Go语言操作MySQL数据库的增、删、改、查等。

数据库连接封装：weiyigeek.top/studygo/Day09/MySQL/mypkg

```go
// 自定义mypkg包 initdb.go
package mypkg
import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

)

// 定义一个mysqlObj结构体
type MysqlObj struct {
	Mysql_host             string
	Mysql_port             uint16
	Mysql_user, Mysql_pass string
	Database               string
	Db                     *sql.DB
}

// 定一个Person结构体
type Person struct {
	Uid  int
	Name string
	Age  int
}

// 定义一个初始化数据库的函数
func (conn *MysqlObj) InitDB() (err error) {

	// DSN(Data Source Name) 数据库连接字符串
	// MySQL 5.7.X 与 MySQL 8.x 都是支持的
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", conn.Mysql_user, conn.Mysql_pass, conn.Mysql_host, conn.Mysql_port, conn.Database)
	
	// 注册第三方mysql驱动到sql中，此处并不会校验账号密码是否正确，此处赋值给全局变量db。
	conn.Db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("DSN : %s Format failed\n%v \n", dsn, err)
		return err
	}
	
	// 尝试与数据库建立连接（校验DSN是否正确）
	err = conn.Db.Ping()
	if err != nil {
		fmt.Printf("Connection %s Failed,\n%v \n", dsn, err)
		return err
	}
	
	// 设置与数据库建立连接的最大数目
	conn.Db.SetMaxOpenConns(1024)
	
	// 设置连接池中的最大闲置连接数
	conn.Db.SetMaxIdleConns(0) // 不会保留闲置
	
	return nil

}
```


实践 main 入口函数:

```go
package main
import (
	"database/sql"
	"fmt"
	db "weiyigeek.top/studygo/Day09/MySQL/mypkg"
)

// 单结果语句查询函数示例
func queryPersonOne(conn *sql.DB, Uid int) (res db.Person) {
	// 1.单条SQL语句
	sqlStr := `select Uid,name,age from test.user where Uid=?;`
	// 2.执行SQL语句并返回一条结果
	rowObj := conn.QueryRow(sqlStr, Uid)
	// 3.必须对rowObj调用Scan方法，因为查询后我们需要释放数据库连接对象，而它调用后会自动释放。
	rowObj.Scan(&res.Uid, &res.Name, &res.Age)
	// 4.返回一个person对象
	return res
}

// 多结果语句查询函数示例
func queryPersonMore(conn *sql.DB, id int) {
	// 1.SQL 语句
	sqlStr := `select Uid,name,age from test.user where Uid > ?;`
	// 2.执行 SQL
	rows, err := conn.Query(sqlStr, id)
	if err != nil {
		fmt.Printf("Exec %s query failed！,err : %v \n", sqlStr, err)
		return
	}
	// 3.调用结束后关闭rows，释放数据库连接资源
	defer rows.Close()
	// 4.循环读取结果集中的数据
	for rows.Next() {
		var u db.Person
		err := rows.Scan(&u.Uid, &u.Name, &u.Age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("Uid:%d name:%s age:%d\n", u.Uid, u.Name, u.Age)
	}
}

// 执行插入操作的函数示例
func insertPerson(conn *sql.DB) {
	// 1.SQL 语句
	sqlStr := `insert into user(name,age) values("Go语言",15)`
	// 2.执行插入语句
	ret, err := conn.Exec(sqlStr)
	if err != nil {
		fmt.Printf("Insert Failed, err : %v \n", err)
		return
	}
	// 3.插入数据操作，拿到插入数据库的id值
	Uid, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("Get Id Failed, err : %v \n", err)
		return
	}
	// 4.打印插入数据的id值
	fmt.Println("插入语句Uid值: ", Uid)
}

// 执行更新操作的函数示例
func updatePerson(conn *sql.DB, age, Uid int) {
	// 1.SQL 语句
	sqlStr := `update user set age=? where Uid = ?`
	// 2.执行插入语句
	ret, err := conn.Exec(sqlStr, age, Uid)
	if err != nil {
		fmt.Printf("Update Failed, err : %v \n", err)
		return
	}
	// 3.更新数据操作，获取到受影响的行数
	count, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("Get Id Failed, err : %v \n", err)
		return
	}
	// 4.打印数据影响的行数
	fmt.Println("更新数据影响的行数: ", count)
}

// 执行删除数据的操作函数示例
func deletePerson(conn *sql.DB, Uid int) {
	// 1.SQL 语句
	sqlStr := `delete from user where Uid > ?`
	// 2.执行删除的语句
	ret, err := conn.Exec(sqlStr, Uid)
	if err != nil {
		fmt.Printf("Delete Failed, err : %v \n", err)
		return
	}
	// 3.删除数据操作，获取到受影响的行数
	count, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("Get Id Failed, err : %v \n", err)
		return
	}
	// 4.打印删除数据的影响的行数:
	fmt.Println("删除数据影响的行数: ", count)
}

func main() {
	// 1.mysqlObj 结构体实例化
	conn := &db.MysqlObj{
		Mysql_host: "10.20.172.248",
		Mysql_port: 3306,
		Mysql_user: "root",
		Mysql_pass: "weiyigeek.top",
		Database:   "test",
	}
	// 2.初始化数据库
	err := conn.InitDB()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("数据库初始化连接成功!")
	}

	// 3.程序结束时关闭数据库连接
	defer conn.Db.Close()
	
	// 4.单行查询
	res := queryPersonOne(conn.Db, 1)
	fmt.Printf("单行查询: %#v\n", res)
	
	// 5.多行查询
	fmt.Println("多行查询")
	queryPersonMore(conn.Db, 6)
	
	// 6.插入数据
	fmt.Println("插入数据")
	insertPerson(conn.Db)
	
	// 7.更新数据
	fmt.Println("更新数据")
	updatePerson(conn.Db, 16, 9)
	
	// 8.删除数据
	fmt.Println("删除数据")
	deletePerson(conn.Db, 10)

}
```

执行结果&数据库查询结果:

```sh
数据库初始化连接成功!
单行查询: main.person{uid:1, name:"WeiyiGeek", age:20}
多行查询
uid:7 name:C++ age:25
uid:8 name:Python age:26
uid:9 name:Golang age:15
插入数据
插入语句uid值:  10
更新数据
更新数据影响的行数:  1
删除数据
删除数据影响的行数:  1
```

![img](https://i0.hdslb.com/bfs/article/1f2f90c82962987eeecfc2ce2a6d076d9263c609.png@801w_773h_progressive.png)

5.MySQL预处理



基础介绍

什么是预处理？

普通SQL语句执行过程：

> 客户端对SQL语句进行占位符替换得到完整的SQL语句。
> 客户端发送完整SQL语句到MySQL服务端
> MySQL服务端执行完整的SQL语句并将结果返回给客户端。

预处理执行过程：

> 把SQL语句分成两部分，命令部分与数据部分。
> 先把命令部分发送给MySQL服务端，MySQL服务端进行SQL预处理。
> 然后把数据部分发送给MySQL服务端，MySQL服务端对SQL语句进行占位符替换。
> MySQL服务端执行完整的SQL语句并将结果返回给客户端。


为什么要预处理？

> 优化MySQL服务器重复执行SQL的方法，可以提升服务器性能，提前让服务器编译，一次编译多次执行，节省后续编译的成本。
> 避免SQL注入问题。

SQL注入

描述: 非常注意, 我们任何时候都不应该自己拼接SQL语句, 可能会导致SQL注入的问题。

此处演示一个自行拼接SQL语句的示例，编写一个根据name字段查询user表的函数如下：

```go
// 可被 sql 注入示例
func sqlInjectDemo(name string) {
  var u user
	sqlStr := fmt.Sprintf("select id, name, age from user where name='%s'", name)  // 关键点
	fmt.Printf("SQL:%s\n", sqlStr)
	err := db.QueryRow(sqlStr).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	fmt.Printf("user:%#v\n", u)
}
```

当name变量输入以下字符串时便会引发SQL注入问题:

```sh
sqlInjectDemo("xxx' or 1=1#")
sqlInjectDemo("xxx' union select * from user #")
sqlInjectDemo("xxx' and (select count(*) from user) <10 #")
```


示例演示

Go是如何实现MySQL预处理
描述: database/sql 中使用下面的Prepare方法来实现预处理操作。
函数原型: func (db *DB) Prepare(query string) (*Stmt, error)
函数说明: Prepare方法会先将sql语句发送给MySQL服务端，返回一个准备好的状态用于之后的查询和命令。返回值可以同时执行多个查询和命令。

示例演示:
描述: 此处引用上面封装的结构体成员以及方法,进行数据库的初始化操作。

```go
package main

import (
	"database/sql"
	"fmt"

	db "weiyigeek.top/studygo/Day09/MySQL/mypkg"

)

// ## 预处理查询示例函数
func prepareQuery(conn *sql.DB, id int) {
	// SQL语句
	sqlStr := "select uid,name,age from user where uid > ?;"
	// 预处理
	stmt, err := conn.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	// 释放预处理
	defer stmt.Close()

	// 查询 uid 为 id 以上的数据
	rows, err := stmt.Query(id)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 释放 rows
	defer rows.Close()
	
	// 循环读取结果集中的数据，此处利用map来装我们遍历获取到的数据，注意内存申请。
	res := make(map[int]db.Person, 5)
	for rows.Next() {
		var u db.Person
		err := rows.Scan(&u.Uid, &u.Name, &u.Age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		_, ok := res[u.Uid]
		if !ok {
			res[u.Uid] = u
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.Uid, u.Name, u.Age)
	}
	fmt.Printf("%#v\n", res)

}


// ## 插入、更新和删除操作的预处理十分类似，这里以插入操作的预处理为例：
func prepareInsert(conn *sql.DB) {
	// 插入的SQL语句
	sqlStr := "insert into user(name,age) values (?,?)"
	// 进行SQL语句的预处理
	stmt, err := conn.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	// 释放 stmt 资源
	defer stmt.Close()

	// 执行预处理后的SQL (可以多次执行)
	_, err = stmt.Exec("WeiyiGeek", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	// 执行预处理后的SQL
	_, err = stmt.Exec("插入示例", 82)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	// 插入成功会显示如下
	fmt.Println("insert success.")

}

// 入口函数
func main() {
	// MysqlObj 结构体初始化
	conn := &db.MysqlObj{
		Mysql_host: "10.20.172.248",
		Mysql_port: 3306,
		Mysql_user: "root",
		Mysql_pass: "weiyigeek.top",
		Database:   "test",
	}
	// 数据库初始化
	err := conn.InitDB()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("[INFO] - 已成功连接到数据库!")
	}
	// 关闭数据库对象
	defer conn.Db.Close()

	// 预处理查询
	fmt.Println("预处理查询示例函数 prepareQuery:")
	prepareQuery(conn.Db, 5)
	
	// 预处理插入
	fmt.Println("预处理插入示例函数 prepareInsert:")
	prepareInsert(conn.Db)

}
```

执行结果:

```sh
[INFO] - 已成功连接到数据库!

-- 预处理查询示例函数 prepareQuery:
id:6 name:C age:25
id:7 name:C++ age:25
id:8 name:Python age:26
id:9 name:Golang age:16
id:12 name:WeiyiGeek age:18
id:13 name:插入示例 age:82
map[int]mypkg.Person{6:mypkg.Person{Uid:6, Name:"C", Age:25}, 7:mypkg.Person{Uid:7, Name:"C++", Age:25}, 8:mypkg.Person{Uid:8, Name:"Python", Age:26}, 9:mypkg.Person{Uid:9, Name:"Golang", Age:16}, 12:mypkg.Person{Uid:12, Name:"WeiyiGeek", Age:18}, 13:mypkg.Person{Uid:13, Name:"插入示例", Age:18}}

-- 预处理插入示例函数 prepareInsert:
insert success.
```

Tips：不同的数据库中，SQL语句使用的占位符语法不尽相同，例如下表所示。

![img](https://i0.hdslb.com/bfs/article/ba914aeb3ff57be5f81c22264a1c7d36575c9755.png@369w_318h_progressive.png)


6.MySQL事务处理

什么是事务？

事务：一个最小的不可再分的工作单元；通常一个事务对应一个完整的业务(例如银行账户转账业务，该业务就是一个最小的工作单元)，同时这个完整的业务需要执行多次的DML(insert、update、delete)语句共同联合完成。A转账给B，这里面就需要执行两次update操作。

在MySQL中只有使用了Innodb数据库引擎的数据库或表才支持事务, 事务处理可以用来维护数据库的完整性，保证成批的SQL语句要么全部执行，要么全部不执行。



事务特性复习 ACID

描述: 通常事务必须满足4个条件（ACID）：原子性（Atomicity，或称不可分割性）、一致性（Consistency）、隔离性（Isolation，又称独立性）、持久性（Durability）。

> 原子性: 一个事务（transaction）中的所有操作，要么全部完成，要么全部不完成，不会结束在中间某个环节。事务在执行过程中发生错误，会被回滚（Rollback）到事务开始前的状态，就像这个事务从来没有执行过一样。
> 一致性: 在事务开始之前和事务结束以后，数据库的完整性没有被破坏。这表示写入的资料必须完全符合所有的预设规则，这包含资料的精确度、串联性以及后续数据库可以自发性地完成预定的工作。
> 隔离性: 数据库允许多个并发事务同时对其数据进行读写和修改的能力，隔离性可以防止多个事务并发执行时由于交叉执行而导致数据的不一致。事务隔离分为不同级别，包括读未提交（Read uncommitted）、读提交（read committed）、可重复读（repeatable read）和串行化（Serializable）。
> 持久性: 事务处理结束后，对数据的修改就是永久的，即便系统故障也不会丢失。


事务方法原型
描述：Go语言中使用以下三个方法实现MySQL中的事务操作。

> func (db *DB) Begin() (*Tx, error) : 开始事务
> func (tx *Tx) Commit() error : 提交事务
> func (tx *Tx) Rollback() error : 回滚事务


实践示例
描述: 下面的代码演示了一个简单的事务操作，该事物操作能够确保两次更新操作要么同时成功要么同时失败，不会存在中间状态。
例如: A 转账给 B 50 RMB,即从A账号余额-50,B账号余额+50。

数据库表创建:

```go
-- 测试表
create table `money` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(20) DEFAULT '',
  `balance` INT(16) DEFAULT '0',
  PRIMARY KEY(`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- 测试数据
insert into `test`.`money`(`name`,`balance`) values("WeiyiGeek",1200);
insert into `test`.`money`(`name`,`balance`) values("辛勤的小蜜蜂",3650);

-- 查看插入的测试数据
SELECT * from money;
1	WeiyiGeek	1200
2	辛勤的小蜜蜂	3650
```

示例代码:

```go
package main
import (
	"database/sql"
	"fmt"

	"weiyigeek.top/studygo/Day09/MySQL/mypkg"

)

// ## 事务操作示例
func transactionDemo(conn *sql.DB, money int) {
	// 开启事务
	tx, err := conn.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}

	// (1) A 用户转账 50 给 B 则 - 50
	sqlStr1 := "UPDATE `money` SET balance=balance-? WHERE id=?;"
	ret1, err := tx.Exec(sqlStr1, money, 1)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}
	
	// B 用户接收到 A 转账的 50 给 则 + 50
	sqlStr2 := "UPDATE `money` SET balance=balance+? WHERE id=?;"
	ret2, err := tx.Exec(sqlStr2, money, 2)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql2 failed, err:%v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}
	
	// 事务处理影响行数判断是否修改成功
	fmt.Println("事务处理影响行数判断是否修改成功: ", affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务正在提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚啦...")
	}
	
	fmt.Println("[INFO] - 事务完成了 ，exec trans success!")

}

func main() {
	// (1) MysqlObj 结构体初始化
	conn := &mypkg.MysqlObj{
		Mysql_host: "10.20.172.248",
		Mysql_port: 3306,
		Mysql_user: "root",
		Mysql_pass: "weiyigeek.top",
		Database:   "test",
	}

	// (2) 数据库初始化
	err := conn.InitDB()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("[INFO] - 已成功连接到数据库!")
	}
	// 关闭数据库对象
	defer conn.Db.Close()
	
	// (3) 简单的事务操作示例
	transactionDemo(conn.Db, 50)

}
```

执行结果:

```sh
[INFO] - 已成功连接到数据库!
事务处理影响行数判断是否修改成功:  1 1
事务正在提交啦...
[INFO] - 事务完成了 ，exec trans success!

# 可以看到用户的在数据库中金额变化

1	WeiyiGeek	1150
2	辛勤的小蜜蜂	3700
```


至此使用database/sql标准库操作MySQL数据库完毕！ 