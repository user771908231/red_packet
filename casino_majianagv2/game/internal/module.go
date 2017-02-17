package internal

import (
	"github.com/name5566/leaf/module"
	"casino_majianagv2/base"
	"casino_majianagv2/core/api"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
	roomMgr  = OinitMJRoomMgr(skeleton)
)

type Module struct {
	*module.Skeleton
	roomMgr api.MjRoomMgr
}

func (m *Module) GetRoomMgr() api.MjRoomMgr {
	return m.roomMgr
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	m.roomMgr = roomMgr
	//

}

func (m *Module) OnDestroy() {

}
