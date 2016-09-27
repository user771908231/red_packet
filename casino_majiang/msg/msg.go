package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	casinoProto "casino_server/msg/bbprotogo"
	majiangProto "casino_majiang/msg/protogo"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&casinoProto.NullMsg{})                      //0
	Processor.Register(&majiangProto.Game_CreateRoom{})             //1 创建房间		game_RoomTypeInfo {
	Processor.Register(&majiangProto.Game_AckCreateRoom{})			//2

	Processor.Register(&majiangProto.Game_EnterRoom{})          	//3 进入房间
	Processor.Register(&majiangProto.Game_AckEnterRoom{})			//4

	Processor.Register(&majiangProto.Game_Ready{})                  //5 准备
	Processor.Register(&majiangProto.Game_ExchangeCards{})          //6 换牌
	Processor.Register(&majiangProto.Game_AckExchangeCards{})		//7 换3张
	Processor.Register(&majiangProto.Game_DingQue{})                //8 定缺
	Processor.Register(&majiangProto.Game_Opening{})                //9 开始(表示都已经准备完了)

}
