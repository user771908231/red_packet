package internal

import (
	"github.com/name5566/leaf/module"
	"casino_server/base"
	"casino_server/common/log"
	"casino_server/service/room"
	"fmt"
)


/**
初始化模块的不走
1,注册模块,就是把模块放置在一个list中方便管理
2,初始化模块,此处的OnInit()方法,需要注意的是,OnInit中调用了Skeleton的的很多初始化函数
3,执行Run方法,由于此处没有Run方法,所以实行的是组合*module.Skeleton的Run方法
 */
var (
	skeleton = base.NewSkeleton()
	ChanRPC = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	log.T("game internal OnInit()")
	m.Skeleton = skeleton
	//开始锦标赛的游戏
	room.ChampionshipRoom.Begin()        //锦标赛开始

	//thgamero
	room.ThGameRoomIns.Recovery()        //恢复上一场的游戏数据
	
}

func (m *Module) OnDestroy() {
	fmt.Print("game module onDestroy--")
	// todo 这里需要做一些异常保存的操作... 现在还不知道crash的时候,是否可以保存
}
