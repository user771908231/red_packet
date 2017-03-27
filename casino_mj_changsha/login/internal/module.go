package internal

import (
	"github.com/name5566/leaf/module"
	"casino_mj_changsha/base"
	"fmt"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
}

func (m *Module) OnDestroy() {
	fmt.Println("")

}
