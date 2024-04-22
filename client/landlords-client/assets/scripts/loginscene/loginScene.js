/*
 * @Author: X1-EXT\lylin lylin888@163.com
 * @Date: 2022-05-09 13:33:49
 * @LastEditors: X1-EXT\lylin lylin888@163.com
 * @LastEditTime: 2024-04-18 15:05:52
 * @FilePath: \landlords-client\assets\scripts\loginscene\loginScene.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
cc.Class({
  extends: cc.Component,

  properties: {
    wait_node: cc.Node,
    protocaolNode: cc.Node,
    loginPanel: cc.Node,
    registerPanel: cc.Node
  },

  // LIFE-CYCLE CALLBACKS:

  onLoad() {
    cc.director.preloadScene("hallScene")
  },

  start() { },
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
      case 'user_login':
        this.loginPanel.active = true
        break
      case 'register':
        this.registerPanel.active = true
        break
      default:
        break
    }
  }
});
