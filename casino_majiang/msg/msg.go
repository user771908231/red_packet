package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	casinoProto "casino_server/msg/bbprotogo"
	majiangProto "casino_majiang/msg/protogo"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&casinoProto.NullMsg{})                      //0
	Processor.Register(&majiangProto.Game_QuickConn{})             	//1 接入服务器
	Processor.Register(&majiangProto.Game_AckQuickConn{})          	//2
	Processor.Register(&majiangProto.Game_Login{})             		//3 登录游戏
	Processor.Register(&majiangProto.Game_AckLogin{})             	//4
	Processor.Register(&majiangProto.Game_CreateRoom{})             //5 创建房间
	Processor.Register(&majiangProto.Game_AckCreateRoom{})			//6

	Processor.Register(&majiangProto.Game_EnterRoom{})          	//7 进入房间
	Processor.Register(&majiangProto.Game_AckEnterRoom{})          	//8
	Processor.Register(&majiangProto.Game_SendGameInfo{})           //9 卓内游戏数据

	Processor.Register(&majiangProto.Game_Ready{})                  //10 准备
	Processor.Register(&majiangProto.Game_AckReady{})               //11
	Processor.Register(&majiangProto.Game_ExchangeCards{})          //12 换3张
	Processor.Register(&majiangProto.Game_AckExchangeCards{})		//13 换3张-回复
	Processor.Register(&majiangProto.Game_DingQue{})                //14 定缺
	Processor.Register(&majiangProto.Game_Opening{})                //15 开始(表示都已经准备完了)
	Processor.Register(&majiangProto.Game_DealCards{})              //16 发牌
	Processor.Register(&majiangProto.Game_GetInCard{})              //17 摸牌
	Processor.Register(&majiangProto.Game_SendOutCard{})            //18 出牌
	Processor.Register(&majiangProto.Game_ActPeng{})            	//19 碰
	Processor.Register(&majiangProto.Game_AckActPeng{})           	//20
	Processor.Register(&majiangProto.Game_ActGang{})            	//21 杠
	Processor.Register(&majiangProto.Game_AckActGang{})            	//22
	Processor.Register(&majiangProto.Game_ActGuo{})            		//23 过
	Processor.Register(&majiangProto.Game_AckActGuo{})            	//24
	Processor.Register(&majiangProto.Game_ActHu{})            		//25 胡
	Processor.Register(&majiangProto.Game_AckActHu{})            	//26
	Processor.Register(&majiangProto.Game_BroadcastBeginDingQue{})	//27 开始定缺(广播)
	Processor.Register(&majiangProto.Game_BroadcastBeginExchange{})	//28 开始换牌(广播)
	Processor.Register(&majiangProto.Game_OverTurn{})        		//29 轮到下一人
	Processor.Register(&majiangProto.Game_CurrentResult{})    		//30 本局结果
	Processor.Register(&majiangProto.Game_SendEndLottery{})    		//31 牌局结束
	Processor.Register(&majiangProto.Game_DissolveDesk{})    		//32 解散房间
	Processor.Register(&majiangProto.Game_AckDissolveDesk{})    	//33
	Processor.Register(&majiangProto.Game_LeaveDesk{})    			//34 离开房间
	Processor.Register(&majiangProto.Game_AckLeaveDesk{})    		//35
	Processor.Register(&majiangProto.Game_Message{})    			//36 发送聊天消息
	Processor.Register(&majiangProto.Game_SendMessage{})    		//37 广播聊天
}
