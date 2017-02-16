package api

import "casino_majianagv2/core/data"

type MjUser interface {
	GetStatus() *data.MjUserStatus
	GetGameData() *data.MJUserGameData
	GetUserData() *data.MJUserData
	Ready()
	GetUserId() uint32 //
}
