package ws

import (
	"github.com/astaxie/beego/logs"
	"landlords/internal/common"
	"landlords/internal/game/player"
	"landlords/internal/game/room"
)

type Request struct {
	Action string      `json:"action"` // ws请求接口
	Data   interface{} `json:"data"`
}

// wsRequest 处理websocket请求
func WSRequest(r *Request, c *player.Player) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error("wsRequest panic:%v ", err)
			c.SendMsg(r.Action, 500, err)
		}
	}()
	switch r.Action {
	case common.PlayerReady:
		c.Ready = 1
		data := map[int]int{c.ID: c.Ready}

		if c.Next != nil {
			data[c.Next.ID] = c.Next.Ready
			if c.Next.Next != nil {
				data[c.Next.Next.ID] = c.Next.Next.Ready
			}
		}
		c.AllSendMsg(common.PlayerReady, data)
		if len(data) == 3 {
			for _, v := range data {
				if v == 0 {
					return
				}
			}
			r, err := room.GetTableData(c.TABLEID)
			if err == nil {
				go r.Table.StartGame()
			}
		}
	}
}
