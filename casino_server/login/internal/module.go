package internal

import (
	"github.com/name5566/leaf/module"
	"casino_server/base"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

/**
初始化模块的不走
1,注册模块,就是把模块放置在一个list中方便管理
2,初始化模块,此处的OnInit()方法,需要注意的是,OnInit中调用了Skeleton的的很多初始化函数
3,执行Run方法,由于此处没有Run方法,所以实行的是组合*module.Skeleton的Run方法
 */

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
}

func (m *Module) OnDestroy() {

}
