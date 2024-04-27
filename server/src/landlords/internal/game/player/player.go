package player

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"landlords/internal/db"
	p "landlords/internal/poker"
	"strconv"
	"sync"
	"time"
)

// Player 游戏中的玩家信息
type Player struct {
	*db.User
	Conn  *websocket.Conn `json:"-"`
	cards []p.Poker
	Mux   sync.RWMutex `json:"-"`
}

// Players information
type Players struct {
	user1     Player
	user2     Player
	user3     Player
	HoleCards p.Pokers
}

type Request struct {
	Action string      `json:"action"` // ws请求接口
	Data   interface{} `json:"data"`
}

type Response struct {
	Action string      `json:"action"` // ws推送接口
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
}

const (
	writeWait      = 1 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512

	RoleFarmer   = 0
	RoleLandlord = 1

	ReqHeart = "1"
	ResHeart = "2"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Player) SendMsg(action string, code int, data interface{}) {
	res := &Response{
		Action: action,
		Code:   code,
		Data:   data,
	}
	var msgByte []byte
	var err error
	if action == ResHeart {
		heart, _ := strconv.Atoi(ResHeart)
		msgByte, _ = json.Marshal(heart)
	} else {
		msgByte, err = json.Marshal(res)
		if err != nil {
			logs.Error("send msg [%v] marsha1 err:%v", string(msgByte), err)
			return
		}
	}
	//err = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	//if err != nil {
	//	logs.Error("send msg SetWriteDeadline [%v] err:%v", string(msgByte), err)
	//	return
	//}
	c.Mux.Lock()
	w, err := c.Conn.NextWriter(websocket.BinaryMessage)

	if err != nil {
		err = c.Conn.Close()
		if err != nil {
			logs.Error("close client err: %v", err)
		}
	}
	_, err = w.Write(msgByte)
	c.Mux.Unlock()

	if err != nil {
		logs.Error("Write msg [%v] err: %v", string(msgByte), err)
	}
	if err = w.Close(); err != nil {
		err = c.Conn.Close()
		if err != nil {
			logs.Error("close err: %v", err)
		}
	}
}

func (c *Player) ReadPump() {
	defer func() {
		logs.Debug("readPump exit")
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	//c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logs.Error("websocket user_id[%d] unexpected close error: %v", c.ID, err)
			}
			break
		}
		if string(message) == ReqHeart {
			c.SendMsg(ResHeart, 0, nil)
			continue
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		req := &Request{}
		//var data []interface{}
		err = json.Unmarshal(message, &req)
		if err != nil {
			logs.Error("message unmarshal1 err, user_id[%d] err:%v", c.ID, err)
		} else {
			wsRequest(req, c)
		}
	}
}

// wsRequest 处理websocket请求
func wsRequest(r *Request, client *Player) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error("wsRequest panic:%v ", err)
			client.SendMsg(r.Action, 500, err)
		}
	}()
	switch r.Action {
	}
}
