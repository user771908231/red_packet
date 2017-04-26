package internal

import (
	"casino_templet/core/api"
	"github.com/name5566/leaf/module"
)

type MJMgr struct {
	froom api.RoomApi
	*module.Skeleton
}
