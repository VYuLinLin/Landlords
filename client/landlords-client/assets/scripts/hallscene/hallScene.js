/*
 * @Author: X1-EXT\lylin lylin888@163.com
 * @Date: 2022-05-09 13:33:49
 * @LastEditors: X1-EXT\lylin lylin888@163.com
 * @LastEditTime: 2024-04-20 19:24:35
 * @FilePath: \landlords-client\assets\scripts\hallscene\hallScene.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import myglobal from "mygolbal.js";
import http from "http.js";
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
    start() {},
    onLoad() {
        this.nickname_label.string = myglobal.playerData.userName;
        cc.director.preloadScene("gameScene");
    },
    onDestroy() {
        console.log('onDestroy')
    },
    onButtonClick(event, customData) {
        switch (customData) {
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
                http.post(http.logout, data, (res) => {
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
