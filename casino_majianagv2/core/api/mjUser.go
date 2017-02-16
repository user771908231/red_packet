package api

import (
	"casino_majianagv2/core/data"
	"casino_majiang/service/majiang"
)

type MjUser interface {
	GetStatus() *data.MjUserStatus
	GetGameData() *data.MJUserGameData
	Ready()
	GetUserId() uint32 //
	DelBillBean(pai *majiang.MJPai) (error, *majiang.BillBean)
}
