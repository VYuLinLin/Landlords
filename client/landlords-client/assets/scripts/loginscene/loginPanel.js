/*
 * @Author: X1-EXT\lylin lylin888@163.com
 * @Date: 2022-05-09 13:33:49
 * @LastEditors: X1-EXT\lylin lylin888@163.com
 * @LastEditTime: 2024-04-18 15:10:25
 * @FilePath: \landlords-client\assets\scripts\loginscene\loginPanel.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

const http = require("../util/http");

cc.Class({
    extends: cc.Component,

    properties: {
        account: cc.Label,
        accountTip: cc.Node,
        password: cc.Label,
        passwordTip: cc.Node,
        tip: cc.Prefab,
        wait: cc.Node,
    },

    start() {},
    verifyHandler(e, val) {
        const account = this.account.string.trim()
        const password = this.password.string.trim()
        if (!val || val === 'account') {
            this.accountTip.active = !account
        }
        if (!val || val === 'password') {
            this.passwordTip.active = !password
        }
        return !!(account && password)
    },
    // 注册、登录
    loginHandler(e, val) {
        const data = {
            account: this.account.string.trim(),
            password: this.password.string.trim()
        }
        if (!this.verifyHandler()) return
        this.wait.active = true
        const url = val === 'login' ? http.login : http.register
        http.post(url, data, res => {
            this.wait.active = false
            if (res.code) {
                this.tipNode && this.tipNode.destroy()
                this.tipNode = cc.instantiate(this.tip)
                this.tipNode.getComponent(cc.Label).string = res.msg
                this.tipNode.parent = this.node
                return
            }
            const { id, name, coin, score, status, room_id, table_id } = res.data || {}
            myglobal.playerData.coin = coin
            myglobal.playerData.userId = id
            myglobal.playerData.userName = name
            myglobal.playerData.userScore = score
            myglobal.playerData.userStatus = status
            myglobal.playerData.tableId = table_id
            myglobal.playerData.roomId = room_id
            cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
            if (table_id) {
                cc.director.loadScene("gameScene")
            } else {
                cc.director.loadScene("hallScene")
            }
        })
    },
    closeHandler() {
        this.node.active = false
    }
});
