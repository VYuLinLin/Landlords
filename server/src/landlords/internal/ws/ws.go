package ws

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"landlords/internal/api/ws"
	"landlords/internal/common"
	"landlords/internal/db"
	"landlords/internal/game/player"
	"landlords/internal/game/room"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

// ws 配置相关
const (
	writeWait      = 1 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512

	ReqHeart = "1"
	ResHeart = "2"
)

var (
	oldLine = []byte{'\n'}
	newLine = []byte{' '}
)

type Response struct {
	Action string      `json:"action"` // ws推送接口
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
}

// WsServer 创建WsServer结构体
type WsServer struct {
	addr    string
	upgrade *websocket.Upgrader
}

// 初始化WsServer
func NewWsServer() *WsServer {
	log.Println("ws init start")
	ws := new(WsServer)
	ws.addr = ":8089"
	ws.upgrade = &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if r.Method != "GET" {
				logs.Error("method is not GET")
				return false
			}
			if r.URL.Path != "/ws" {
				logs.Error("path error")
				return false
			}
			return true
		},
	}
	return ws
}

// 处理WebSocket连接
func (self *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ws" {
		httpCode := http.StatusInternalServerError
		reasePhrase := http.StatusText(httpCode)
		logs.Error("path error ", reasePhrase)
		http.Error(w, reasePhrase, httpCode)
		return
	}
	params := r.URL.Query()
	userId, err := strconv.Atoi(params.Get("id"))
	if err != nil {
		logs.Error("id conversion err:%v\n\n", err)
		return
	}
	user := room.GetPlayerData(userId)
	if user == nil {
		logs.Debug("未在座位中【%v】", userId, err)
		err = db.UpdateUserRoomIdAndTableId(0, 0, &db.User{
			ID: userId,
		})
		if err != nil {
			logs.Debug("退出失败", err)
		}
		data := &common.ResponseMsg{Code: 1, Msg: "该玩家未在座位中"}
		strParams, _ := json.Marshal(data)
		w.Write(strParams)
		return
	}

	conn, err := self.upgrade.Upgrade(w, r, w.Header())
	if err != nil {
		logs.Error("websocket error:", err)
		return
	}
	if user.Conn != nil {
		err = user.CloseWS()
		if err != nil {
			logs.Error("upgrade err:%v\n", err)
			return
		}
	}
	user.Conn = conn
	logs.Debug("client connect:", conn.RemoteAddr())
	go self.connHandle(conn, user)
}

// 处理WebSocket连接中的消息
func (self *WsServer) connHandle(conn *websocket.Conn, user *player.Player) {
	defer func() {
		err := conn.Close()
		if err != nil {
			logs.Error("ws 关闭错误：", err)
		}
	}()
	conn.SetReadLimit(maxMessageSize)
	stopCh := make(chan any)
	go self.send(conn, stopCh)
	for {
		//conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(5000)))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			close(stopCh)
			var netErr net.Error
			if errors.As(err, &netErr) {
				if netErr.Timeout() {
					fmt.Printf("ReadMessage timeout remote: %v\n", conn.RemoteAddr())
					return
				}
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				fmt.Printf("ReadMessage other remote:%v error: %v \n", conn.RemoteAddr(), err)
			}
			return
		}
		if string(msg) == ReqHeart {
			heart, _ := strconv.Atoi(ResHeart)
			stopCh <- heart
			continue
		}
		fmt.Println("Received message:", string(msg))
		msg = bytes.TrimSpace(bytes.Replace(msg, oldLine, newLine, -1))
		req := &ws.Request{}
		err = json.Unmarshal(msg, &req)
		if err != nil {
			logs.Error("message unmarshal1 err, user_id[%d] err:%v", user.ID, err)
		} else {
			ws.WSRequest(req, user)
		}
	}
}

// 向客户端发送消息
func (self *WsServer) send(conn *websocket.Conn, stopCh chan any) {
	for {
		select {
		case data := <-stopCh:
			msgByte, err := json.Marshal(data)
			if err != nil {
				logs.Error("send msg [%v] marsha1 err:%v", string(msgByte), err)
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, msgByte)
			//fmt.Println("Sending....")
			if err != nil {
				fmt.Println("send message failed ", err)
				return
			}
			//case <-time.After(time.Second * 1):
			//	data := fmt.Sprintf("Hello WebSocket test from server %v", time.Now().UnixNano())
			//	err := conn.WriteMessage(1, []byte(data))
			//	if err != nil {
			//		fmt.Println("send message failed ", err)
			//		return
			//	}
		}
	}
}

// Start 启动WebSocket服务器
/*示例：
ws := NewWsServer()
err := ws.Start()
if err != nil {
	logs.Error("ws init error: ", err)
}*/
func (self *WsServer) Start() (err error) {
	http.HandleFunc("/ws", self.ServeHTTP)
	err = http.ListenAndServe(self.addr, nil)
	return nil
}
