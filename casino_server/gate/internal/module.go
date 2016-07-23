package internal

import (
	"github.com/name5566/leaf/gate"
	"casino_server/conf"
	"casino_server/msg"
	"casino_server/game"
)


/**
初始化模块的不走
1,注册模块,就是把模块放置在一个list中方便管理
2,初始化模块,此处的OnInit()方法,需要注意的是,OnInit中调用了Gate的的很多初始化函数
3,执行Run方法,由于此处没有Run方法,所以实行的是组合*gate.Gate的Run方法
 */


type Module struct {
	*gate.Gate
}

func (m *Module) OnInit() {
	m.Gate = &gate.Gate{
		MaxConnNum:      conf.Server.MaxConnNum,
		PendingWriteNum: conf.PendingWriteNum,
		MaxMsgLen:       conf.MaxMsgLen,
		WSAddr:          conf.Server.WSAddr,
		HTTPTimeout:     conf.HTTPTimeout,
		TCPAddr:         conf.Server.TCPAddr,
		LenMsgLen:       conf.LenMsgLen,
		LittleEndian:    conf.LittleEndian,
		//Processor:       msg.Processor,
		Processor:       msg.ProtoProcessor,
		AgentChanRPC:    game.ChanRPC,
	}
}
