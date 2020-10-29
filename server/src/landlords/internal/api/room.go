package api

import (
	"fmt"
	"landlords/internal/common"
	"landlords/internal/mysql"
	"strconv"
)

type roomData struct {
	level float64
	name  string
}

var rooms = map[string]string{
	"1": "初级房",
	"2": "中级房",
	"3": "高级房",
	"4": "大师房",
}

// CreateRoom 创建房间
func CreateRoom(p interface{}) (d Msgs, err error) {
	d = make(Msgs, 5)
	m := p.(map[string]interface{})
	roomLevel := m["roomLevel"].(string)
	userID := m["userId"].(int)
	level, _ := strconv.Atoi(roomLevel)
	var roots = [2]common.User{
		*NewUser("guest1"),
		*NewUser("guest2"),
	}

	d["id"] = userID
	d["name"] = rooms[roomLevel]
	d["rate"] = level
	d["bottom"] = level * 10
	d["roots"] = roots

	err = mysql.InsertRoom(d)
	fmt.Println(d)
	return d, err
}
