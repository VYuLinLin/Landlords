package api

import (
	"errors"
	"fmt"
	"landlords/internal/db"
	"landlords/internal/game/room"
)

type RoomApi struct{}

type roomData struct {
	RoomName  string `json:"room_name"`
	RoomLevel int    `json:"room_level"`
	TableID   int64  `json:"table_id"`
}

// JoinRoom 进入房间
func (r *RoomApi) JoinRoom(p interface{}) (data *roomData, err error) {
	m := p.(map[string]interface{})
	if m["roomLevel"] == nil || m["userId"] == nil {
		return nil, errors.New("用户ID或房间等级不能为空")
	}
	level := m["roomLevel"].(string)
	userID := m["userId"].(float64)
	User := &db.User{
		ID: int(userID),
	}
	user, err := db.QueryUserId(User)
	fmt.Println("RoomApi JoinRoom:", *user, *User)
	if err != nil {
		return nil, err
	}
	roomInfo, err := room.JoinRoom(user, level)
	data = &roomData{
		RoomName:  roomInfo.RoomName,
		RoomLevel: roomInfo.RoomLevel,
		TableID:   roomInfo.Table.TableID,
	}
	return data, err
}

// ExitRoom 退出房间
func (r *RoomApi) ExitRoom(p interface{}) (err error) {
	m := p.(map[string]interface{})
	if m["userId"] == nil || m["userId"] == "" {
		return errors.New("userId不能为空")
	}
	userID := m["userId"].(float64)
	User := &db.User{
		ID: int(userID),
	}
	user, err := db.QueryUserId(User)
	fmt.Println("RoomApi ExitRoom:", *user, err)
	if err != nil {
		return err
	}
	err = room.ExitRoom(user)
	return err
}

// GetTable 根据桌子id获取桌子信息
func (r *RoomApi) GetTable(p interface{}) (d *room.Info, err error) {
	m := p.(map[string]interface{})
	if m["roomLevel"] == nil || m["tableId"] == nil || m["tableId"] == "" {
		return nil, errors.New("tableId不能为空")
	}
	id := m["tableId"].(int64)
	d, err = room.GetTableData(id)
	return d, err
}
