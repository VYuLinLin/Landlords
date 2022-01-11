// Learn cc.Class:
//  - https://docs.cocos.com/creator/manual/en/scripting/class.html
// Learn Attribute:
//  - https://docs.cocos.com/creator/manual/en/scripting/reference/attributes.html
// Learn life-cycle callbacks:
//  - https://docs.cocos.com/creator/manual/en/scripting/life-cycle-callbacks.html

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

    start () {

    },
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
            const {id, username, coin} = res.data || {}
            myglobal.playerData.coin = coin
            myglobal.playerData.userId = id
            myglobal.playerData.userName = username
            cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
            cc.director.loadScene("hallScene")
        })
    },
    closeHandler() {
        this.node.active = false
    }
});
