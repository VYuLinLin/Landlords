package player

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"landlords/internal/db"
	p "landlords/internal/game/poker"
	"sync"
)

// Player 游戏中的玩家信息
type Player struct {
	*db.User
	GameStatus int             `json:"game_status"` // 默认0-未准备 1-已准备 2叫地主 3-不叫 4-抢地主 5-不抢 6-出牌 7-不出牌 8-输 9-赢
	Ready      int             `json:"ready"`
	CardCount  int             `json:"card_count"`
	Cards      []p.Poker       `json:"cards"`
	NextID     int             `json:"next_id"`
	Next       *Player         `json:"-"`
	Conn       *websocket.Conn `json:"-"`
	Mux        sync.RWMutex    `json:"-"`
}

// Players information
type Players struct {
	user1     Player
	user2     Player
	user3     Player
	HoleCards p.Pokers
}

type Response struct {
	Action string      `json:"action"` // ws推送接口
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
}

// AllSendMsg 推送给座位上的所有玩家(与*table.AllSendMsg方法功能相同)
func (c *Player) AllSendMsg(action string, data interface{}) {
	id1 := c.ID
	c.SendMsg(action, 200, data)
	if c.Next != nil {
		id2 := c.Next.ID
		if id1 != id2 {
			c.Next.SendMsg(action, 200, data)
		}
		if c.Next.Next != nil {
			id3 := c.Next.Next.ID
			if id1 != id3 && id2 != id3 {
				c.Next.Next.SendMsg(action, 200, data)
			}
		}
	}
}

// SendMsg 推送给座位上的某一位玩家
func (c *Player) SendMsg(action string, code int, data interface{}) (err error) {
	if c.Conn == nil {
		return err
	}
	res := &Response{
		Action: action,
		Code:   code,
		Data:   data,
	}
	msgByte, err := json.Marshal(res)
	if err != nil {
		logs.Error("send msg [%v] marsha1 err:%v", string(msgByte), err)
		return
	}
	//err = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	//if err != nil {
	//	logs.Error("send msg SetWriteDeadline [%v] err:%v", string(msgByte), err)
	//	return
	//}
	c.Mux.Lock()
	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	_, err = w.Write(msgByte)
	c.Mux.Unlock()

	if err != nil {
		logs.Error("Write msg [%v] err: %v", string(msgByte), err)
		return err
	}
	return w.Close()
}

// CloseWS 关闭websocket
func (c *Player) CloseWS() (err error) {
	if c.Conn != nil {
		err = c.Conn.Close()
		c.Conn = nil
	}
	return err
}

// SetCards 设置玩家手牌
func (c *Player) SetCards(cards p.Pokers) (err error) {
	c.Cards = cards
	c.CardCount = len(cards)
	return err
}
