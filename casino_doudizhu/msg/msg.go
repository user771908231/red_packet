package msg

import (
	"github.com/name5566/leaf/network/protobuf"
)

var Processor = protobuf.NewProcessor()

/**

 PID_QUICK_CONN = 1; //连接服务器
    PID_QUICK_CONN_ACK = 2; //
    PID_GAME_LOGIN = 3; //登录游戏
    PID_GAME_LOGIN_ACK = 4; //
    PID_CREATEROOM = 5; //创建房间
    PID_CREATEROOM_ACK = 6; //创建房间-回复
    PID_ENTER_ROOM = 7; //进入房间
    PID_ENTER_ROOM_ACK = 8; //进入房间-回复
    PID_SEND_GAMEINFO = 9; //游戏信息
    PID_READY = 10; //准备
    PID_READY_ACK = 11; //准备-ack
    PID_OPENING = 12; //开局
    PID_DEAL_CARDS = 13; //发牌

    PID_CALL_DIZHU = 14; //叫地主
    PID_CALL_DIZHU_ACK = 15; //叫地主-ack

    ////////////////////////////////////////////
    //欢乐斗地主
    PID_ROB_DIZHU = 16; //抢地主
    PID_ROB_DIZHU_ACK = 17; //抢地主-ack
    PID_DOUBLE = 18; //加倍
    PID_DOUBLE_ACK = 19; //加倍-ack
    PID_NOT_DOUBLE = 20; //不加倍
    PID_NOT_DOUBLE_ACK = 21; //不加倍-ack

    PID_SHOW_HANDPOKERS = 22; //明牌
    PID_SHOW_HANDPOKERS_ACK = 23; //明牌-ack

    ////////////////////////////////////////////
    //四川斗地主
    PID_MEN_ZHUA = 24; //闷抓
    PID_MEN_ZHUA_ACK = 25; //闷抓-ack
    PID_LOOK_POKER = 26; //看牌
    PID_LOOK_POKER_ACK = 27; //看牌-ack
    PID_DAO = 28; //倒
    PID_DAO_ACK = 29; //倒-ack
    PID_LA = 30; //拉
    PID_LA_ACK = 31; //拉-ack

    PID_SEND_OUT_POKER = 32; //出牌
    PID_SEND_OUT_POKER_ACK = 33; //出牌-ack
    PID_GUO_POKER = 34; //过
    PID_GUO_POKER_ACK = 35; //过-ack

    ////////////////////////////////////////////

    PID_OVERTURN = 36; //轮到下一人操作(广播)

    PID_CURRENTRESULT = 37; //本局结果
    PID_SENDENDLOTTERY = 38; //牌局结束

    PID_DISSOLVE_DESK = 39; //离开房间
    PID_DISSOLVE_DESK_ACK = 40; //
    PID_LEAVE_DESK = 41; //解散房间
    PID_LEAVE_DESK_ACK = 42; //

    PID_MESSAGE = 43; //聊天信息
    PID_SEND_MESSAGE = 44; //广播聊天

    PID_GAME_GAMERECORD = 45; //查询战绩
    PID_GAME_ACKGAMERECORD = 46; //战绩回复
 */
func init() {

}
