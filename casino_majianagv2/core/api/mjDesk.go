package api

import (
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/ins/skeleton"
)

type MjDesk interface {
	EnterUser(userId uint32) error               //进入游戏
	Ready(userId uint32) error                   //
	DingQue(userId uint32, color int32) error    //定缺
	ActOut(userId uint32, cardId int32) error    //打牌
	Leave(userId uint32) error                   //进入游戏
	GetSkeletonMJDesk() *skeleton.SkeletonMJDesk //骨架
	//基本方法
	GetMJConfig() *data.SkeletonMJConfig //获得一个麻将的配置
	GetStatus() *data.MjDeskStatus       //桌子的状态
}
