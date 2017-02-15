package api

import "casino_majianagv2/core/data"

type MjDesk interface {
	EnterUser(userId uint32) error            //进入游戏
	Leave(userId uint32) error                //进入游戏
	GetMJConfig() data.SkeletonMJConfig       //获得一个麻将的配置
	Ready(userId uint32) error                //
	DingQue(userId uint32, color int32) error //定缺
	GetStatus() *data.MjDeskStatus            //桌子的状态
}
