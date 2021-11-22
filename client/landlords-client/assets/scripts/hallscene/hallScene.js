import myglobal from "../mygolbal.js";
import http from "../util/http.js";
import { WS } from '../util/websocket.js'
import wsAPI from '../util/wsAPI.js'
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
        this.nickname_label.string = myglobal.playerData.userName;
        cc.director.preloadScene("gameScene");
        console.log('onLoad')
        if (!cc.ws) {
            cc.wsApi = wsAPI
            cc.ws = new WS().connect()
            // cc.ws.on(cc.ws.MESSAGE, this.onMessage.bind(this))
        }
    },
    onDestroy() {
        console.log('onDestroy', cc.ws._isConnected)
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
                http.get(http.logout, (res) => {
                    console.log(res);
                    if (res.code) {
                        this.tipNode && this.tipNode.destroy();
                        this.tipNode = cc.instantiate(this.tip);
                        this.tipNode.getComponent(cc.Label).string = res.msg;
                        this.tipNode.parent = this.node;
                        return;
                    }
                    cc.sys.localStorage.clear()
                    cc.director.loadScene('loginScene')
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
