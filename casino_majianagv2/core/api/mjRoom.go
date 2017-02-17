package api

type MjRoom interface {
	CreateDesk(config interface{}) (MjDesk, error)
	GetDesk() MjDesk
	EnterUser(userId uint32, key string) error
	DissolveDesk(desk MjDesk, f bool) error //解散方剂爱你
	GetRoomMgr() MjRoomMgr                  //room管理器
}
