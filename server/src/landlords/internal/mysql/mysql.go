package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//数据库配置
const (
	userName = "root"
	password = "121212"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "landlords"
)

var (
	db                = &sql.DB{}
	err         error = nil
	createTable       = `CREATE TABLE IF NOT EXISTS users(
		id INT(10) PRIMARY KEY AUTO_INCREMENT NOT NULL,
		name VARCHAR(64),
		password VARCHAR(64),
		score INT(5) DEFAULT 0,
		status INT(4) DEFAULT 0,
		createtime BIGINT DEFAULT 0
		); `
)

type User struct {
	ID         int64
	NAME       string
	PASSWORD   string
	SCORE      int8
	STATUS     int
	CREATEDATE int64
}

func init() {
	// 构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println("path", path)
	// 打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	db, _ = sql.Open("mysql", path)

	// 设置数据库最大连接数
	db.SetConnMaxLifetime(100)

	// 设置数据库最大闲置连接数
	db.SetMaxIdleConns(10)

	// 验证连接
	err = db.Ping()
	errorHandler(err)

	// 创建表
	_, err := db.Exec(createTable)
	if err == nil {
		fmt.Println("create table users successd")
		// 关闭数据库
		// defer db.Close()
	} else {
		errorHandler(err)
	}

}

func errorHandler(err error) {
	if err != nil {
		// fmt.Println(err.Error())
		panic(err.Error())
	}
}

var INSERT_DATA = `INSERT INTO users(id,name,password,createtime) VALUES(?,?,?,?);`

// InsertUser 新增用户
func InsertUser(id int, name, password string) {
	println(id)
	result, err := db.Exec(INSERT_DATA, id, name, password, int32(time.Now().Unix()))
	if err != nil {
		fmt.Printf("Insert data failed, err:%v \n", err)
		return
	}
	// 获取插入数据的自增ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Get insert id failed, err:%v \n", err)
		return
	}
	fmt.Println("Insert data id:", lastInsertID)
	// 通过RowsAffected获取受影响的行数
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed, err:%v \n", err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
}

var UPDATE_DATA = `UPDATE users SET score=28 WHERE name="唐僧";`

// Update 修改数据
func Update() {
	db.Exec(UPDATE_DATA)

}

var DELETE_DATA = `DELETE FROM users WHERE score>=30`

// Delete 删除记录
func Delete() {
	db.Exec(DELETE_DATA)
}

var DELETE_TABLE = `DROP TABLE users;`

// DeleteTable 删除表
func DeleteTable() {
	db.Exec(DELETE_TABLE)
}

var QUERY_DATA = `SELECT * FROM users;`

// QueryOne 查询数据
func QueryOne(account string) (user *User, err error) {
	user = new(User) // 用new()函数初始化一个结构体对象
	row := db.QueryRow("SELECT * FROM users WHERE name=?", account)
	// row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	if err := row.Scan(&user.ID, &user.NAME, &user.PASSWORD, &user.SCORE, &user.STATUS, &user.CREATEDATE); err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return user, err
	}
	return user, err
}

// Query 查询数据
func Query(key string) {
	if key != "" {
		QUERY_DATA = QUERY_DATA + " WHERE name=" + key
	}
	rows, err := db.Query(QUERY_DATA)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var name string
		var id int
		var score int
		if err := rows.Scan(&id, &name, &score); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s is %d\n", name, score)
	}
}

func main() {
	// 建立数据连接
	// db := setupConnect()
	// 创建数据库表
	// CreateTable(db, CREATE_TABLE)
	// 插入数据
	// Insert(db)
	// // 查询数据
	// Query(db)
	// // 删除数据
	// Delete(db)
	// // 插入数据
	// Insert(db)
	// // 修改数据
	// Update(db)
	// // 查询数据
	// Query(db)
	// // 删除表
	// DeleteTable(db)
	// 关闭数据库连接
	// db.Close()
}
