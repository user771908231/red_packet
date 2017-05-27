package internal

import (
	"github.com/name5566/leaf/module"
	"casino_templet/base"
	"casino_templet/core/api"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
	roomMgr  = OinitMJRoomMgr(skeleton)
)

type Module struct {
	*module.Skeleton
	roomMgr api.RoomManagerApi
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	m.roomMgr = roomMgr
	//

}

func (m *Module) OnDestroy() {

}
