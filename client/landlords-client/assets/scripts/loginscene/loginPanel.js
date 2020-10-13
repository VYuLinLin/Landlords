// Learn cc.Class:
//  - https://docs.cocos.com/creator/manual/en/scripting/class.html
// Learn Attribute:
//  - https://docs.cocos.com/creator/manual/en/scripting/reference/attributes.html
// Learn life-cycle callbacks:
//  - https://docs.cocos.com/creator/manual/en/scripting/life-cycle-callbacks.html

cc.Class({
    extends: cc.Component,

    properties: {
        account: cc.Label,
        accountTip: cc.Node,
        password: cc.Label,
        passwordTip: cc.Node,
        tip: cc.Prefab,
        wait: cc.Node,
        // bar: {
        //     get () {
        //         return this._bar;
        //     },
        //     set (value) {
        //         this._bar = value;
        //     }
        // },
    },

    // LIFE-CYCLE CALLBACKS:

    // onLoad () {},

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
        myglobal.socket[`request_${val}`](data, function(status, res) {
            this.wait.active = false
            if (!status) {
            //   console.log("err:" + res)
              this.tipNode && this.tipNode.destroy()
              this.tipNode = cc.instantiate(this.tip)
              this.tipNode.getComponent(cc.Label).string = res
              this.tipNode.parent = this.node
              return
            }
            const {userId, userName} = res || {}
            myglobal.playerData.userId = userId
            myglobal.playerData.userName = userName
            cc.sys.localStorage.setItem('userData', JSON.stringify(myglobal.playerData))
            cc.director.loadScene("hallScene")
          }.bind(this))
    },
    closeHandler() {
        this.node.active = false
    }
});
