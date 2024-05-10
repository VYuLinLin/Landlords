package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

const createTableUsers = `CREATE TABLE IF NOT EXISTS users(
	id INT(10) AUTO_INCREMENT NOT NULL,
	name VARCHAR(64) NOT NULL,
	password VARCHAR(64),
	coin INT(64) DEFAULT 0,
	score INT(5) DEFAULT 0,
	status INT(4) DEFAULT 0,
	create_time VARCHAR(64) NOT NULL,
    room_id INT DEFAULT 0,
    table_id INT(64) DEFAULT 0,
	PRIMARY KEY (id, name, create_time)
); `

var insertUserSql = `INSERT INTO users(coin,name,password,create_time) VALUES(?,?,?,?);`

type User struct {
	ID         int    `json:"id"`
	NAME       string `json:"name"`
	PASSWORD   string `json:"-"`
	COIN       int8   `json:"coin"`
	SCORE      int8   `json:"score"`
	STATUS     int    `json:"status"`
	CREATETIME string `json:"create_time"`
	ROOMID     int    `json:"room_id"`
	TABLEID    int64  `json:"table_id"`
}

// CreateTableUsers 创建用户表
func CreateTableUsers() {
	_, err := db.Exec(createTableUsers)
	if err == nil {
		fmt.Println("create table users success")
	} else {
		// 关闭数据库
		defer db.Close()
		errorHandler(err)
	}
}

// InsertUser 新增用户
func InsertUser(name, password string) {
	result, err := db.Exec(insertUserSql, 0, name, password, time.Now().Format("2006-01-02 15:04:05"))
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

// QueryUser 查询用户
func QueryUser(p *User) (user *User, err error) {
	var row *sql.Row
	if p.ID > 0 {
		row = db.QueryRow("SELECT * FROM users WHERE name=? AND id=?", p.NAME, p.ID)
	} else {
		row = db.QueryRow("SELECT * FROM users WHERE name=?", p.NAME)
	}
	fmt.Println("QueryUser:", *p)
	fmt.Println("QueryRow:", *row)
	user = new(User) // 用new()函数初始化一个结构体对象
	// row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	err = row.Scan(&user.ID, &user.NAME, &user.PASSWORD, &user.COIN, &user.SCORE, &user.STATUS, &user.CREATETIME, &user.ROOMID, &user.TABLEID)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("未查询到此用户")
	}
	return user, err
}

// QueryUserId 根据ID查询用户
func QueryUserId(p *User) (user *User, err error) {
	var row *sql.Row
	if p.ID > 0 {
		row = db.QueryRow("SELECT * FROM users WHERE id=?", p.ID)
	} else {
		return user, errors.New("用户ID不能为空")
	}
	fmt.Println("QueryUser:", *p)
	fmt.Println("QueryRow:", *row)
	user = new(User) // 用new()函数初始化一个结构体对象
	// row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	err = row.Scan(&user.ID, &user.NAME, &user.PASSWORD, &user.COIN, &user.SCORE, &user.STATUS, &user.CREATETIME, &user.ROOMID, &user.TABLEID)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("未查询到此用户")
	}
	return user, err
}

// UpdateUserRoomIdAndTableId 更新用户房间和桌子id
func UpdateUserRoomIdAndTableId(roomID int, tableId int64, p *User) (err error) {
	if p.ID <= 0 {
		return errors.New("用户ID不能为空")
	}
	_, err = db.Exec(`UPDATE users SET room_id=?, table_id=? WHERE id=?;`, roomID, tableId, p.ID)
	logs.Info("UpdateUserRoomIdAndTableId", err)
	return err
}
