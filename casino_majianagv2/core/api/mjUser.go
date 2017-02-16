package api

import (
	"casino_majianagv2/core/data"
	"casino_majiang/service/majiang"
)

type MjUser interface {
	GetStatus() *data.MjUserStatus
	GetGameData() *data.MJUserGameData
	GetUserData() *data.MJUserData
	Ready()
	GetUserId() uint32 //
	DelBillBean(pai *majiang.MJPai) (error, *majiang.BillBean)
	AddBill(relationUserid uint32, billType int32, des string, score int64, pai *majiang.MJPai, roomType int32) error
}
