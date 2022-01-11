/**
 * @description 斗地主常量
 */
module.exports = {
  roomNames: ["初级房", "中级房", "高级房", "大师房"],
  /**
   * @description 游戏状态
   */
  gameStatus: {
    INVALID: 0, // 无效
    WAITREADY: 1, // 进入房间，等待游戏
    GAMESTART: 2, // 开始游戏，已准备
    PUSHCARD: 3, // 发牌
    ROBSTATE: 4, // 抢地主
    SHOWBOTTOMCARD: 5, // 显示底牌
    PLAYING: 6, // 出牌阶段
    GAMEEND: 7, // 游戏结束
  },
  /**
   * @description 游戏状态
   */
  playerStatus: {
    invalid: 0, // 无效
    unready: 1, // 未准备
    ready: 2, // 已准备
    calling: 3, // 抢地主（叫分）
    double: 4, // 加倍
    playing: 5, // 出牌阶段
    over: 6, // 游戏结束
  },
  /**
   * @description 扑克牌纹理
   */
  _pokersFrame: null,
  _chipsFrame: null,
};
