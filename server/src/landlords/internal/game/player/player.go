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

type Request struct {
	Action string      `json:"action"` // ws请求接口
	Data   interface{} `json:"data"`
}

type Response struct {
	Action string      `json:"action"` // ws推送接口
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
}

// ws 配置相关
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

// ws 推送主题
const (
	UserUpdate    = "user/update"     // 更新用户信息
	RoomList      = "room/list"       // 房间列表
	RoomJoin      = "room/join"       // 进入房间
	RoomJoinOther = "room/join/other" // 其他玩家进入房间
	RoomLeave     = "room/leave"      // 离开房间
	TableInfo     = "table/info"      // 桌子信息
	TableStatus   = "table/status"    // 桌子状态

	TableUpdate      = "table/data/update" // 桌子状态
	TableJoin        = "table/join"        // 进入桌子
	PlayerReady      = "player/ready"      // 玩家准备
	PlayerDeal       = "player/deal"       // 发牌
	TableCallPoints  = "table/callPoints"  // 抢地主
	PlayerCallPoints = "player/callPoints" // 玩家叫分
	TableHoleCards   = "table/holeCards"   // 显示底牌
)

var (
	oldLine = []byte{'\n'}
	newLine = []byte{' '}
)

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
	var msgByte []byte
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

// ReadPump 心跳、消息接受
func (c *Player) ReadPump() {
	defer func() {
		logs.Debug("readPump exit", *c.Conn)
		err := c.CloseWS()
		if err != nil {
			logs.Error("ws 关闭错误：", err)
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	//c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		if c.Conn == nil {
			return
		}
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
		message = bytes.TrimSpace(bytes.Replace(message, oldLine, newLine, -1))
		req := &Request{}
		err = json.Unmarshal(message, &req)
		if err != nil {
			logs.Error("message unmarshal1 err, user_id[%d] err:%v", c.ID, err)
		} else {
			c.wsRequest(req)
		}
	}
}

// wsRequest 处理websocket请求
func (c *Player) wsRequest(r *Request) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error("wsRequest panic:%v ", err)
			c.SendMsg(r.Action, 500, err)
		}
	}()
	switch r.Action {
	case PlayerReady:
		c.Ready = 1
		data := map[int]int{c.ID: c.Ready}

		if c.Next != nil {
			data[c.Next.ID] = c.Next.Ready
			if c.Next.Next != nil {
				data[c.Next.Next.ID] = c.Next.Next.Ready
			}
		}
		c.AllSendMsg(PlayerReady, data)
	}
}
