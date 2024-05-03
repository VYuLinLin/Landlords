package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// 数据库配置
const (
	userName = "root"
	password = "123456"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "landlords"
)

var (
	db        = &sql.DB{}
	err error = nil
)

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

	fmt.Println("连接数据库成功")

	CreateTableUsers()
}

func errorHandler(err error) {
	if err != nil {
		panic(err.Error())
	}
}
