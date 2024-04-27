/*
 * @Author: v vvv@888.com
 * @Date: 2022-05-09 13:33:49
 * @LastEditors: v vvv@888.com
 * @LastEditTime: 2024-04-26 12:44:17
 * @FilePath: \landlords-client\assets\scripts\common\loadingLayer.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import myglobal from "../mygolbal.js"

var LoadingLayer = cc.Class({
  extends: cc.Component,
  properties: {},
  __preload() {
    this.init()
  },
  start() {
  },
  init() {
    // 调整画布前的回调
    cc.view.resizeWithBrowserSize(true);
    cc.view.setResizeCallback(this.resizeCallback)
    cc.director.on(cc.Director.EVENT_AFTER_SCENE_LAUNCH, this.resizeCallback);
    cc.game.addPersistRootNode(this.node);
    // 场景跳转
    const { userId, roomId } = myglobal.playerData
    console.log('userId = ', userId)
    console.log('roomId = ', roomId)
    if (!userId) {
      cc.director.loadScene('loginScene')
    } else if (!roomId) {
      cc.director.loadScene('hallScene')
    } else {
      cc.director.loadScene('gameScene')
    }
  },
  resizeCallback() {
    var canvas = cc.find("Canvas").getComponent(cc.Canvas)
    var t = cc.winSize.width / canvas.designResolution.width
    var n = cc.winSize.height / canvas.designResolution.height;
    t < n
      ? (canvas.fitWidth = !0, canvas.fitHeight = !1) : n < t
        ? (canvas.fitWidth = !1, canvas.fitHeight = !0) : (canvas.fitWidth = !1, canvas.fitHeight = !1)
  }
});