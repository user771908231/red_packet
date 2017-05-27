package internal

import (
	"github.com/name5566/leaf/module"
	"casino_paodekuai/base"
	"casino_common/api/common/roomMgr"
)

var (
	skeleton   = base.NewSkeleton()
	ChanRPC    = skeleton.ChanRPCServer
	roomMgrIns = OnInitRoomMgr(skeleton)
)

type Module struct {
	*module.Skeleton
	roomMgr roomMgr.RoomMgrApi
}

func (m *Module) GetRoomMgr() roomMgr.RoomMgrApi {
	return m.roomMgr
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton //赋值骨架
	if roomMgr != nil {
		m.roomMgr = roomMgrIns //赋值
		m.roomMgr.OnInit()     //初始化roommgr
	}
}

func (m *Module) OnDestroy() {

}
