package table

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"landlords/internal/db"
	"landlords/internal/game/player"
	"time"
)

type Table struct {
	TableID    int64            `json:"table_id"`
	CreateTime string           `json:"create_time"`
	Status     int              `json:"status"`
	Creator    *player.Player   `json:"creator"`
	Players    []*player.Player `json:"players"`
	//Player1    *player.Player `json:"player_1"`
	//Player2    *player.Player `json:"player_2"`
	//Player3    *player.Player `json:"player_3"`
}

// IsAtTable 是否已经在座位中
func (t *Table) IsAtTable(id int) *player.Player {
	for i := 0; i < len(t.Players); i++ {
		if t.Players[i] != nil && t.Players[i].ID == id {
			return t.Players[i]
		}
	}
	return nil
}

// 重置用户前后顺序
func (t *Table) setPlayerNext() {
	pCount := len(t.Players)
	for i := 0; i < pCount; i++ {
		player := t.Players[i]
		player.Next = t.Players[(i+1)%pCount]
	}
}

// JoinTable 加入空位
func (t *Table) JoinTable(u *db.User) (err error) {
	if len(t.Players) < 3 {
		return errors.New("桌子人员已满")
	}

	t.Players = append(t.Players, &player.Player{User: u})

	t.setPlayerNext()

	return err
}

// LeaveTable 离开桌子
func (t *Table) LeaveTable(u *db.User) (ok bool) {
	var newPlayers []*player.Player
	for i := 0; i < len(t.Players); i++ {
		player := t.Players[i]
		if u.ID == player.ID {
			// 关闭ws
			err := player.Conn.Close()
			logs.Info("LeaveTable Player[%s]:", player.ID, err)

			// 删除房主
			if u.ID == t.Creator.ID {
				t.Creator = nil
			}

			ok = true
			break
		}
		newPlayers = append(newPlayers, player)
	}
	t.Players = newPlayers
	t.setPlayerNext()

	// 转移房主
	if t.Creator == nil && len(newPlayers) > 0 {
		t.Creator = newPlayers[0]
	}

	return ok
}

func JoinNewTable(u *db.User) (t *Table, err error) {
	now := time.Now()
	creator := &player.Player{User: u}
	return &Table{
		TableID:    now.Unix(),
		CreateTime: now.Format("2006-01-02 15:04:05"),
		Creator:    creator,
		Players:    []*player.Player{creator},
	}, err
}
