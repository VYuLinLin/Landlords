import myglobal from "../../mygolbal.js"
cc.Class({
  extends: cc.Component,

  onBtnClose() {
    this.node.destroy()
  },
  // 进入游戏房间
  onButtonClick(event, value) {
    const data = {
      roomLevel: value,
      userId: myglobal.playerData.userId
    }
    myglobal.socket.request_creatroom(data, function (err, res) {
      console.log("创建房间", res)
      const {rate, bottom, id, name} = res
      // const { bottom, rate } = defines.jdRoomConfig['rate_' + value]
      // const roomId = `${rate}_${bottom}_${Math.floor(Math.random() * 1000)}`
      myglobal.playerData.bottom = bottom
      myglobal.playerData.rate = rate
      myglobal.playerData.roomId = id
      myglobal.playerData.roomName = name
      cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
      cc.director.loadScene("gameScene")
      this.node.destroy()
    }.bind(this))
  }

});
