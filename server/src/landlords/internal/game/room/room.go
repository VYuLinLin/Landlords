package room

import (
	"errors"
	"fmt"
	"landlords/internal/common"
	"landlords/internal/db"
	"landlords/internal/game/player"
	"landlords/internal/game/table"
	"strconv"
)

type Room struct {
	RoomName  string         `json:"room_name"`
	RoomLevel int            `json:"room_level"`
	Tables    []*table.Table `json:"tables"`
}

type RoomMap map[string]*Room

type Info struct {
	RoomName  string       `json:"room_name"`
	RoomLevel int          `json:"room_level"`
	Table     *table.Table `json:"table"`
}

var Rooms *RoomMap

func init() {
	Rooms = &RoomMap{
		"1": &Room{
			RoomLevel: 1,
			RoomName:  "初级房",
		},
		"2": &Room{
			RoomLevel: 2,
			RoomName:  "中级房",
		},
		"3": &Room{
			RoomLevel: 3,
			RoomName:  "高级房",
		},
		"4": &Room{
			RoomLevel: 4,
			RoomName:  "大师房",
		},
	}
}

// JoinRoom 进入房间
func JoinRoom(u *db.User, level string) (room *Info, err error) {
	r := *Rooms
	room = &Info{
		RoomName:  r[level].RoomName,
		RoomLevel: r[level].RoomLevel,
	}
	tables := r[level].Tables
	fmt.Println("JoinRoom", len(tables), cap(tables), tables, room.Table == nil, room.Table)
	if room.Table == nil {
		// 是否已经在座位中
		for i, l := 0, len(tables); i < l; i++ {
			t := tables[i]
			if t.IsAtTable(u.ID) != nil {
				room.Table = t
				break
			}
		}
	}
	if room.Table == nil {
		// 寻找空位的桌子
		for i, l := 0, len(tables); i < l; i++ {
			t := tables[i]
			err := t.JoinTable(u)
			if err == nil {
				room.Table = t
				break
			}
		}
	}
	if room.Table == nil {
		// 无空位，添加新桌子
		t, _ := table.JoinNewTable(u)
		r[level].Tables = append(tables, t)
		Rooms = &r
		fmt.Println("JoinRoom", *t)
		room.Table = t
	}
	if room.Table != nil {
		roomId := 0
		tableId := room.Table.TableID
		if roomId, err = strconv.Atoi(level); err == nil {
			// 更新mysql
			err = db.UpdateUserRoomIdAndTableId(roomId, tableId, u)
			// ws推送
			if p := GetPlayerData(u.ID); err == nil && p != nil {
				p.ROOMID = roomId
				p.TABLEID = tableId
				data := map[string]any{
					"id":       p.ID,
					"name":     p.NAME,
					"next_id":  p.NextID,
					"table_id": tableId,
				}
				room.Table.AllSendMsg(common.RoomJoin, data)
			}
		} else {
			err = room.Table.LeaveTable(u)
		}
		fmt.Println("UpdateUserTableID", err)
	}
	return room, err
}

// ExitRoom 离开房间
func ExitRoom(u *db.User) (err error) {
	if u.ROOMID <= 0 {
		return err
	}
	r := *Rooms
	tables := r[strconv.Itoa(u.ROOMID)].Tables
	for i, l := 0, len(tables); i < l; i++ {
		t := tables[i]
		if t.IsAtTable(u.ID) != nil {
			err = t.LeaveTable(u)
			break
		}
	}
	return err
}

// GetTableData 根据桌子id获取桌面信息
func GetTableData(id int64) (room *Info, err error) {
	room = &Info{}
	r := *Rooms
Exit:
	for _, val := range r {
		for i, l := 0, len(val.Tables); i < l; i++ {
			t := val.Tables[i]
			if t.TableID == id {
				room.RoomName = val.RoomName
				room.RoomLevel = val.RoomLevel
				room.Table = t
				break Exit
			}
		}
	}
	// 匹配桌子id

	if room.Table == nil {
		err = errors.New("桌子id匹配为空")
	}
	return room, err
}

// GetPlayerData 根据用户id查找当前游戏中的用户
func GetPlayerData(id int) *player.Player {
	r := *Rooms
	for _, val := range r {
		for i, l := 0, len(val.Tables); i < l; i++ {
			t := val.Tables[i]
			u := t.IsAtTable(id)
			if u != nil {
				return u
			}
		}
	}
	return nil
}
