package internal

import (
	"github.com/name5566/leaf/module"
	"casino_paodekuai/core/api"
	"casino_common/common/log"
	"casino_paodekuai/core/ins/friend"
	"casino_common/common/Error"
	"casino_common/common/consts"
)

func OnInitRoomMgr(s *module.Skeleton) api.PDKRoomMgr {
	ins := &PdkRoomIns{
		s:s,
	}
	ins.OnInitFRoom()
	return ins
}

type PdkRoomIns struct {
	s      *module.Skeleton //leaf骨架
	froom  api.PDKRoom      //朋友桌
	crooms []api.PDKRoom    //金币场
}


/************************ 所有方法的各类错误信息 统一定义在这里 ************************/
var ERR_ONINIT_1 = Error.NewError(consts.ACK_RESULT_ERROR, "初始化错误, 房间管理器为空")
var ERR_GETROOM_1 = Error.NewError(consts.ACK_RESULT_ERROR,"系统错误")


//初始化room管理器
func (m *PdkRoomIns) OnInit() error {
	//判断空
	if m == nil {
		log.F("初始化麻将管理器的时候失败，启动程序失败...")
		return ERR_ONINIT_1
	}

	//初始化朋友桌
	return nil
}

//初始化朋友桌
func (m *PdkRoomIns) OnInitFRoom() error {
	m.froom = friend.NewPDKFRoom(m.s)
	return nil
}

//通过room类型和level等级获得一个room
func (m *PdkRoomIns) GetRoom(roomType int32, roomLevel int32) (api.PDKRoom, error) {
	//todo 暂时返回朋友桌
	return m.froom, nil
}

func (m *PdkRoomIns) GetDeskBySession(userId uint32) (api.PDKDesk, error) {
	return nil, nil
}
