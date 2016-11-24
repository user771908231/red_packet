package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	majiangProto "casino_majiang/msg/protogo"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&majiangProto.Heartbeat{})                      //0
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
	Processor.Register(&majiangProto.Game_AckSendOutCard{})         //19 出牌-ack
	Processor.Register(&majiangProto.Game_ActPeng{})            	//20 碰
	Processor.Register(&majiangProto.Game_AckActPeng{})           	//21
	Processor.Register(&majiangProto.Game_ActGang{})            	//22 杠
	Processor.Register(&majiangProto.Game_AckActGang{})            	//23
	Processor.Register(&majiangProto.Game_ActGuo{})            		//24 过
	Processor.Register(&majiangProto.Game_AckActGuo{})            	//25
	Processor.Register(&majiangProto.Game_ActHu{})            		//26 胡
	Processor.Register(&majiangProto.Game_AckActHu{})            	//27
	Processor.Register(&majiangProto.Game_BroadcastBeginDingQue{})	//28 开始定缺(广播)
	Processor.Register(&majiangProto.Game_BroadcastBeginExchange{})	//29 开始换牌(广播)
	Processor.Register(&majiangProto.Game_OverTurn{})        		//30 轮到下一人
	Processor.Register(&majiangProto.Game_SendCurrentResult{})    	//31 本局结果
	Processor.Register(&majiangProto.Game_SendEndLottery{})    		//32 牌局结束
	Processor.Register(&majiangProto.Game_DissolveDesk{})    		//33 解散房间
	Processor.Register(&majiangProto.Game_AckDissolveDesk{})    	//34
	Processor.Register(&majiangProto.Game_LeaveDesk{})    			//35 离开房间
	Processor.Register(&majiangProto.Game_AckLeaveDesk{})    		//36
	Processor.Register(&majiangProto.Game_Message{})    			//37 发送聊天消息
	Processor.Register(&majiangProto.Game_SendMessage{})    		//38 广播聊天
	Processor.Register(&majiangProto.Game_DingQueEnd{})			//39 定缺结束
	Processor.Register(&majiangProto.Game_GameRecord{})			//40 查询战绩
	Processor.Register(&majiangProto.Game_AckGameRecord{})			//41 战绩回复
	Processor.Register(&majiangProto.Game_ExchangeCardsEnd{})		//42 换三张 结束之后的广播
}
