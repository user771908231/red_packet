package api

type MjDesk interface {
	EnterUser(userId uint32) error //进入游戏
	Leave(userId uint32) error     //进入游戏
	GetMJConfig() interface{}      //获得一个麻将的配置
}
