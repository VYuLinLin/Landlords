import myglobal from "../../mygolbal.js"
cc.Class({
  extends: cc.Component,
  properties: {
    roomList: [cc.Node],
    tip: cc.Prefab,
  },
  onLoad() {
    this._event = cc.ws.on(cc.ws.MESSAGE, this.onMessage.bind(this))
    cc.ws.send({ action: cc.wsApi.roomList })
  },
  onDestroy() {
    cc.ws.off(this._event)
  },
  onMessage(e) {
    if (!e) return
    if (e.code === 500) {
      this.tipNode && this.tipNode.destroy()
      this.tipNode = cc.instantiate(this.tip)
      this.tipNode.getComponent(cc.Label).string = e.data
      this.tipNode.parent = this.node
      return
    }
    if (e.action === cc.wsApi.roomList) {
      const roomList = this.roomList
      roomList.forEach(a => {
        a.active = false
      });
      const list = e.data
      for (let i = 0; i < list.length; i++) {
        const {room_id, entrance_fee} = list[i];
        roomList[room_id-1].active = true
        const label = roomList[room_id-1].children[0].getComponent(cc.Label)
        label.string = entrance_fee
      }
    }
    if (e.action === cc.wsApi.roomJoinSelf) {
      if (e.code === 200) {
        const {room_id, table_id} = e.data
        myglobal.playerData.roomId = room_id
        myglobal.playerData.tableId = table_id
        cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
        cc.director.loadScene("gameScene")
        this.node.destroy()
      }
    }
  },
  onBtnClose() {
    this.node.destroy()
  },
  // 进入游戏房间
  onButtonClick(event, room_id) {
    const data = {
      action: cc.wsApi.roomJoinSelf,
      data: room_id
    }
    cc.ws.send(data)
  },

});
