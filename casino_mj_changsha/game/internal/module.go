package internal

import (
	"github.com/name5566/leaf/module"
	"casino_mj_changsha/base"
	"casino_mj_changsha/service/majiang"
	"casino_common/common/service/taskService/task"
	"casino_common/common/service/pushService"
	"casino_mj_changsha/conf"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	//回复数据
	//majiang.RecoverFMJ()        //回复麻将朋友桌的数据
	majiang.OnitRoomMnager(skeleton) //初始化麻将管理器

	//初始化任务系统
	task.InitTask()

	//载入push service初始化配置
	pushService.PoolInit(conf.Server.HallTcpAddr)
}

func (m *Module) OnDestroy() {
	//销毁时候需要做的操作
	//db.CloseMGO() //销毁数据库
}
