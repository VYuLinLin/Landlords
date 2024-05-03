import { EventEmitter } from "./eventEmitter";
import wsAPI from "./wsAPI"
cc.wsApi = wsAPI

export default {
    initWS() {
        if (cc.ws) return
        cc.ws = new WS().connect()
    }
}
/**
 * WebSocket的工具类
 *
 * @export
 * @class WS
 * @extends {EventEmitter}
 */
export class WS extends EventEmitter {
    constructor() {
        super();
        this._sock = null; // socket实例
        this._isConnected = false; // 是否连接成功
        this._pingIntervalSeconds = 3000; // 心跳连接时间
        this._heartTimer = null; // 心跳定时器
        this._serverTimer = null; // 服务器超时 定时器
        this._reconnectTimer = null; // 断开 重连倒计时
        this._sendFixHeartTimer = null; // 20s固定发送心跳定时器
    }
    OPEN = "open";
    ERROR = "error";
    CLOSE = "close";
    MESSAGE = "message";
    /**
     * 连接的方法
     *
     * @param {string} url
     * @memberof WS
     */
    connect(url = 'ws://localhost/ws') {
        if (!this._sock || [WebSocket.CLOSING, WebSocket.CLOSED].includes(this._sock.readyState)) {
            
            url += `?id=${myglobal.playerData.userId}`
            this._sock = new WebSocket(url);
            // this._sock.binaryType = 'arraybuffer';
            this._sock.onopen = this._onOpen.bind(this);
            this._sock.onclose = this._onClose.bind(this);
            this._sock.onerror = this._onError.bind(this);
            this._sock.onmessage = this._onMessage.bind(this);
        }
        console.log(this)
        return this;
    }
    /**
     * 开始连接的方法
     *
     * @private
     * @memberof WS
     */
    _onOpen(event) {
        this._isConnected = true;
        this.emit(this.OPEN, event);
        this._start()
        this._sendFixHeart()
    }
    /**
     * 错误的方法
     *
     * @private
     * @memberof WS
     */
    _onError(event) {
        console.log(event.type);
        this._isConnected = false;
        this.emit(this.ERROR, event);
        this._reconnect();
    }
    /**
     * 关闭的方法
     *
     * @private
     * @memberof WS
     */
    _onClose(event) {
        console.log(event.type);
        this._isConnected = false;
        this.emit(this.CLOSE, event);
        this._reconnect();
    }
    /**
     * 信息的方法
     *
     * @private
     * @param {MessageEvent} event
     * @memberof WS
     */
    _onMessage(event) {
        this._reset()
        const _this = this;
        if (event.data instanceof Blob) {
            if (event.data.text) {
                event.data.text().then(_this._parseMsgAndEmit);
            } else {
                // 兼容QQ浏览器等无Blob实例方法的浏览器
                var reader = new FileReader();
                reader.onload = function (event) {
                    _this._parseMsgAndEmit(event.target.result)
                };
                reader.readAsText(event.data);
            }
        } else {
            _this._parseMsgAndEmit(event.data)
        }
    }
    _parseMsgAndEmit(res) {
        res = JSON.parse(res)
        if (res === wsAPI.RES_HEART) return
        console.log("ws 消息: ", res)
        this.emit(this.MESSAGE, res);
    }
    /**
     * 发送数据
     *
     * @param {(string | object)} message
     * @memberof WS
     */
    send(message) {
        if (!this._isConnected) return;
        if (typeof message == "string") {
            this._sock.send(message);
        } else {
            let jsonStr = JSON.stringify(message);
            this._sock.send(jsonStr);
        }
    }
    /**
     * 发送二进制数据
     * @param message
     */
    sendBinary(message) {
        if (!this._isConnected) return;
        if (typeof message != "string") {
            this._sock.send(message);
        }
    }
    /**
     * 关闭连接
     *
     * @memberof WS
     */
    close() {
        this._sock.close();
        this._sock = null;
        this.allOff();
        
        this._isConnected = false;
        
        clearTimeout(this._heartTimer);
        clearTimeout(this._serverTimer);
        clearInterval(this._sendFixHeartTimer)
    }
    // 开始心跳
    _start() {
        this._serverTimer && clearTimeout(this._serverTimer);
        this._heartTimer && clearTimeout(this._heartTimer);
        const _this = this
        this._heartTimer = setTimeout(function() {
            _this.send(wsAPI.REQ_HEART);
            //超时关闭，超时时间为5s
            _this._serverTimer = setTimeout(function(){
                _this._sock.close();
            }, 5000);
        }, this._pingIntervalSeconds);
    }
    // 重新连接  3000-5000之间，设置延迟避免请求过多
    _reconnect(){
        this._reconnectTimer && clearTimeout(this._reconnectTimer);
        const _this = this
        this._reconnectTimer = setTimeout(()=> {
            _this._sock && _this.connect();
        }, parseInt(Math.random()*2000 + 3000));
    }
    // 重置心跳
    _reset(){
        this._start();
    }
    // 20s固定发送心跳
    _sendFixHeart(){
        this._sendFixHeartTimer && clearInterval(this._sendFixHeartTimer);
        this._sendFixHeartTimer = setInterval(()=>{
            this.send(wsAPI.REQ_HEART);
        }, 20000);
    }
}
