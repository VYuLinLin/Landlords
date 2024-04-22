package mysql

import "fmt"

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

// CreateTableUsers 创建用户表
func CreateTableUsers() {
	_, err := db.Exec(createTableUsers)
	if err == nil {
		fmt.Println("create table users successd")
		// 关闭数据库
		// defer db.Close()
	} else {
		errorHandler(err)
	}
}

const createTableRooms = `CREATE TABLE IF NOT EXISTS rooms(
	id INT(10) AUTO_INCREMENT NOT NULL,
	name VARCHAR(64),
	rate INT(4) DEFAULT 0,
	bottom INT(4) DEFAULT 0,
	status INT(4) DEFAULT 0,
	roots VARCHAR(64),
	create_time VARCHAR(64) NOT NULL,
	PRIMARY KEY (id, create_time)
); `

// CreateTableRooms 创建房间表
func CreateTableRooms() {
	_, err := db.Exec(createTableRooms)
	if err == nil {
		fmt.Println("create table rooms successd")
		// 关闭数据库
		// defer db.Close()
	} else {
		errorHandler(err)
	}
}
