package api

import (
	"casino_majianagv2/core/data"
	"github.com/golang/protobuf/proto"
	"casino_majiang/service/majiang"
	"github.com/name5566/leaf/gate"
)

type MjUser interface {
	GetStatus() *data.MjUserStatus
	GetGameData() *data.MJUserGameData
	GetSkeletonUser() interface{} //返回骨架User
	ActReady()
	BeginInit(CurrPlayCount int32, banker uint32)
	GetUserId() uint32 //
	GetAgent() gate.Agent
	WriteMsg(p proto.Message) error
	DelBillBean(pai *majiang.MJPai) (error, *majiang.BillBean)
	AddBill(relationUserid uint32, billType int32, des string, score int64, pai *majiang.MJPai, roomType int32) error
	SendOverTurn(p proto.Message) error
	SetCoin(coin int64)             //设置用户金币
	ReEnterDesk(a gate.Agent) error //重新连接的处理
}
