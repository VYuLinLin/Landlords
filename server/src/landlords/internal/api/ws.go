package api

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"landlords/internal/db"
	"landlords/internal/game/room"
	"log"
	"net/http"
	"strconv"
)

var (
	upGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if r.Method != "GET" {
				fmt.Println("method is not GET")
				return false
			}
			if r.URL.Path != "/ws/" {
				fmt.Println("path error")
				return false
			}
			return true
		},
	} //不验证origin
)

// StartServeWs 初始化websocket长链接
func StartServeWs(w http.ResponseWriter, r *http.Request) (err error) {
	//defer func() {
	//	fmt.Println("StartServeWs, defer", err.Error())
	//	w.Write([]byte(err.Error()))
	//}()
	params := r.URL.Query()
	userId, err := strconv.Atoi(params.Get("id"))
	logs.Debug("strconv.Atoi(params.Get", userId, err)
	if err != nil {
		log.Panicf("id conversion err:%v\n\n", err)
		return err
	}
	user := room.GetPlayerData(userId)
	if user == nil {
		logs.Debug("未在座位中", err)
		err = db.UpdateUserRoomIdAndTableId(0, 0, &db.User{
			ID: userId,
		})
		if err != nil {
			logs.Debug("退出失败", err)
		}
		return err
	}
	conn, err := upGrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Panicf("upgrade err:%v\n", err)
		return err
	}

	user.Conn = conn
	go user.ReadPump()
	return nil
}
