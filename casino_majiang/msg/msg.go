package msg

import (
	"github.com/name5566/leaf/network/protobuf"
	casinoProto "casino_server/msg/bbprotogo"
	majiangProto "casino_majiang/msg/bbprotogo"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&casinoProto.NullMsg{})                      //0
	Processor.Register(&majiangProto.Game_CreateRoom{})             //1 创建房间		game_RoomTypeInfo {
	Processor.Register(&majiangProto.Game_AckCreateRoom{})		//2

	Processor.Register(&majiangProto.Game_Ready{})                  //
	Processor.Register(&majiangProto.Game_ExchangeCards{})          //换牌
	Processor.Register(&majiangProto.Game_AckExchangeCards{})
	Processor.Register(&majiangProto.Game_DingQue{})                //定缺
	Processor.Register(&majiangProto.Game_Opening{})                //开始(表示都已经准备完了)

}
