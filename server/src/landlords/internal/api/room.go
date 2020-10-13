package api

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
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
func CreateRoom(p interface{}) (d Msgs) {
	d = make(Msgs, 3)
	m := p.(map[string]interface{})
	roomLevel := m["roomLevel"].(string)
	// userId := m["userId"].(int)
	level, _ := strconv.Atoi(roomLevel)
	rand.Seed(time.Now().Unix())
	var roots = new([2]Msgs)
	roots[0] = NewUser("guest")
	roots[1] = NewUser("guest")

	d["id"] = rand.Intn(100000)
	d["name"] = rooms[roomLevel]
	d["rate"] = level
	d["bottom"] = level * 10
	d["roots"] = roots

	fmt.Println(d)
	return d
}
