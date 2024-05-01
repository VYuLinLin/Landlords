export default {

    REQ_HEART: 1,
    RES_HEART: 2,

    REQ_CHEAT: 1,
    RSP_CHEAT: 2,

    REQ_LOGIN: 11,
    RSP_LOGIN: 12,

    REQ_ROOM_LIST: 13,
    RSP_ROOM_LIST: 14,

    REQ_TABLE_LIST: 15,
    RSP_TABLE_LIST: 16,

    REQ_JOIN_ROOM: 17,
    RSP_JOIN_ROOM: 18,

    REQ_JOIN_TABLE: 19,
    RSP_JOIN_TABLE: 20,

    REQ_NEW_TABLE: 21,
    RSP_NEW_TABLE: 22,

    REQ_DEAL_POKER: 31,
    RSP_DEAL_POKER: 32,

    REQ_CALL_SCORE: 33,
    RSP_CALL_SCORE: 34,

    REQ_SHOW_POKER: 35,
    RSP_SHOW_POKER: 36,

    REQ_SHOT_POKER: 37,
    RSP_SHOT_POKER: 38,

    REQ_GAME_OVER: 41,
    RSP_GAME_OVER: 42,

    REQ_CHAT: 43,
    RSP_CHAT: 44,

    REQ_RESTART: 45,
    RSP_RESTART: 46,

    roomList: 'room/list',
    roomJoin: 'room/join',    // 玩家进入房间
    roomJoinOther: 'room/join/other',    // 其他玩家进入房间
    roomLeave: 'room/leave',    // 玩家离开
    tableInfo: 'table/info', // 当前桌子和玩家信息
    tableUpdate: 'table/update', // 桌子状态更新
    userUpdate: 'user/update',
    playerReady: 'player/ready',
    playerDeal: 'player/deal', // 玩家发牌
    tableCallPoints: 'table/callPoints', // 抢地主
    PlayerCallPoints: 'player/callPoints', // 玩家叫分

}	
