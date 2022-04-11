第三方sqlx库操作MySQL数据库

描述: 在项目中我们通常可能会使用database/sql(原生库)连接MySQL数据库,而我们可以采用sqlx来替代它, 它可以简化代码量、从而使查询更加方便。

你可以认为sqlx是Go语言内置database/sql的超集，它在优秀的内置database/sql基础上提供了一组扩展。
例如: 常用来查询的 Get(dest interface{}, ...) error 和 Select(dest interface{}, ...) error 外还有很多其他强大的功能。

本文借助使用sqlx实现批量插入数据的例子，介绍了sqlx中可能被你忽视了的sqlx.In和DB.NamedExec方法。

第三方sqlx库主页: http://jmoiron.github.io/sqlx/

sqlx 安装&语法

描述: 在shell或者cmd终端中执行如下命令进行sqlx的安装: **go get github.com/jmoiron/sqlx**

语法原型:

> func (db *DB) Get(dest interface{}, query string, args ...interface{}) error: 执行SQL并绑定单行结果查询到指定类型变量中, 占位符参数都将替换为提供的参数，如果结果集为空则返回错误。
> func (db *DB) Select(dest interface{}, query string, args ...interface{}) error : 执行SQL并绑定多行结果查询到指定类型变量中。
> func (db *DB) Exec(query string, args ...interface{}) (Result, error): Exec 执行查询时不返回任何行但可以获取影响的行数, 支持插入、更新、删除等SQL语句


sqlx 数据库初始化

描述: 我们可以利用下述示例看到sqlx与sql之间的小小区别.

```go
// weiyigeek.top/studygo/Day09/MySQL/mypkg/initsqlx.go
package mypkg
import (
	"fmt"

	"github.com/jmoiron/sqlx"

)

// 定义一个MysqlObj结构体
type SqlObj struct {
	Mysql_host             string
	Mysql_port             uint16
	Mysql_user, Mysql_pass string
	Database               string
	DB                     *sqlx.DB
}

// 定一个Person结构体
type User struct {
	Uid  int
	Name string
	Age  int
}

func (conn *SqlObj) InitDB() (err error) {
	// DSN(Data Source Name) 数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", conn.Mysql_user, conn.Mysql_pass, conn.Mysql_host, conn.Mysql_port, conn.Database)
	// 注册第三方mysql驱动到sqlx中并连接到dsn数据源设定的数据库中(与database/sql不同点，代码更加精简)
	conn.DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("Connect %s DB Failed\n%v \n", dsn, err)
		return err
	}
	// 设置与数据库建立连接的最大数目
	conn.DB.SetMaxOpenConns(1024)
	// 设置连接池中的最大闲置连接数
	conn.DB.SetMaxIdleConns(10)
	return nil
}
```


sqlx CRUD操作

描述: 在测试使用sqlx针对MySQL数据库进行CRUD操作时,我们需要准备
库表准备

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


单行查询结果

```go
// ## 查询显示单行数据示例代码
func queryRow(db *sqlx.DB) {
	// User 结构体类型声明
	var u mypkg.User
	sqlStr := "SELECT uid,name,age FROM user WHERE uid=?"
	// 执行查询语句并通过反射reflect将查询结果进行一一绑定,返回单行数据
	err := db.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.Uid, u.Name, u.Age)
}
```


多行查询结果

```go
// ## 查询显示多行数据示例代码
func queryMultiRow(db *sqlx.DB) {
	// User 结构体类型数组声明
	var u []mypkg.User
	sqlStr := "select uid, name, age from user where uid > ?"

	// 执行多行数据结果查询
	err := db.Select(&u, sqlStr, 8)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	fmt.Printf("users:%#v\n", u)

}
```


插入数据示例

```go
// 插入数据方法示例
func insertRow(db *sqlx.DB) {
	sqlStr := "insert into user(name, age) values (?,?)"
	// EXEC 方法执行的SQL语句包括 插入/更新和删除
	ret, err := db.Exec(sqlStr, "我爱学Go", 19)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}
```


更新数据示例

