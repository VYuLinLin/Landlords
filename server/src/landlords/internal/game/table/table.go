package table

import (
	"errors"
	"landlords/internal/mysql"
	"time"
)

type Table struct {
	TableID    int64       `json:"table_id"`
	CreateTime string      `json:"create_time"`
	Player1    *mysql.User `json:"player_1"`
	Player2    *mysql.User `json:"player_2"`
	Player3    *mysql.User `json:"player_3"`
}

// IsAtTable 是否已经在座位中
func (t *Table) IsAtTable(u *mysql.User) bool {
	if t.Player1 != nil && t.Player1.ID == u.ID {
		return true
	} else if t.Player2 != nil && t.Player2.ID == u.ID {
		return true
	} else if t.Player3 != nil && t.Player3.ID == u.ID {
		return true
	}
	return false
}

// JoinTable 加入空位
func (t *Table) JoinTable(u *mysql.User) error {
	err := errors.New("桌子人员已满")
	if t.Player1 == nil {
		t.Player1 = u
		err = nil
	} else if t.Player2 == nil {
		t.Player2 = u
		err = nil
	} else if t.Player3 == nil {
		t.Player3 = u
		err = nil
	}
	return err
}

func JoinNewTable(u *mysql.User) (t *Table, err error) {
	now := time.Now()
	return &Table{
		TableID:    now.Unix(),
		CreateTime: now.Format("2006-01-02 15:04:05"),
		Player1:    u,
	}, err
}
