package mysql

import (
	"fmt"
	"time"
)

func errorHandler(err error) {
	if err != nil {
		panic(err.Error())
	}
}

/*
 * Example
 */

var INSERT_DATA = `INSERT INTO users(id,name,password,createtime) VALUES(?,?,?,?);`

// InsertData 新增
func InsertData(id int, name, password string) {
	db.Exec(INSERT_DATA, id, name, password, int32(time.Now().Unix()))
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