```go
// 更新数据
func updateRow(db *sqlx.DB) {
	sqlStr := "update user set age=? where uid = ?"
	ret, err := db.Exec(sqlStr, 39, 8)
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

删除数据示例

```go
// 删除数据
func deleteRow(db *sqlx.DB) {
	sqlStr := "delete from user where uid = ?"
	ret, err := db.Exec(sqlStr, 16)
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

主函数数据库初始化和调用示例

```go
func main() {
	// 1.sqlx结构体初始化
	conn := &mypkg.SqlObj{
		Mysql_host: "10.20.172.248",
		Mysql_port: 3306,
		Mysql_user: "root",
		Mysql_pass: "weiyigeek.top",
		Database:   "test",
	}

	// 2.连接数据库初始化操作
	err := conn.InitDB()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("[INFO] - 数据库已连接成功！")
	}
	
	// 3.关闭sqlx.DB数据连接对象(资源释放)
	defer conn.DB.Close()
	
	// 4.单行数据查询
	fmt.Println("单行数据查询结果：")
	queryRow(conn.DB)
	
	// 5.多行数据查询
	fmt.Println("多行数据查询结果：")
	queryMultiRow(conn.DB)
	
	// 6.插入数据
	fmt.Println("输入数据操作：")
	insertRow(conn.DB)
	
	// 7.更新数据
	fmt.Println("更新数据操作: ")
	updateRow(conn.DB)
	
	// 8.删除数据
	fmt.Println("删除数据操作: ")
	deleteRow(conn.DB)

}
```


执行结果:

```sh
[INFO] - 数据库已连接成功！
单行数据查询结果：
id:1 name:WeiyiGeek age:20
多行数据查询结果：
users:[]mypkg.User{mypkg.User{Uid:16, Name:"我爱学Go", Age:19}}
输入数据操作：
insert success, the id is 17.
更新数据操作:
update success, affected rows:1
删除数据操作:
delete success, affected rows:1

# 查看数据库中存储的数据结果

select uid,name,age from `test`.`user`
1	WeiyiGeek	20
2	Elastic	18
3	Logstash	20
4	Beats	10
5	Kibana	19
6	C	25
7	C++	25
8	Python	39
17	我爱学Go	19
```


sqlx 绑定SQL语句到同名字段

我们可以使用 DB.NamedExec 和 DB.NamedQuery 方法用来绑定SQL语句与结构体或map中的同名字段，来分别进行操作字段里面的值或者将查询的结果赋予这些字段。

函数原型:

```go
func (db *DB) NamedQuery(query string, arg interface{}) (*Rows, error) - 执行查询语句返回*rows类型的数据。
func (db *DB) NamedExec(query string, arg interface{}) (sql.Result, error)  - 执行操作语句单行sql.Result结果集。
```


示例演示:

```go
// # NamedQuery
func namedQuery(){
	sqlStr := "SELECT * FROM user WHERE name=:name"
	// 1.使用 map 做命名查询
	rows, err := db.NamedQuery(sqlStr, map[string]interface{}{"name": "WeiyiGeek"})
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return
	}
  // 2.程序结束后释放资源给连接池
	defer rows.Close()

  // 3.遍历查询结果
	for rows.Next(){
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}


	// 4.使用结构体命名查询，根据结构体字段的 db tag进行映射

  u := user{
		Name: "WeiyiGeek",
	}
	rows, err = db.NamedQuery(sqlStr, u)
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next(){
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}
}

// # NamedExec
func insertUserDemo()(err error){
	sqlStr := "INSERT INTO user (name,age) VALUES (:name,:age)"
  // 执行后不返回结果，但可以通过调用.RowsAffected()了解到影响的行数。
	_, err = db.NamedExec(sqlStr,
		map[string]interface{}{
			"name": "WeiyiGeek",
			"age": 28,
		})
	return
}
```


sqlx 事务处理

描述: 对于事务操作sqlx中为我们提供了db.Beginx()和tx.Exec()等方法。

函数原型:

```go
// Beginx开始一个事务并返回一个*sqlx.Tx而不是*sql.Tx。
func (db *DB) Beginx() (*Tx, error)
```


测试库表:

```go
// # 插入测试数据
INSERT INTO `test`.`money`(`id`, `name`, `balance`) VALUES (1, 'WeiyiGeek', 1100);
INSERT INTO `test`.`money`(`id`, `name`, `balance`) VALUES (2, '辛勤的小蜜蜂', 3800);
```


实际案例:

```go
// 事务处理
func transactionSqlx(db *sqlx.DB) (err error) {
	// 开启事务
	tx, err := db.Beginx()
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}
	// 任务执行完毕后判断是否进行rollback
	defer func() {
		if p := recover(); p != nil {
			// 回滚操作并抛出异常
			tx.Rollback()
			panic(p)
		} else if err != nil {
			fmt.Println("rollback")
			// 当错误不为nil则进行回滚操作
			tx.Rollback()
		} else {
			// 提交操作
			err = tx.Commit()
			fmt.Println("commit")
		}
	}()

	// A 用户向 B用户转账 50 rmb
	sqlStr1 := "UPDATE `money` SET balance=balance-50 WHERE id=?"
	rs, err := tx.Exec(sqlStr1, 1) // 执行更新语句
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()   // 获得影响行数
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	
	// B 接收到 A用户的转账 50 rmb
	sqlStr2 := "UPDATE `money` SET balance=balance+50 WHERE id=?;"
	rs, err = tx.Exec(sqlStr2, 2)  // 执行更新语句
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()     // 获得影响行数
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	return err

}
```

执行后结果:

```sh
# 表示已提交事务

commit

# 数据库中可以看到A和B的balance都发送了变化。
1	WeiyiGeek	1050
2	辛勤的小蜜蜂	3850
```


sqlx 批量执行

描述: sqlx 为我们提供了一个非常方便的函数sqlx.In使得我们可以批量插入,使用的函数原型格式如下:

查询占位符-bindvars

描述: 例如此处查询占位符?在内部称为bindvars（查询占位符）它非常重要, 由于通过字符串格式 database/sql 不尝试对查询文本进行任何验证, 而利用查询占位符进行预处理，可以极大的防止SQL注入攻击。

> MySQL 中使用?
> PostgreSQL 中使用枚举的1、1、2
> SQLite 中?和$1的语法都支持
> Oracle 中使用:name的语法
> Tips: 非常注意bindvars的一个常见误解是，它们用来在sql语句中插入值,它们其实仅用于参数化，不允许更改SQL语句的结构。

例如，使用bindvars尝试参数化列名或表名将不起作用：

```go
// ？不能用来插入表名（做SQL语句中表名的占位符）
db.Query("SELECT * FROM ?", "mytable")

// ？也不能用来插入列名（做SQL语句中列名的占位符）
db.Query("SELECT ?, ? FROM people", "name", "location")
```


测试表库

```go
-- 为了方便演示插入数据操作，这里创建一个user表，表结构如下：
CREATE TABLE `user` (
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(20) DEFAULT '',
    `age` INT(11) DEFAULT '0',
    PRIMARY KEY(`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- 结构体：定义一个user结构体，字段通过tag与数据库中user表的列一致。
type User struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}
```


自定义批量插入

描述: 通常拼接语句实现批量插入, 方法可能比较笨但是很好理解，就是有多少个User就拼接多少个(?, ?)。

```go
// BatchInsertUsers 自行构造批量插入的语句
func BatchInsertUsers1(users []*User, db *sqlx.DB) error {
	// 1.存放 (?, ?) 的slice
	valueStrings := make([]string, 0, len(users))
	// 2.存放values的slice
	valueArgs := make([]interface{}, 0, len(users)*2)

	// 3.遍历users准备相关数据
	for _, u := range users {
		// 此处占位符要与插入值的个数对应
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, u.Name)
		valueArgs = append(valueArgs, u.Age)
	}
	fmt.Printf("%#v\n%#v\n", valueStrings, valueArgs)
	
	// 4.自行拼接要执行的具体语句
	stmt := fmt.Sprintf("INSERT INTO user (name, age) VALUES %s",
		strings.Join(valueStrings, ","))
	fmt.Println(stmt)
	res, err := db.Exec(stmt, valueArgs...)
	if err != nil {
		fmt.Printf("Exec Batch Insert Users SQL Failed, %v\n", err)
		return err
	}
	
	// 5.输出插入成功的行函数(影响行)
	count, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Get Rows Affected Failed, %v\n", err)
		return err
	} else {
		fmt.Println("Insert Rows Affected ：", count)
		return nil
	}

}

