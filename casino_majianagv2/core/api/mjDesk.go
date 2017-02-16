package api

import (
	"casino_majianagv2/core/data"
)

type MjDesk interface {
	EnterUser(userId uint32) error                       //进入游戏
	Ready(userId uint32) error                           //
	DingQue(userId uint32, color int32) error            //定缺
	ActOut(userId uint32, cardId int32, auto bool) error //打牌
	ActPeng(userId uint32) error                         //碰牌
	ActGang(userId uint32, id int32, bu bool) error      //杠牌
	ActHu(userId uint32) error                           //胡牌
	Leave(userId uint32) error                           //进入游戏
	GetMJConfig() *data.SkeletonMJConfig                 //获得一个麻将的配置
	GetStatus() *data.MjDeskStatus                       //桌子的状态
	GetHuParser() HuPaerApi                              //得到胡牌解析器
}
