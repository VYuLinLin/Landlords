/*
 * @Author: X1-EXT\lylin lylin888@163.com
 * @Date: 2022-05-09 13:33:49
 * @LastEditors: X1-EXT\lylin lylin888@163.com
 * @LastEditTime: 2024-04-18 14:57:22
 * @FilePath: \landlords-client\assets\scripts\data\player.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
const getRandomStr = function (count) {
  var str = '';
  for (var i = 0; i < count; i++) {
    str += Math.floor(Math.random() * 10);
  }
  return str;
};

const playerData = function () {
  const userData = JSON.parse(cc.sys.localStorage.getItem('userData'))
  const [rootId1, rootId2] = [getRandomStr(5), getRandomStr(5)]
  var that = userData || {
    coin: 0, // 余额
    userId: '', // 用户id
    userName: '', // 用户名称，guest_ 开头
    userScore: 0, // 用户积分
    userStatus: 0, // 用户状态
    roomId: '',// 游戏房间id
    // roomName: '',
    creator: '', // 房主
    tableId: '', // 游戏桌子id
    tableStatus: 0, // 桌子状态
    // seatindex: 0, // 作为id
    avatarUrl: 'avatar_1', // 头像
    clients: [],
    // rootList: [
    //   { seatindex: 1, userId: rootId1, userName: `guest_${rootId1}`, "avatarUrl": "avatar_2", "coin": getRandomStr(4) },
    //   { seatindex: 2, userId: rootId2, userName: `guest_${rootId2}`, "avatarUrl": "avatar_3", "coin": getRandomStr(4) }
    // ],
    masterUserId: '', // 地主id
  }
  // that.uniqueID = 1 + getRandomStr(6)
  that.gobal_count = cc.sys.localStorage.getItem('user_count')
  // that.master_accountid = 0
  if (!userData) {
    cc.sys.localStorage.setItem('userData', JSON.stringify(that))
  }
  return that;
}
export default playerData