func main() {
	// 1.sqlx结构体初始化
	conn := &mypkg.SqlObj{
		Mysql_host: "10.20.172.248",
		Mysql_port: 3306,
		Mysql_user: "root",
		Mysql_pass: "weiyigeek.top",
		Database:   "test",
	}

	// 2.连接数据库初始化操作
	err := conn.InitDB()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("[INFO] - 数据库已连接成功！")
	}
	
	// 3.关闭sqlx.DB数据连接对象(资源释放)
	defer conn.DB.Close()

  // 4.自定义利用占位符进行批量插入
	userInsert := make([]*User, 0)
	userInsert = append(userInsert, &User{Name: "WeiyiGeek-20", Age: 20})
	userInsert = append(userInsert, &User{Name: "WeiyiGeek-21", Age: 21})
	userInsert = append(userInsert, &User{Name: "WeiyiGeek-22", Age: 22})
	err = BatchInsertUsers1(userInsert, conn.DB)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("批量插入执行完毕!")
	}
}
```


执行结果:

```sh
[INFO] - 数据库已连接成功！
插入占位符:
[]string{"(?, ?)", "(?, ?)", "(?, ?)"}
[]interface {}{"WeiyiGeek-20", 20, "WeiyiGeek-21", 21, "WeiyiGeek-22", 22}
INSERT INTO user (name, age) VALUES (?, ?),(?, ?),(?, ?)
Insert Rows Affected ： 3
批量插入执行完毕!

