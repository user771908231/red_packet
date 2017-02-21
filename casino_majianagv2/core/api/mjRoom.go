package api

type MjRoom interface {
	CalcCreateFee(c int32) int64                   //计算创建房间需要的费用
	CreateDesk(config interface{}) (MjDesk, error) //创建房间
	GetRoomId() int32
	GetRoomLevel() int32
	GetDesk(id int32) MjDesk                   //得到一个desk
	EnterUser(userId uint32, key string) error //进入一个玩家
	DissolveDesk(desk MjDesk, f bool) error    //解散方剂爱你
	GetRoomMgr() MjRoomMgr                     //room管理器
}