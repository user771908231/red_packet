package api

import (
	"casino_majianagv2/core/data"
	"github.com/golang/protobuf/proto"
)

type MjUser interface {
	GetStatus() *data.MjUserStatus
	GetGameData() *data.MJUserGameData
	GetUserData() *data.MJUserData
	Ready()
	BeginInit(CurrPlayCount int32, banker uint32)
	GetUserId() uint32 //
	WriteMsg(p proto.Message) error
}
