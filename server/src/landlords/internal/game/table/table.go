package table

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"landlords/internal/common"
	"landlords/internal/db"
	"landlords/internal/game/logic/pokerlogic"
	"landlords/internal/game/player"
	"landlords/internal/game/poker"
	"time"
)

const (
	GameWaitting       = iota // 准备
	GamePushCard              // 发牌
	GameCallScore             // 叫分
	GameShowBottomCard        // 显示底牌
	GamePlaying               // 出牌
	GameEnd                   // 结束
)

type Table struct {
	TableID    int64            `json:"table_id"`
	CreateTime string           `json:"create_time"`
	Status     int              `json:"status"`
	Creator    *player.Player   `json:"creator"`
	Players    []*player.Player `json:"players"`
	pokers     *pokerlogic.Card `json:"-"`
}

func JoinNewTable(u *db.User) (t *Table, err error) {
	now := time.Now()
	creator := &player.Player{User: u}
	return &Table{
		TableID:    now.Unix(),
		CreateTime: now.Format("2006-01-02 15:04:05"),
		Creator:    creator,
		Players:    []*player.Player{creator},
		Status:     GameWaitting,
	}, err
}

// AllSendMsg 推送给座位上的所有玩家(与*player.AllSendMsg方法功能相同)
func (t *Table) AllSendMsg(action string, data interface{}) {
	for i := 0; i < len(t.Players); i++ {
		c := t.Players[i]
		if c != nil {
			c.SendMsg(action, 200, data)
		}
	}
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

// 重置房主、用户前后顺序
func (t *Table) setPlayerNext() {
	defer func() {
		logs.Error("setPlayerNext")
	}()
	pCount := len(t.Players)
	if pCount > 1 {
		// 重置顺序
		for i := 0; i < pCount; i++ {
			player := t.Players[i]
			player.Next = t.Players[(i+1)%pCount]
			player.NextID = t.Players[(i+1)%pCount].ID
		}
	}

	//	重置房主
	if pCount == 0 {
		t.Creator = nil
	} else {
		t.Creator = t.Players[0]
	}
}

// JoinTable 加入空位
func (t *Table) JoinTable(u *db.User) (err error) {
	if len(t.Players) >= 3 {
		return errors.New("桌子人员已满")
	}

	t.Players = append(t.Players, &player.Player{User: u})

	t.setPlayerNext()

	return err
}

// LeaveTable 离开桌子
func (t *Table) LeaveTable(u *db.User) (err error) {
	var newPlayers []*player.Player
	for i := 0; i < len(t.Players); i++ {
		p := t.Players[i]
		if u.ID == p.ID {
			// 关闭ws
			err = p.CloseWS()
			logs.Info("LeaveTable Player[%v]: %v", p.ID, err)
			// 更新mysql
			err = db.UpdateUserRoomIdAndTableId(0, 0, u)
		} else {
			newPlayers = append(newPlayers, p)
		}
	}

	t.Players = newPlayers
	t.setPlayerNext()
	//	推送
	if err == nil && len(t.Players) > 0 {
		data := map[string]int{
			"user_id":    u.ID,
			"creator_id": t.Creator.ID,
		}
		t.AllSendMsg(common.RoomLeave, data)
	}
	return err
}

// StartGame 所有玩家准备完成，开始游戏
// 发牌 =》 叫地主 =》抢地主 =》 出牌 =》 结束
func (t *Table) StartGame() (err error) {
	pokers := &pokerlogic.Card{}
	pokers.GetNewPokers()
	t.pokers = pokers
	t.Status = GamePushCard
	// 推送
	for i := 0; i < len(t.Players); i++ {
		p := t.Players[i]
		data := map[string]poker.Pokers{
			"pokers": pokers.Cards[i],
		}
		p.SendMsg(common.PlayerDeal, 200, data)
	}
	return err
}
