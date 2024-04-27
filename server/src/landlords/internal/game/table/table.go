package table

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"landlords/internal/db"
	"landlords/internal/game/player"
	"time"
)

type Table struct {
	TableID    int64          `json:"table_id"`
	CreateTime string         `json:"create_time"`
	Player1    *player.Player `json:"player_1"`
	Player2    *player.Player `json:"player_2"`
	Player3    *player.Player `json:"player_3"`
}

// IsAtTable 是否已经在座位中
func (t *Table) IsAtTable(id int) *player.Player {
	if t.Player1 != nil && t.Player1.ID == id {
		return t.Player1
	} else if t.Player2 != nil && t.Player2.ID == id {
		return t.Player2
	} else if t.Player3 != nil && t.Player3.ID == id {
		return t.Player3
	}
	return nil
}

// JoinTable 加入空位
func (t *Table) JoinTable(u *db.User) error {
	err := errors.New("桌子人员已满")
	if t.Player1 == nil {
		t.Player1 = &player.Player{User: u}
		err = nil
	} else if t.Player2 == nil {
		t.Player2 = &player.Player{User: u}
		err = nil
	} else if t.Player3 == nil {
		t.Player3 = &player.Player{User: u}
		err = nil
	}
	return err
}

// LeaveTable 离开桌子
func (t *Table) LeaveTable(u *db.User) (ok bool) {
	if t.Player1.ID == u.ID {
		if t.Player1.Conn != nil {
			err := t.Player1.Conn.Close()
			logs.Info("LeaveTable Player1", err)
		}
		t.Player1 = nil
		ok = true
	} else if t.Player2.ID == u.ID {
		if t.Player2.Conn != nil {
			err := t.Player2.Conn.Close()
			logs.Info("LeaveTable Player2", err)
		}
		t.Player2 = nil
		ok = true
	} else if t.Player3.ID == u.ID {
		if t.Player3.Conn != nil {
			err := t.Player3.Conn.Close()
			logs.Info("LeaveTable Player3", err)
		}
		t.Player3 = nil
		ok = true
	}
	return ok
}

func JoinNewTable(u *db.User) (t *Table, err error) {
	now := time.Now()
	return &Table{
		TableID:    now.Unix(),
		CreateTime: now.Format("2006-01-02 15:04:05"),
		Player1:    &player.Player{User: u},
	}, err
}
