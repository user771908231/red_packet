package redModel

import (
	"casino_common/common/log"
)

type MsgContentType string

//标签类型
const (
	TagRoomId = "TagRoomId"
	TagSessionId = "TagSessionId"
	TagUserId = "UserId"
)

//消息类型
const (
	//红包
	ContentTypeRedPack MsgContentType = "RedPack"
	//文字
	ContentTypeText MsgContentType = "Text"

)

//通用消息体
type MsgHeader struct {
	FromUserId uint32
	FromClientId uint32

	ToUserId uint32
	ToClientId uint32
	ToGroupId int64

	Time int64

	//消息内容类型
	ContentType MsgContentType

	//------------文本----------
	//文本内容
	TextContent string

	//------------红包----------
	RedMoney float64  //金额

}

//处理消息
func (conn WsConn) HandlMsg(msg []byte) {
	log.T("msg:%v", msg)
	//屏蔽收消息
	return
	//广播给本群
	conn.BroadToCurrRoom(msg)
}

//新建立的链接
func (conn WsConn) HandConnect() {
	roomId := conn.Get(TagRoomId).(int32)
	room := GetRoomById(roomId)
	if room == nil {
		return
	}
	conn.WriteMsg(GetClientRedpackListJson(room.RedpackList))
}
