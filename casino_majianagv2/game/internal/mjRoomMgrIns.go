package internal

import (
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/ins/friendPlay"
	"github.com/name5566/leaf/module"
	"casino_common/common/service/robotService"
	"casino_common/proto/ddproto"
)

type MJMgr struct {
	froom       api.MjRoom
	RobotManger robotService.RobotsMgrApi
	*module.Skeleton
}

//todo
func (m *MJMgr) GetDesk() api.MjDesk {
	return nil
}

//todo
func (m *MJMgr) GetRoom(roomType int32, roomLevel int32) api.MjRoom {
	return m.froom
}

//room管理器的初始化工作应该放在这里
func (m *MJMgr) OnInit() error {
	m.RobotManger = robotService.NewRobotManager(ddproto.CommonEnumGame_GID_MAHJONG)
	return nil
}

func (m *MJMgr) GetMjDeskBySession(userId uint32) api.MjDesk {
	return nil
}
func (m *MJMgr) GetRobotManger() robotService.RobotsMgrApi {
	return m.RobotManger
}

//返回一个默认的mjroom管理器
func OinitMJRoomMgr(s *module.Skeleton) api.MjRoomMgr {
	return &MJMgr{
		froom:    friendPlay.NewDefaultFMJRoom(s),
		Skeleton: s,
	}
}
