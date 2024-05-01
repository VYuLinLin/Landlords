/*
 * @Author: X1-EXT\lylin lylin888@163.com
 * @Date: 2022-05-09 13:33:49
 * @LastEditors: X1-EXT\lylin lylin888@163.com
 * @LastEditTime: 2024-04-22 17:28:11
 * @FilePath: \landlords-client\assets\scripts\hallscene\prefabs_script\creatRoom.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import myglobal from "mygolbal.js"
import http from "http.js";

cc.Class({
  extends: cc.Component,
  properties: {
    roomList: [cc.Node],
    tip: cc.Prefab,
  },
  onLoad() {
    // this._event = cc.ws.on(cc.ws.MESSAGE, this.onMessage.bind(this))
    // cc.ws.send({ action: cc.wsApi.roomList })
  },
  onDestroy() {
    // cc.ws.off(this._event)
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
    // if (e.action === cc.wsApi.roomJoinSelf) {
    //   if (e.code === 200) {
    //     const {room_id, table_id} = e.data
    //     myglobal.playerData.roomId = room_id
    //     myglobal.playerData.tableId = table_id
    //     cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
    //     cc.director.loadScene("gameScene")
    //     this.node.destroy()
    //   }
    // }
  },
  onBtnClose() {
    this.node.destroy()
  },
  // 进入游戏房间
  onJoinRoomHandler(event, level) {
    const data = {
      userId: myglobal.playerData.userId,
      roomLevel: level,
    }
    http.post(http.joinRoom, data, (res) => {
        console.log(res)
        if (res.code) {
            this.tipNode && this.tipNode.destroy();
            this.tipNode = cc.instantiate(this.tip);
            this.tipNode.getComponent(cc.Label).string = res.msg;
            this.tipNode.parent = this.node;
            return;
        }
        myglobal.playerData.roomId = res.data.room_level
        myglobal.playerData.tableId = res.data.table_id
        cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
        cc.director.loadScene('gameScene')
        this.onBtnClose()
    });
    // const data = {
    //   action: cc.wsApi.roomJoinSelf,
    //   data: level
    // }
    // cc.ws.send(data)
  },

});
