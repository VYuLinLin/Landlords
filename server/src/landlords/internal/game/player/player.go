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
	Conn   *websocket.Conn `json:"-"`
	Ready  int             `json:"ready"`
	Cards  []p.Poker       `json:"cards"`
	Next   *Player         `json:"-"`
	NextID int             `json:"next_id"`
	Mux    sync.RWMutex    `json:"-"`
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
