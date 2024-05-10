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
	if m["userId"] == nil || m["userId"] == "" {
		return d, errors.New("userId不能为空")
	}
	if m["tableId"] == nil || m["tableId"] == "" {
		return d, errors.New("tableId不能为空")
	}

	userID := m["userId"].(float64)
	tableId := m["tableId"].(float64)
	roomInfo, err := room.GetDeepTableData(int64(tableId))
	if err != nil {
		return d, err
	}
	// 玩家只能看到自己的手牌
	players := roomInfo.Table.Players
	for i := 0; i < len(players); i++ {
		if players[i].ID != int(userID) {
			players[i].Cards = nil
		}
	}
	return roomInfo, err
}