# 数据库中的结果

19	WeiyiGeek-20	20
20	WeiyiGeek-21	21
21	WeiyiGeek-22	22
```


使用 sqlx.In 实现批量插入

描述: 我们除了使用自定义的还可以使用sqlx.In方法与NamedExec方法实现批量插入，下面我们来实践sqlx.In的批量插入。

步骤01.插入实例前提是需要我们的结构体实现driver.Valuer接口(类似于Java中的重写), 此处将字段值包装为空接口进返回。

```go
func (u User) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}
```

步骤02.使用sqlx.In实现批量插入代码如下：

```go
// BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
func BatchInsertUsers2(users []interface{}, db *sqlx.DB) error {
	// 1.预处理SQL将参数与占位符绑定。
	query, args, _ := sqlx.In(
		"INSERT INTO user (name, age) VALUES (?), (?), (?)",
		users..., // 如果arg实现了 driver.Valuer, sqlx.In 会通过调用 Value()来展开它
	)
	fmt.Println(query) // 查看生成的querystring
	fmt.Println(args)  // 查看生成的args

	// 2.执行批量插入。
	res, err := db.Exec(query, args...)
	if err != nil {
		fmt.Printf("Exec Batch Insert Users SQL Failed, %v\n", err)
		return err
	}
	
	// 3.输出插入成功的行函数(影响行)。
	count, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Get Rows Affected Failed, %v\n", err)
		return err
	} else {
		fmt.Println("Insert Rows Affected ：", count)
		return nil
	}

}

// 在前面的Main 函数中调用执行如下代码块，我就不再重新写了。
......
// 主要功能：使用sqlx.in进行批量插入
userInsert := make([]interface{}, 0) // 空接口数组内存申请
userInsert = append(userInsert, &User{Name: "Gooo-20", Age: 20})
userInsert = append(userInsert, &User{Name: "R-21", Age: 21})
userInsert = append(userInsert, &User{Name: "Javascript-22", Age: 22})
err = BatchInsertUsers2(userInsert, conn.DB)
if err != nil {
  panic(err) // 在进行开发测试代码时使用，正式环境中请勿使用。
} else {
  fmt.Println("sqlx.In - 批量插入执行完毕!")
}
```


执行结果:

```sh
[INFO] - 数据库已连接成功！
INSERT INTO user (name, age) VALUES (?, ?), (?, ?), (?, ?)
[Go-20 20 R-21 21 Javascript-22 22]
Insert Rows Affected ： 3
sqlx.In - 批量插入执行完毕!

# 数据库插入结果查询

25	Go-20	20
26	R-21	21
27	Javascript-22	22
```




扩展学习之 sqlx.In 的查询示例
在sqlx查询语句中实现In查询和 FIND_IN_SET函数, 即实现 SELECT * FROM user WHERE id in (3, 2, 1); 和 SELECT * FROM user WHERE id in (3, 2, 1) ORDER BY FIND_IN_SET(id, '3,2,1');.

In查询: IN 操作符允许我们在 WHERE 子句中规定多个值

```go
// QueryByIDs 根据给定ID查询
func QueryByIDs(ids []int)(users []User, err error){
  // 动态填充id
  query, args, err := sqlx.In("SELECT name, age FROM user WHERE id IN (?)", ids)
  if err != nil {
    return
  }
  // sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
  query = DB.Rebind(query)
    err = DB.Select(&users, query, args...)
    return
}
```


in 查询 和 FIND_IN_SET 查询: 查询id在给定id集合的数据并维持给定id集合的顺序。

```go
// QueryAndOrderByIDs 按照指定id查询并维护顺序
func QueryAndOrderByIDs(ids []int)(users []User, err error){
  // 动态填充id
  strIDs := make([]string, 0, len(ids))
  for _, id := range ids {
    strIDs = append(strIDs, fmt.Sprintf("%d", id))
  }
  query, args, err := sqlx.In("SELECT name, age FROM user WHERE id IN (?) ORDER BY FIND_IN_SET(id, ?)", ids, strings.Join(strIDs, ","))
  if err != nil {
    return
  }
  // sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
  query = DB.Rebind(query)

  err = DB.Select(&users, query, args...)
  return
}
```


当然，在这个例子里面你也可以先使用IN查询，然后通过代码按给定的ids对查询结果进行排序。

Tips: 上述SQL执行结果以及IN 关键字 与 FIND_IN_SET 区别如下:

```sh
-- # IN 关键字段
SELECT uid, name, age FROM user WHERE uid IN (1,2)
-- uid	name	age
-- 1	WeiyiGeek	20
-- 2	Elastic	18
SELECT uid, name, age FROM user WHERE 8 IN (uid)
-- uid	name	age
-- 8	Python	39
SELECT uid, name, age FROM user WHERE 1 IN (2,3,4)
-- 字段返回为空


