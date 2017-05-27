package internal

import (
	"github.com/name5566/leaf/module"
	"casino_paodekuai/base"
	"casino_paodekuai/core/api"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
	roomMgr  = OnInitRoomMgr(skeleton)
)

type Module struct {
	*module.Skeleton
	roomMgr api.PDKRoomMgr
}

func (m *Module) GetRoomMgr() api.PDKRoomMgr {
	return m.roomMgr
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton //赋值骨架
	if roomMgr != nil {
		m.roomMgr = roomMgr //赋值
		m.roomMgr.OnInit()  //初始化roommgr
	}
}

func (m *Module) OnDestroy() {

}
