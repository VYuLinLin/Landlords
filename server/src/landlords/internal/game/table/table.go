package table

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"landlords/internal/common"
	"landlords/internal/db"
	"landlords/internal/game/logic/pokerlogic"
	"landlords/internal/game/player"
	"landlords/internal/game/poker"
	"math/rand"
	"time"
)

const (
	GameWaitting       = iota // 准备
	GamePushCard              // 发牌
	GameCalling               // 叫地主
	GameSnatch                // 抢地主
	GameShowBottomCard        // 显示底牌
	GamePlaying               // 出牌
	GameEnd                   // 结束
)

type Table struct {
	TableID     int64            `json:"table_id"`
	CreateTime  string           `json:"create_time"`
	Status      int              `json:"status"`
	Creator     *player.Player   `json:"creator"`
	Players     []*player.Player `json:"players"`
	holePokers  poker.Pokers     `json:"-"`
	activeIndex int              `json:"-"`
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("StartGame panic:", err)
		}
	}()

	t.Status = 0
	t.NextStatus()

	return err
}

// NextStatus 更新桌子状态到下一步
func (t *Table) NextStatus(arg ...int) (ok bool) {
	len := len(arg)
	if len > 0 {
		status := arg[0]
		if t.Status > status {
			return ok
		}
		t.Status = status
	} else {
		t.Status += 1
	}

	switch t.Status {
	case GamePushCard:
		// 发牌
		t.gamePushCard()
	case GameCalling:
		// 叫地主
		t.activeIndex = rand.Intn(3)
		go t.gameCalling(t.activeIndex, 1)
	case GameSnatch:
		// 抢地主
		go t.gameSnatch(t.activeIndex, 1)
	}

	ok = true
	return ok
}

// 发牌
func (t *Table) gamePushCard() (err error) {
	if t.Status > GamePushCard {
		return err
	}
	pokers := &pokerlogic.Card{}
	pokers.GetNewPokers()
	t.holePokers = pokers.HoleCards
	for i := 0; i < len(t.Players); i++ {
		p := t.Players[i]
		cards := pokers.Cards[i]
		data := map[string]poker.Pokers{"pokers": cards}
		err = p.SendMsg(common.PlayerDeal, 200, data)
		logs.Error("[%v]PlayerDeal Error:", p.ID, err)
		err = p.SetCards(cards)
		logs.Error("[%v]SetCards Error:", p.ID, err)
	}
	t.NextStatus()
	return err
}

// 叫地主
func (t *Table) gameCalling(index, count int) (err error) {
	if t.Status > GameCalling {
		return err
	}
	if count > 3 {
		t.NextStatus(GameSnatch)
		return err
	}

	callPlayer := t.Players[index]
	data := map[string]int{"id": callPlayer.ID}
	for i := 0; i < len(t.Players); i++ {
		p := t.Players[i]
		err = p.SendMsg(common.TableCalling, 200, data)
		logs.Error("[%v]TableCalling Error:", p.ID, err)
	}

	time.Sleep(10 * time.Second)

	nextActive := (index + 1) % 3
	addCount := count + 1
	if t.Players[nextActive].GameStatus < 2 {
		t.Players[index].GameStatus = 3 // 不叫
		go t.gameCalling(nextActive, addCount)
	}
	return err
}

// 抢地主
func (t *Table) gameSnatch(index, count int) (err error) {
	if t.Status > GameSnatch {
		return err
	}
	if count > 3 {
		t.NextStatus(GameShowBottomCard)
		return err
	}

	callPlayer := t.Players[index]
	data := map[string]int{"id": callPlayer.ID}
	for i := 0; i < len(t.Players); i++ {
		p := t.Players[i]
		err = p.SendMsg(common.TableSnatch, 200, data)
		logs.Error("[%v]TableSnatch Error:", p.ID, err)
	}

	time.Sleep(10 * time.Second)

	nextActive := (index + 1) % 3
	addCount := count + 1
	if t.Players[nextActive].GameStatus < 4 {
		t.Players[index].GameStatus = 5 // 不抢
		go t.gameSnatch(nextActive, addCount)
	}
	return err
}

// 显示底牌
func (t *Table) gameShowBottomCard(index, count int) (err error) {
	if t.Status > GameShowBottomCard {
		return err
	}
	data := map[string]poker.Pokers{"pokers": t.holePokers}
	for i := 0; i < len(t.Players); i++ {
		p := t.Players[i]
		err = p.SendMsg(common.TableShowHolePokers, 200, data)
		logs.Error("[%v]PlayerSnatch Error:", p.ID, err)
	}

	return err
}
