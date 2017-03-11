package internal

import (
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/ins/friendPlay"
	"github.com/name5566/leaf/module"
	"casino_common/common/service/robotService"
	"casino_common/proto/ddproto"
	"casino_majianagv2/core/ins/coinPlay"
	"casino_majiang/service/majiang"
	"casino_common/common/sessionService"
	"casino_common/common/log"
	"casino_majiang/gamedata/dao"
)

type MJMgr struct {
	froom       api.MjRoom
	croom       []api.MjRoom
	RobotManger robotService.RobotsMgrApi
	*module.Skeleton
}

//todo
func (m *MJMgr) GetDesk() api.MjDesk {
	return nil
}

//todo
func (m *MJMgr) GetRoom(roomType int32, roomLevel int32) api.MjRoom {
	if roomType == majiang.ROOMTYPE_FRIEND {
		return m.froom
	} else {
		for _, r := range m.croom {
			if r.GetRoomLevel() == roomLevel {
				return r
			}
		}
	}

	return nil
}

//room管理器的初始化工作应该放在这里
func (m *MJMgr) OnInit() error {
	m.initFriendDesk()  //初始化朋友桌
	m.initCoinDesk()    //初始化金币场
	m.initRobotManger() //初始化机器人
	return nil
}

//初始化朋友桌
func (m *MJMgr) initFriendDesk() {
	m.froom = friendPlay.NewDefaultFMJRoom(m, m.Skeleton)
}

//初始化金币场
func (m *MJMgr) initCoinDesk() {
	roomConfigs := dao.GetMJRoomConfigData()
	log.T("开始初始化麻将的金币场room：roomConfigs:%v", roomConfigs)
	for _, d := range roomConfigs {
		m.croom = append(m.croom, coinPlay.NewDefaultCMjRoom(m, m.Skeleton, d))
	}
	log.T("麻将的金币场room初始化腕尺:%v", m.croom)
}

//初始化机器人管理器
func (m *MJMgr) initRobotManger() {
	m.RobotManger = robotService.NewRobotManager(ddproto.CommonEnumGame_GID_MAHJONG)
}

func (m *MJMgr) GetMjDeskBySession(userId uint32) api.MjDesk {
	session := sessionService.GetSessionAuto(userId)
	if session != nil {
		return m.getDeskBySession(session)
	}
	return nil
}

func (m *MJMgr) getDeskBySession(session *ddproto.GameSession) api.MjDesk {
	log.T("通过玩家的session查找房间:%v", session)
	if session == nil {
		return nil
	}
	//这里需用通过roomId 得到朋友桌还是金币场
	roomId := session.GetRoomId()
	deskId := session.GetDeskId()
	roomType := session.GetRoomType()
	room := m.GetRoom(roomType, roomId)
	if room == nil {
		//没有找到合适的room
		log.E("没有找到合适的room")
		return nil
	} else {
		desk := room.GetDesk(deskId)
		if desk == nil {
			sessionService.DelSession(session)
		}
		return desk
	}
}

func (m *MJMgr) GetRobotManger() robotService.RobotsMgrApi {
	return m.RobotManger
}

//返回一个默认的mjroom管理器
func OinitMJRoomMgr(s *module.Skeleton) api.MjRoomMgr {
	return &MJMgr{
		Skeleton: s,
	}
}
