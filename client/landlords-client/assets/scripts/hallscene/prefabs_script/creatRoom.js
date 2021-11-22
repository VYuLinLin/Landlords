import myglobal from "../../mygolbal.js"
import http from "../../util/http.js";
cc.Class({
  extends: cc.Component,
  properties: {
    roomList: [cc.Node],
  },
  onLoad() {
    this._event = cc.ws.on(cc.ws.MESSAGE, this.onMessage.bind(this))
    cc.ws.send([cc.wsApi.REQ_ROOM_LIST])
  },
  onDestroy() {
    cc.ws.off(this._event)
  },
  onMessage(e) {
    if (e[0] === cc.wsApi.RSP_ROOM_LIST) {
      const roomList = this.roomList
      roomList.forEach(a => {
        a.active = false
      });
      const list = e[1]
      for (let i = 0; i < list.length; i++) {
        const {room_id, entrance_fee} = list[i];
        roomList[room_id-1].active = true
        const label = roomList[room_id-1].children[0].getComponent(cc.Label)
        label.string = entrance_fee
      }
    }
    if (e[0] === cc.wsApi.RSP_JOIN_TABLE) {
      const tableId = e[1]
      const tableUserList = e[2]
      console.log(this._event, tableId, tableUserList)
      if (tableId) {
        myglobal.playerData.roomId = tableId
        myglobal.playerData.roots = e[2]
        cc.director.loadScene("gameScene")
        this.node.destroy()
      }
      // const {rate, bottom, id, name, roots} = res
    //   myglobal.playerData.bottom = bottom
    //   myglobal.playerData.rate = rate
    //   myglobal.playerData.roomId = id
    //   myglobal.playerData.roomName = name
    //   myglobal.playerData.roots = roots
    //   cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
    //   cc.director.loadScene("gameScene")
    //   this.node.destroy()
    }
  },
  onBtnClose() {
    this.node.destroy()
  },
  // 进入游戏房间
  onButtonClick(event, room_id) {
    cc.ws.send([cc.wsApi.REQ_JOIN_ROOM, room_id])


    // const data = {
    //   roomLevel: room_id,
    //   userId: myglobal.playerData.userId
    // }
    // myglobal.socket.request_creatroom(data, function (err, res) {
    //   console.log("创建房间", res)
    //   const {rate, bottom, id, name, roots} = res
    //   myglobal.playerData.bottom = bottom
    //   myglobal.playerData.rate = rate
    //   myglobal.playerData.roomId = id
    //   myglobal.playerData.roomName = name
    //   myglobal.playerData.roots = roots
    //   cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
    //   cc.director.loadScene("gameScene")
    //   this.node.destroy()
    // }.bind(this))
  }

});
