package api

import "casino_common/common/service/robotService"

//麻将的房间管理器
type MjRoomMgr interface {
	GetDesk() MjDesk
	GetRoom(int32, int32) MjRoom
	OnInit() error
	//SetSkeleton(*module.Skeleton)
	GetMjDeskBySession(userId uint32) MjDesk
	GetRobotManger() robotService.RobotsMgrApi //得到机器人管理器
}
