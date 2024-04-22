package mysql

import (
	"fmt"
	"landlords/internal/common"
	"strconv"
	"time"
)

var INSERT_ROOM = `INSERT INTO rooms(name,create_time) VALUES(?,?,?);`

// InsertRoom 新增游戏房间
func InsertRoom(d map[string]interface{}) (err error) {
	var keys string
	var vals []interface{}
	for key, val := range d {
		keys = keys + key + ","
		if key == "roots" {
			v := val.([2]common.User)
			rootsID := strconv.Itoa(v[0].UserID) + "," + strconv.Itoa(v[1].UserID)
			vals = append(vals, rootsID)
			continue
		}
		vals = append(vals, val)
	}
	vals = append(vals, int32(time.Now().Unix()))
	var insertVals = ""
	for i := 1; i < len(vals); i++ {
		insertVals += "?,"
	}
	// INSERT_ROOM
	var INSERT_ROOM = "INSERT INTO rooms(" + keys + "create_time) VALUES(" + insertVals + "?);"
	result, err := db.Exec(INSERT_ROOM, vals...)
	if err != nil {
		fmt.Printf("Insert data failed, err:%v \n", err)
		return err
	}
	// 获取插入数据的自增ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Get insert id failed, err:%v \n", err)
		return err
	}
	fmt.Println("Insert data id:", lastInsertID)
	// 通过RowsAffected获取受影响的行数
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed, err:%v \n", err)
		return err
	}
	fmt.Println("Affected rows:", rowsaffected)
	return nil
}
