import myglobal from "../mygolbal.js";
import http from "../util/http.js";
import ws from '../util/websocket.js'
cc.Class({
    extends: cc.Component,

    properties: {
        nickname_label: cc.Label,
        headimage: cc.Sprite,
        gobal_count: cc.Label,
        creatroom_prefabs: cc.Prefab,
        joinroom_prefabs: cc.Prefab,
        moreTip: cc.Node,
        tip: cc.Prefab,
    },

    // LIFE-CYCLE CALLBACKS:
    start() {
        console.log('start', cc.ws._isConnected)
    },
    onLoad() {
        ws.initWS()

        this.nickname_label.string = myglobal.playerData.userName;
        cc.director.preloadScene("gameScene");
    },
    onDestroy() {
        console.log('onDestroy')
    },
    
    onMessage(e) {
        console.log('onMessage', e, this)
    },
    // update (dt) {},

    onButtonClick(event, customData) {
        switch (customData) {
            case "create_room":
                var creator_Room = cc.instantiate(this.creatroom_prefabs);
                creator_Room.parent = this.node;
                creator_Room.zIndex = 100;
                break;
            case "join_room":
                var join_Room = cc.instantiate(this.joinroom_prefabs);
                join_Room.parent = this.node;
                join_Room.zIndex = 100;
                break;
            case "logout":
                const data = {
                    id: myglobal.playerData.userId,
                    name: myglobal.playerData.userName,
                }
                http.get(http.logout, data, (res) => {
                    console.log(res);
                    if (res.code) {
                        this.tipNode && this.tipNode.destroy();
                        this.tipNode = cc.instantiate(this.tip);
                        this.tipNode.getComponent(cc.Label).string = res.msg;
                        this.tipNode.parent = this.node;
                        return;
                    }
                    myglobal.playerData = {}
                    cc.sys.localStorage.clear()
                    cc.director.loadScene('loginScene')
                    cc.ws.close()
                    cc.ws = null
                });
                break;
            default:
                break;
        }
    },
    onBtnJingdian() {
        const creator_Room = cc.instantiate(this.creatroom_prefabs);
        creator_Room.parent = this.node;
        creator_Room.zIndex = 100;
    },
    onBtnLaizi() {
        this.moreTip.active = !this.moreTip.active;
    },
});
