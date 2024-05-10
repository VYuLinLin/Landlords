/*
 * @Author: v vvv@888.com
 * @Date: 2022-05-09 13:33:49
 * @LastEditors: v vvv@888.com
 * @LastEditTime: 2024-05-07 13:03:27
 * @FilePath: \landlords-client\assets\scripts\data\ddzData.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/**
 * Date: 2020/2/21
 * Author: 
 * Desc:斗地主 数据层
 */
const DataNotify = require('DataNotify')

module.exports = {
  /**
   * @description 当前游戏状态
   */
  gameState: -1,

  initData() {
    this.gameStateNotify = DataNotify.create(this, 'gameState', JSON.parse(cc.sys.localStorage.getItem('gameState')))
    this.gameStateNotify.addListener(value => {
      cc.sys.localStorage.setItem('gameState', value)
    })
  }
}