-- # FIND_IN_SET 函数使用
SELECT FIND_IN_SET (5, '1,5,6,18') as 'Index';
-- Index
-- 2
SELECT uid, name, age FROM user WHERE FIND_IN_SET (1,uid);
-- uid	name	age
-- 1	WeiyiGeek	20

-- # 组合使用 : 安装顺序数组但将设定的FIND_IN_SET的uid那一行值放在末尾。
SELECT uid, name, age FROM user WHERE uid IN (1,5,6,18) ORDER BY FIND_IN_SET (1,uid);
-- uid	name	age
-- 5	Kibana	19
-- 6	C	25
-- 18	我爱学Go	19
-- 1	WeiyiGeek	20
```


使用 NamedExec 实现批量插入

注意：该功能需1.3.1版本以上并在1.3.1版本目前还有点问题sql语句最后不能有空格和, 不过当前版本 v1.3.4 中已解决;

使用NamedExec实现批量插入示例如下:

```go
// BatchInsertUsers3 使用NamedExec实现批量插入函数
func BatchInsertUsers3(users []*User, db *sqlx.DB) error {
	// 1.SQL预处理以及执行批量插入
	res, err := db.NamedExec("INSERT INTO user (name, age) VALUES (:name, :age)", users)
	if err != nil {
		fmt.Printf("Exec Batch Insert Users SQL Failed, %v\n", err)
		return err
	}

	// 2.输出插入成功的行函数(影响行)。
	count, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Get Rows Affected Failed, %v\n", err)
		return err
	} else {
		fmt.Println("Insert Rows Affected ：", count)
		return nil
	}

}

// 在 Main 函数中执行执行如下代码片段进行使用 NamedExec 实现批量插入
userInsert := make([]*User, 0)
userInsert = append(userInsert, &User{Name: "小红", Age: 20})
userInsert = append(userInsert, &User{Name: "小南", Age: 21})
userInsert = append(userInsert, &User{Name: "小白", Age: 22})
err = BatchInsertUsers3(userInsert, conn.DB)
if err != nil {
  fmt.Printf("[Error] - %v\n", err)
} else {
  fmt.Println("NamedExec - 批量插入执行完毕!")
}
```


执行结果:

```go
[INFO] - 数据库已连接成功！
Insert Rows Affected ： 3
NamedExec - 批量插入执行完毕!

# 数据库中插入的数据查看

28	小红	20
29	小南	21
30	小白	22
```


此处将上面三种方法综合起来试一下：

```go
func main() {
  err := initDB()
  if err != nil {
    panic(err)
  }
  defer DB.Close()
  u1 := User{Name: "WeiyiGeek", Age: 18}
  u2 := User{Name: "weiy_", Age: 28}
  u3 := User{Name: "weiyi", Age: 38}

  // 方法1.User类型的指针数组
  users := []*User{&u1, &u2, &u3}
  err = BatchInsertUsers(users)
  if err != nil {
    fmt.Printf("BatchInsertUsers failed, err:%v\n", err)
  }

  // 方法2.空接口类型的数组
  users2 := []interface{}{u1, u2, u3}
  err = BatchInsertUsers2(users2)
  if err != nil {
    fmt.Printf("BatchInsertUsers2 failed, err:%v\n", err)
  }

  // 方法3.User类型的指针数组
  users3 := []*User{&u1, &u2, &u3}
  err = BatchInsertUsers3(users3)
  if err != nil {
    fmt.Printf("BatchInsertUsers3 failed, err:%v\n", err)
  }
}
```


至此使用sqlx操作MyDSQL数据库实践完毕！ 