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
)

type CMjRoom struct {
	*module.Skeleton //leaf 的骨架
	*skeleton.SkeletonMJRoom
	RoomLevel int32 //金币场等级
	EnterFee  int64 //进房费用
}

func NewDefaultCMjRoom(s *module.Skeleton, l int32) api.MjRoom {
	ret := &CMjRoom{
		Skeleton:       s,
		RoomLevel:      l,
		SkeletonMJRoom: skeleton.NewSkeletonMJRoom(l),
	}
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
	desk.SetRoom(r)
	return desk, nil
}

//金币场 进入一个User
func (r *CMjRoom) EnterUser(userId uint32, key string, a gate.Agent) error {
	//金币场加入房间的逻辑
	desk := r.GetAbleCoinDesk(userId)
	if desk == nil {
		log.T("通过key[%v]没有找到对应的desk,删除session", key)
		//删除玩家对应的session
		sessionService.DelSessionByKey(userId, majiang.ROOMTYPE_COINPLAY)
		return Error.NewError(consts.ACK_RESULT_FAIL, "长时间没有动作，请换房间")
	}

	err := desk.EnterUser(userId, a) //普通玩家进入房间
	if err != nil {
		//用户加入房间失败...
		sessionService.DelSession(sessionService.GetSession(userId, majiang.ROOMTYPE_COINPLAY))
		log.E("玩家[%v]加入房间失败errMsg[%v],删除session", userId, err)
		return err
	}

	return nil
}

//找到金币场可以玩的桌子
func (r *CMjRoom) GetAbleCoinDesk(userId uint32) api.MjDesk {
	var retDesk api.MjDesk
	//1,通过session来找
	csession := sessionService.GetSession(userId, majiang.ROOMTYPE_COINPLAY)
	//log.T("玩家%v进入金币场的时候，获取到的session：%v", userId, csession)
	if csession != nil {
		retDesk = r.GetDesk(csession.GetDeskId())
	}

	//如果session中寻找到的desk不为空，那么直接返回，如果没有找到，那么再从新分配
	if retDesk != nil {
		return retDesk
	}

	//2,session中没有找到原来的desk，重新分配桌子
	//log.T("开始获取可以使用的金币场level[%v]的desk....", r.GetRoomLevel())
	for _, d := range r.Desks {
		if d != nil && d.GetSkeletonMjDesk() != nil {
			s := d.GetSkeletonMjDesk().(*skeleton.SkeletonMJDesk)
			if !s.IsPlayerEnough() {
				retDesk = d
			}
		}
	}

	//创建一个金币场的桌子
	if retDesk == nil {
		retDesk = r.CreateCoinPlayDesk()
	}
	log.T("获取到的可以使用level[%v]的金币场的desk[%v].", r.GetRoomLevel(), retDesk.GetMJConfig().DeskId)
	return retDesk
}

//创建金币场的desk
func (r *CMjRoom) CreateCoinPlayDesk() api.MjDesk {
	//创建金币场，创建时，需要更新mjroom的配置来创建
	config := &data.SkeletonMJConfig{
		Owner:            0,
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
		CoinLimit:        r.GetCoinLimit(),

	}
	//
	desk, err := r.CreateDesk(config)
	if err != nil {
		log.E("创建金币场的房间的时候出错..err %v ", err)
	}
	return desk
}

//进入房间的费用 todo 需要移动到配置文件
func (r *CMjRoom) GetEnterFee() int64 {
	return 300
}

//金币的限制 todo 需要移动到配置文件
func (r *CMjRoom) GetCoinLimit() int64 {
	retcoin := int64(0)
	switch r.GetRoomLevel() {
	case 1:
		retcoin = 1000
	case 2:
		retcoin = 5000
	case 3:
		retcoin = 10000
	case 4:
		retcoin = 10000
	}

	return retcoin
}
