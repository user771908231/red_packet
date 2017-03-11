package coinPlay

import (
	"github.com/name5566/leaf/module"
	"casino_majianagv2/core/api"
	"casino_majianagv2/core/ins/skeleton"
	"casino_common/common/consts/tableName"
	"casino_majianagv2/core/data"
	"casino_common/common/log"
	"casino_common/utils/db"
	"casino_common/common/sessionService"
	"casino_common/common/consts"
	"casino_common/common/Error"
	"casino_majiang/service/majiang"
	"github.com/name5566/leaf/gate"
	"casino_majiang/msg/protogo"
	"casino_common/proto/ddproto"
	"fmt"
	"casino_common/utils"
	"casino_majiang/gamedata/model"
)

type CMjRoom struct {
	*module.Skeleton //leaf 的骨架
	*skeleton.SkeletonMJRoom
	RoomLevel   int32 //金币场等级
	EnterFee    int64 //进房费用
	CoinLimit   int64
	CoinLimitUL int64
	BaseValue   int64 //底分
}

//new一个room
func NewDefaultCMjRoom(mgr api.MjRoomMgr, s *module.Skeleton, config *model.TMjRoomConfig) api.MjRoom {
	ret := &CMjRoom{
		Skeleton:       s,
		RoomLevel:      config.RoomLevel,
		EnterFee:       config.EnterCoinFee,
		CoinLimit:      config.RoomLimitCoin,
		CoinLimitUL:    config.RoomLimitCoinUL,
		SkeletonMJRoom: skeleton.NewSkeletonMJRoom(mgr, config.RoomLevel),
	}
	ret.RoomName = config.RoomName
	ret.BaseValue = config.RoomBaseValue
	return ret
}

func (r *CMjRoom) GetRoomLevel() int32 {
	return r.RoomLevel
}

//room创建房间
func (r *CMjRoom) CreateDesk(config interface{}) (api.MjDesk, error) {
	c := config.(*data.SkeletonMJConfig)
	c.DeskId, _ = db.GetNextSeq(tableName.DBT_MJ_DESK)
	//根据不同的类型来得到不同地区的麻将
	desk := NewCMJDesk(c, r.Skeleton) //创建成都麻将朋友桌
	desk.SetRoom(r)                   //desk 关联room
	r.AddDesk(desk)                   //room 关联desk
	return desk, nil
}

//金币场 进入一个User
func (r *CMjRoom) EnterUser(userId uint32, key string, a gate.Agent) error {
	//金币场加入房间的逻辑
	desk, err := r.GetAbleCoinDeskBySession(userId)
	if desk == nil {
		log.T("通过key[%v]没有找到对应的desk,删除session", key)
		//删除玩家对应的session
		sessionService.DelSessionByKey(userId, majiang.ROOMTYPE_COINPLAY)
		return Error.NewError(consts.ACK_RESULT_FAIL, "长时间没有动作，请换房间")
	}

	err = desk.EnterUser(userId, a) //普通玩家进入房间
	if err != nil {
		//用户加入房间失败...
		sessionService.DelSession(sessionService.GetSession(userId, majiang.ROOMTYPE_COINPLAY))
		log.E("玩家[%v]加入房间失败errMsg[%v],删除session", userId, err)
		return err
	}

	return nil
}

//找到金币场可以玩的桌子
//通过session 找到对应的金币场的游戏desk
func (r *CMjRoom) GetAbleCoinDeskBySession(userId uint32) (api.MjDesk, error) {
	var retDesk api.MjDesk
	//得到金币场的session
	csession := sessionService.GetSession(userId, majiang.ROOMTYPE_COINPLAY)
	//如果csession == nil 表示之前没有在游戏状态 返回nil
	if csession == nil {
		log.T("没有找到玩家[%v]的session，返回空，需要寻找room中合适的房间...", userId)
		return nil, nil
	}

	//定义错误
	errMsg := "[%v]%v的游戏还没有结束，请完成后再来"
	gameName := utils.GetGameName(csession.GetGameId(), csession.GetRoomType())

	//如果金币的session 不是麻将的，直接返回
	if csession.GetGameId() != int32(ddproto.CommonEnumGame_GID_MAHJONG) {
		log.T("玩家[%v]当前的游戏状态没有在麻将金币场中..返回错误", userId)
		return nil, Error.NewError(consts.ACK_RESULT_ERROR, fmt.Sprintf(errMsg, gameName, ""))
	}

	//
	if csession.GetRoomId() != r.GetRoomId() {
		log.T("玩家[%v]当前的游戏状态room(id==%v)没有在当前room(id==%v)中..返回错误", userId, csession.GetRoomId(), r.GetRoomId())
		return nil, Error.NewError(consts.ACK_RESULT_ERROR,
			fmt.Sprintf(errMsg, gameName, r.GetRoomMgr().GetRoom(csession.GetRoomType(), csession.GetRoomId()).GetRoomName()))
	}

	//进入房间的error
	retDesk = r.GetDesk(csession.GetDeskId())
	if retDesk == nil {
		//玩家的session信息已经过时，直接删除
		log.T("通过玩家[%v]的session[%v]没有找到对应的desk，表示session过期，需要删除对应的session", userId, csession)
		sessionService.DelSessionByKey(userId, majiang.ROOMTYPE_COINPLAY) //通过session没有找到desk，证明session失效,此时删除session
		return nil, nil
	}

	//判断桌子的登记是否一样
	log.T("开始判断找到的desk的level[%v] 和 room 的level[%v] 是否相等", retDesk.GetRoom().GetRoomLevel(), r.GetRoomLevel())
	if retDesk.GetRoom().GetRoomLevel() != r.GetRoomLevel() {
		return nil, Error.NewError(consts.ACK_RESULT_ERROR,
			fmt.Sprintf(errMsg, gameName, r.GetRoomMgr().GetRoom(csession.GetRoomType(), csession.GetRoomId()).GetRoomName()))
	}

	//返回最终的数据
	return retDesk, nil
}

//创建金币场的desk
func (r *CMjRoom) CreateCoinPlayDesk() api.MjDesk {
	log.T("金币场没有空闲的房间，重新新建一个.")
	//创建金币场，创建时，需要更新mjroom的配置来创建
	config := &data.SkeletonMJConfig{
		Owner:            0,
		RoomId:           r.GetRoomId(),
		RoomType:         majiang.ROOMTYPE_COINPLAY,
		Status:           majiang.MJDESK_STATUS_READY,
		MjRoomType:       int32(mjproto.MJRoomType_roomType_xueZhanDaoDi),
		BoardsCout:       1000, //局数，如：4局（房卡 × 2）、8局（房卡 × 3）
		CreateFee:        r.EnterFee,
		CapMax:           8,
		CardsNum:         13,
		Settlement:       0,
		BaseValue:        100,
		ZiMoRadio:        1,
		OthersCheckBox:   nil,
		HuRadio:          1,
		DianGangHuaRadio: 1,
		MJPaiCursor:      0,
		TotalPlayCount:   10,
		CurrPlayCount:    0,
		Banker:           0,
		NextBanker:       0,
		ActiveUser:       0,
		ActUser:          0,
		ActType:          0,
		NInitActionTime:  30,
		RoomLevel:        0,
		PlayerCountLimit: 4,
		FangCount:        3,
		CoinLimit:        r.CoinLimit,
	}
	//
	desk, err := r.CreateDesk(config)
	if err != nil {
		log.E("创建金币场的房间的时候出错..err %v ", err)
	}
	return desk
}
