package api

import "casino_majianagv2/core/data"

type MjUser interface {
	GetStatus() *data.MjUserStatus
	GetGameData() *data.MJUserGameData
	Ready()
	GetUserId() uint32 //
}
