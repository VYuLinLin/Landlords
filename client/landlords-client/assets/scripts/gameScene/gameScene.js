import myglobal from "mygolbal.js"
import http from "http.js";
import ws from "websocket.js"

const ddzConsts= require('ddzConstants')
const ddzData = require('ddzData')

cc.Class({
  extends: cc.Component,
  properties: {
    bjMusic: {
      type: cc.AudioClip, // 背景音乐
      default: null,     // object's default value is null
    },
    di_label: cc.Label,
    beishu_label: cc.Label,
    roomid_label: cc.Label,
    player_node_prefabs: cc.Prefab,
    btn_ready: cc.Node, // 准备按钮
    // 绑定玩家座位,下面有3个子节点
    players_seat_pos: cc.Node,
    gameUiNode: cc.Node,
    tip: cc.Prefab,
  },
  onLoad() {
    this.roomid_label.string = ddzConsts.roomNames[myglobal.playerData.roomId-1] + '：' + myglobal.playerData.tableId
    
    ws.initWS()

    this.playerNodeList = []
    if (cc.ws) {
      this._event = cc.ws.on(cc.ws.MESSAGE, this.onMessage.bind(this))
    }

    this.getTable()
    // const data = {
    //   action: cc.wsApi.tableInfo,
    //   data: {
    //     room_id: myglobal.playerData.roomId,
    //     table_id: myglobal.playerData.tableId,
    //   }
    // }
    // cc.ws.send(data)

    cc.audioEngine.stopAll()
    cc.audioEngine.play(this.bjMusic, true)

    // ddzData.gameState = ddzConsts.gameStatus.WAITREADY
    // if (!CC_EDITOR) {
    //   ddzData.gameStateNotify.addListener(this.gameStateHandler, this)
    // }
    // this.playerNodeList = []
    // const { roomId, roomName, rate, bottom, clients } = myglobal.playerData

    // this.roomid_label.string = roomName
    // this.beishu_label.string = "倍数：" + rate
    // this.di_label.string = "底：" + bottom
    // console.log('重新开始', ddzData.gameState)
    // this.btn_ready.active = ddzData.gameState < ddzConsts.gameStatus.GAMESTART // 准备按钮
    // // 玩家
    // this.addPlayerNode(myglobal.playerData)
    // // 机器人1
    // this.addPlayerNode({
    //   ...clients[0],
    //   seatindex: 1
    // })
    // // 机器人2
    // this.addPlayerNode({
    //   ...clients[1],
    //   seatindex: 2
    // })
    //监听，给其他玩家发牌(内部事件)
    // this.node.on("pushcard_other_event", function () {
    //   console.log('其他玩家发牌')
    //   for (let i = 0; i < this.playerNodeList.length; i++) {
    //     const node = this.playerNodeList[i]
    //     if (node) {
    //       //给playernode节点发送事件
    //       node.emit("push_card_event")
    //     }
    //   }
    // }.bind(this))

    //监听房间状态改变事件
    // myglobal.socket.onRoomChangeState(function (data) {
    //   //回调的函数参数是进入房间用户消息
    //   console.log("onRoomChangeState:" + data)
    //   this.roomstate = data
    // }.bind(this))
    // 抢地主
    // this.node.on("canrob_event", function (event) {
    //   console.log("gamescene canrob_event:" + event)
    //   //通知给playernode子节点
    //   for (var i = 0; i < this.playerNodeList.length; i++) {
    //     var node = this.playerNodeList[i]
    //     if (node) {
    //       //给playernode节点发送事件
    //       node.emit("playernode_canrob_event", event)
    //     }
    //   }
    // }.bind(this))

    // this.node.on("choose_card_event", function (event) {
    //   this.gameUiNode.emit("choose_card_event", event)
    // }.bind(this))

    // this.node.on("unchoose_card_event", function (event) {
    //   this.gameUiNode.emit("unchoose_card_event", event)
    // }.bind(this))
    //监听给玩家添加三张底牌
    // this.node.on("add_three_card",function(event){
    //     console.log("add_three_card:"+event)
    //     for(var i=0;i<this.playerNodeList.length;i++){
    //         var node = this.playerNodeList[i]
    //         if(node){
    //             //给playernode节点发送事件
    //             node.emit("playernode_add_three_card",event)
    //         }
    //     }
    // }.bind(this))
    return

    myglobal.socket.request_enter_room({}, function (err, result) {
      console.log("enter_room_resp" + JSON.stringify(result))
      if (err != 0) {
        console.log("enter_room_resp err:" + err)
      } else {

        //enter_room成功
        //notify ={"seatid":1,"playerdata":[{"accountid":"2117836","userName":"tiny543","avatarUrl":"http://xxx","coin":1000}]}
        var seatid = result.seatindex //自己在房间里的seatid
        this.playerdata_list_pos = []  //3个用户创建一个空用户列表
        this.setPlayerSeatPos(seatid)

        var playerdata_list = result.playerdata
        var roomId = result.roomId
        this.roomid_label.string = "房间号:" + roomId
        myglobal.playerData.housemanageid = result.housemanageid

        for (var i = 0; i < playerdata_list.length; i++) {
          //consol.log("this----"+this)
          this.addPlayerNode(playerdata_list[i])
        }


      }
      var gamebefore_node = this.node.getChildByName("gamebeforeUI")
      gamebefore_node.emit("init")
    }.bind(this))

    //在进入房间后，注册其他玩家进入房间的事件
    myglobal.socket.onPlayerJoinRoom(function (join_playerdata) {
      //回调的函数参数是进入房间用户消息
      console.log("onPlayerJoinRoom:" + JSON.stringify(join_playerdata))
      this.addPlayerNode(join_playerdata)
    }.bind(this))

    //回调参数是发送准备消息的accountid
    myglobal.socket.onPlayerReady(function (data) {
      console.log("-------onPlayerReady:" + data)
      for (var i = 0; i < this.playerNodeList.length; i++) {
        var node = this.playerNodeList[i]
        if (node) {
          node.emit("player_ready_notify", data)
        }
      }
    }.bind(this))

    myglobal.socket.onGameStart(function () {
      for (var i = 0; i < this.playerNodeList.length; i++) {
        var node = this.playerNodeList[i]
        if (node) {
          node.emit("gamestart_event")
        }
      }

      //隐藏gamebeforeUI节点
      var gamebeforeUI = this.node.getChildByName("gamebeforeUI")
      if (gamebeforeUI) {
        gamebeforeUI.active = false
      }
    }.bind(this))

    //监听服务器玩家抢地主消息
    // myglobal.socket.onRobState(function (event) {
    //   console.log("-----onRobState" + JSON.stringify(event))
    //   //onRobState{"accountid":"2162866","state":1}
    //   for (var i = 0; i < this.playerNodeList.length; i++) {
    //     var node = this.playerNodeList[i]
    //     if (node) {
    //       //给playernode节点发送事件
    //       node.emit("playernode_rob_state_event", event)
    //     }
    //   }
    // }.bind(this))

    //注册监听服务器确定地主消息
    myglobal.socket.onChangeMaster(function (event) {
      console.log("onChangeMaster" + event)
      //保存一下地主id
      myglobal.playerData.masterUserId = event
      for (var i = 0; i < this.playerNodeList.length; i++) {
        var node = this.playerNodeList[i]
        if (node) {
          //给playernode节点发送事件
          node.emit("playernode_changemaster_event", event)
        }
      }
    }.bind(this))

    //注册监听服务器显示底牌消息
    // myglobal.socket.onShowBottomCard(function (event) {
    //   console.log("onShowBottomCard---------" + event)
    //   this.gameUiNode.emit("show_bottom_card_event", event)
    // }.bind(this))
  },
  start() {
    // $socket.on('change_master_notify', this.masterNotify, this)
  },
  onDestroy() {
    if (cc.ws) {
      cc.ws.off(this._event)
      cc.ws.close()
      cc.ws = null
    }
    if (!CC_EDITOR) {
      // ddzData.gameStateNotify.removeListener(this.gameStateHandler, this)
      cc.audioEngine.stopAll()
    }
    // $socket.remove('change_master_notify', this)
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
    switch (e.action) {
      case cc.wsApi.tableInfo:
        const {state, room_id, table_id, creator, clients} = e.data

        myglobal.playerData.creator = creator
        ddzData.gameState = state
        this.roomid_label.string = ddzConsts.roomNames[room_id-1] + '：' + table_id
        this.btn_ready.active = state < ddzConsts.gameStatus.GAMESTART // 准备按钮
        for (var i = 0; i < this.playerNodeList.length; i++) {
          var node = this.playerNodeList[i]
          node.destroy()
        }
        this.playerNodeList = []
        clients.map(c => {
          this.addPlayerNode(c)
        })
        break
      case cc.wsApi.tableUpdate:
        myglobal.playerData.creator = e.data.creator
        ddzData.gameState = e.data.state
        break
      case cc.wsApi.roomJoin:
        // 用户进入后需要更新座位顺序，所以这里重新请求获取数据
        this.getTable()
        break
      case cc.wsApi.roomLeave:
        const {user_id, creator_id} = e.data
        myglobal.playerData.creator = creator_id
        if (myglobal.playerData.userId === user_id) {
          myglobal.playerData.roomId = ''
          cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
          cc.director.loadScene("hallScene")
        } else {
          const [player] = this.playerNodeList.filter(a => a.getComponent("player_node").userId === user_id)
          if (!player) return
          player.destroy()
          this.playerNodeList = this.playerNodeList.filter(a => {
            const node = a.getComponent("player_node")
            if (node.userId !== user_id) {
              node.updateCreator()
              return true
            }
            return false
          })
        }
        break
      case cc.wsApi.playerReady:
        for (let i = 0; i < this.playerNodeList.length; i++) {
          const playerNode = this.playerNodeList[i].getComponent("player_node")
          const {seat_index, userId} = playerNode
          if (seat_index === 0) {
            this.btn_ready.active = !e.data[userId]
          } else {
            playerNode.updateReadyStatus(e.data[userId])
          }
        }
        break
      case cc.wsApi.playerDeal:
        for (let i = 0; i < this.playerNodeList.length; i++) {
          const node = this.playerNodeList[i]
          if (node.seat_index === 0) {
            window.$socket.emit('pushcard_notify', e.data.pokers)
          } else {
            //给playernode节点发送事件
            node.emit("push_card_event")
          }
        }
        break
      case cc.wsApi.tableCalling:
        window.$socket.emit('canrob_notify', e.data)
        break
      case cc.wsApi.tableSnatch:
        window.$socket.emit('canrob_notify', e.data)
        break
    }
  },
  getTable() {
    const data = {
      userId: myglobal.playerData.userId,
      tableId: myglobal.playerData.tableId,
    }
    http.post(http.getTable, data, (res) => {
        if (res.code) {
            this.tipNode && this.tipNode.destroy();
            this.tipNode = cc.instantiate(this.tip);
            this.tipNode.getComponent(cc.Label).string = res.msg;
            this.tipNode.parent = this.node;
            return;
        }
        const {status, creator, players} = res.data.table

        myglobal.playerData.creator = creator.id
        ddzData.gameState = status 
        this.gameStateHandler(status)
        for (var i = 0; i < this.playerNodeList.length; i++) {
          var node = this.playerNodeList[i]
          node.destroy()
        }
        this.playerNodeList = []
        players.map(c => {
          this.addPlayerNode(c)
          c.cards && c.cards.length && window.$socket.emit('pushcard_notify', c.cards)
        })
    })
  },
  // 通知谁是地主, 并显示底牌
  masterNotify({ masterId, cards }) {
    // 必须先设置全局地主id
    myglobal.playerData.masterUserId = masterId
    // 显示底牌
    this.gameUiNode.emit("show_bottom_card_event", cards)
    for (var i = 0; i < this.playerNodeList.length; i++) {
      var node = this.playerNodeList[i]
      if (node) {
        // 给playernode节点发送事件
        node.emit("playernode_changemaster_event", masterId)
      }
    }
  },
  gameStateHandler(state) {
    this.btn_ready.active = state < ddzConsts.gameStatus.GAMESTART
    if (state === ddzConsts.gameStatus.WAITREADY) {
      this.btn_ready.active = true
    }
  },
  // 返回大厅
  onGoback() {
    const data = {
        userId: myglobal.playerData.userId,
    }
    http.post(http.exitRoom, data, (res) => {
        if (res.code) {
            this.tipNode && this.tipNode.destroy();
            this.tipNode = cc.instantiate(this.tip);
            this.tipNode.getComponent(cc.Label).string = res.msg;
            this.tipNode.parent = this.node;
            return;
        }
        myglobal.playerData.roomId = ''
        myglobal.playerData.tableId = ''
        cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
        cc.director.loadScene('hallScene')
    })
    // cc.ws.send({ action: cc.wsApi.roomLeave })
  },
  // 准备
  onBtnReadey(event) {
    cc.ws.send({ action: cc.wsApi.playerReady })
  },
  //seat_index自己在房间的位置id
  // setPlayerSeatPos(seat_index) {
  //   if (seat_index < 1 || seat_index > 3) {
  //     console.log("seat_index error" + seat_index)
  //     return
  //   }

  //   console.log("setPlayerSeatPos seat_index:" + seat_index)

  //   //界面位置转化成逻辑位置
  //   switch (seat_index) {
  //     case 1:
  //       this.playerdata_list_pos[1] = 0
  //       this.playerdata_list_pos[2] = 1
  //       this.playerdata_list_pos[3] = 2
  //       break
  //     case 2:


  //       this.playerdata_list_pos[2] = 0
  //       this.playerdata_list_pos[3] = 1
  //       this.playerdata_list_pos[1] = 2
  //       break
  //     case 3:
  //       this.playerdata_list_pos[3] = 0
  //       this.playerdata_list_pos[1] = 1
  //       this.playerdata_list_pos[2] = 2
  //       break
  //     default:
  //       break
  //   }
  // },
  // 添加玩家节点
  addPlayerNode(player_data) {
    const {id, next_id, table_id, ready} = player_data
    if (table_id !== myglobal.playerData.tableId) return

    var index = id === myglobal.playerData.userId ? 0 : next_id === myglobal.playerData.userId ? 2 : 1
    player_data.seat_index = index
    if (index === 0) {
      this.btn_ready.active = !ready // 准备按钮
    }
    
    var playernode_inst = cc.instantiate(this.player_node_prefabs)
    playernode_inst.parent = this.players_seat_pos.children[index]
    playernode_inst.seat_index = index
    //创建的节点存储在gamescene的列表中
    this.playerNodeList.push(playernode_inst)

    playernode_inst.getComponent("player_node").init_data(player_data)

    // myglobal.playerData.playerList[index] = player_data
  },

  /*
   //通过userId获取用户出牌放在gamescend的位置 
   做法：先放3个节点在gameacene的场景中 cardsoutzone(012)
  */
  getUserOutCardPosByAccount(userId) {
    for (var i = 0; i < this.playerNodeList.length; i++) {
      var node = this.playerNodeList[i]
      if (node) {
        //获取节点绑定的组件
        var node_script = node.getComponent("player_node")
        //如果accountid和player_node节点绑定的accountid相同
        //接获取player_node的子节点
        if (node_script.userId === userId) {
          var seat_node = this.players_seat_pos.children[node_script.seat_index].getChildByName('cardsoutzone')
          return seat_node
        }
      }
    }
    return null
  },
  /**
    * @description 通过userId获取玩家头像节点 
    * @param {String} userId 
    * @returns {cc.Node} 玩家节点
    */
  getUserNodeByAccount(userId) {
    for (let i = 0; i < this.playerNodeList.length; i++) {
      const node = this.playerNodeList[i]
      if (node) {
        //获取节点绑定的组件
        const playerNode = node.getComponent("player_node")
        if (playerNode.userId === userId) return playerNode
      }
    }
    return null
  }
});
