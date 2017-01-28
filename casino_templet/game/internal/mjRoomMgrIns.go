package internal

import (
	"casino_templet/core/api"
	"casino_templet/core/ins/friendPlay"
	"github.com/name5566/leaf/module"
)

type MJMgr struct {
	froom api.MjRoom
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

func (m *MJMgr) OnInit() error {

}

//返回一个默认的mjroom管理器
func OinitMJRoomMgr(s *module.Skeleton) api.MjRoomMgr {
	return &MJMgr{
		froom:   friendPlay.NewDefaultFMJRoom(s),
		Skeleton:s,
	}
}
