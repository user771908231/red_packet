package api

type PDKRoomMgr interface {
	OnInit() error                                            //初始化room管理器
	GetRoom(roomType int32, roomLevel int32) (PDKRoom, error) //通过room类型和level等级获得一个room
	GetDeskBySession(userId uint32) (PDKDesk, error)          //通过玩家的session找到一个desk
}
