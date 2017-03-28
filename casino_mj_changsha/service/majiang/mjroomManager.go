package majiang

import (
	"casino_common/common/log"
	"casino_common/common/sessionService"
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
	"casino_common/common/service/robotService"
	"github.com/name5566/leaf/module"
)

var MjroomManagerIns *MjRoomManager

//room管理器初始化
func OnitRoomMnager(s *module.Skeleton) {
	MjroomManagerIns = new(MjRoomManager) //初始化room管理器
	MjroomManagerIns.initFriendDesk(s)    //初始化朋友桌
}

//房间管理器
type MjRoomManager struct {
	FMJRoomIns  *MjRoom                     //朋友桌的room
	CMJRoomIns  []*MjRoom                   //金币场的room
	RobotManger *robotService.RobotsManager //机器人管理器
}

//初始化朋友桌
func (m *MjRoomManager) initFriendDesk(s *module.Skeleton) {
	m.FMJRoomIns = NewMjRoom(s)
	m.FMJRoomIns.RoomType = proto.Int32(ROOMTYPE_FRIEND)
}

//得到朋友桌的room
func (m *MjRoomManager) GetFMJRoom() *MjRoom {
	return MjroomManagerIns.FMJRoomIns
}

/**
	1,首先寻找玩家朋友桌的session
	2,没有找到朋友桌的session开始寻找金币场的session
	3,如果没有session 返回desk 是nil
 */
func (m *MjRoomManager) GetMjDeskBySession(userId uint32) *MjDesk {
	session := sessionService.GetFriendSession(userId, int32(ddproto.CommonEnumGame_GID_MAHJONG))
	if session == nil {
		session = sessionService.GetCoinSession(userId, int32(ddproto.CommonEnumGame_GID_MAHJONG))
	}
	return m.getDeskBySession(session) //
}

//
func (m *MjRoomManager) getDeskBySession(session *ddproto.GameSession) *MjDesk {
	log.T("通过玩家的session查找房间:%v", session)
	if session == nil {
		return nil
	}

	//这里需用通过roomId 得到朋友桌还是金币场
	deskId := session.GetDeskId()
	room := m.GetFMJRoom()
	if room == nil {
		//没有找到合适的room
		log.E("没有找到合适的room")
		return nil
	} else {
		return room.GetDeskByDeskId(deskId)
	}
}

//通过朋友桌的session得到桌子
func (m *MjRoomManager) GetFMjDeskBySession(userId uint32) *MjDesk {
	session := sessionService.GetSession(userId, ROOMTYPE_FRIEND)
	if session == nil || session.GetDeskId() == 0 {
		log.T("没有找到玩家[%v]朋友桌的session，现在开始找金币场的session", userId)
	}
	return m.getDeskBySession(session)

}
