package mysql

import (
	"fmt"
	"time"
)

type User struct {
	ID         int
	NAME       string
	PASSWORD   string
	SCORE      int8
	STATUS     int
	CREATEDATE int64
}

var INSERT_USER = `INSERT INTO users(id,name,password,createtime) VALUES(?,?,?,?);`

// InsertUser 新增用户
func InsertUser(id int, name, password string) {
	println(id)
	result, err := db.Exec(INSERT_USER, id, name, password, int32(time.Now().Unix()))
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

// QueryOneUser 查询用户
func QueryOneUser(account string) (user *User, err error) {
	user = new(User) // 用new()函数初始化一个结构体对象
	row := db.QueryRow("SELECT * FROM users WHERE name=?", account)
	// row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	if err := row.Scan(&user.ID, &user.NAME, &user.PASSWORD, &user.SCORE, &user.STATUS, &user.CREATEDATE); err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return user, err
	}
	return user, err
}
