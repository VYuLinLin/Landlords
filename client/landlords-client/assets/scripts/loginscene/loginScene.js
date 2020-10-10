import myglobal from "../mygolbal.js"
cc.Class({
  extends: cc.Component,

  properties: {
    wait_node: cc.Node,
    protocaolNode: cc.Node
  },

  // LIFE-CYCLE CALLBACKS:

  onLoad() {
    cc.director.preloadScene("hallScene")
  },

  start() {},
  // 确认协议
  confirmProtocol() {
    this.protocaolNode.children[0].active = !this.protocaolNode.children[0].active
  },
  onButtonCilck(event, customData) {
    if (!this.protocaolNode.children[0].active) {
      const anim = this.protocaolNode.getComponent(cc.Animation)
      anim.play('scale')
      return
    }
    switch (customData) {
      case "wx_login":
        console.log("wx_login request")

        //this.wait_node.active = true

        myglobal.socket.request_wxLogin({
          uniqueID: myglobal.playerData.uniqueID,
          // userId: myglobal.playerData.userId,
          userName: myglobal.playerData.userName,
          avatarUrl: myglobal.playerData.avatarUrl,
        }, function (err, result) {
          //请求返回
          //先隐藏等待UI
          //this.wait_node.active = false
          if (err != 0) {
            console.log("err:" + err)
            return
          }

          console.log("login sucess" + JSON.stringify(result))
          myglobal.playerData.gobal_count = result.goldcount
          cc.director.loadScene("hallScene")
        }.bind(this))
        break
      case 'guest_login':
        this.wait_node.active = true
        myglobal.socket.request_login({isGuest: 1}, function(status, res) {
          if (status) {
            console.log("err:" + status+res)
            return
          }
          console.log(res)
          const {userId, userName} = res
          myglobal.playerData.userId = userId
          myglobal.playerData.userName = `${userName}_${userId}`
          cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
          cc.director.loadScene("hallScene")
        }.bind(this))
        // const count = Math.floor(Math.random() * 100000)
        // const userName = `guest_${count}`
        // myglobal.playerData.userId = `${count}`
        // myglobal.playerData.userName = userName
        // cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
        // cc.director.loadScene("hallScene")
      default:
        break
    }
  }
});
