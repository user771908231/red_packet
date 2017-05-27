package api

type PDKDesk interface {
	EnterUser(userId uint32) error
	ActOut(userId uint32, p interface{}) error //出牌的user和牌型
	ActReady(userId uint32) error

	//common func
	GetDeskId() int32 //得到desk id
	//GetDeskCfg() interface{} //得到desk的配置
}
