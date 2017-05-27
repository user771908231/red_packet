package majiang

import (
	"casino_common/common/log"
	mjproto        "casino_mj_changsha/msg/protogo"
	"github.com/golang/protobuf/proto"
	"casino_common/common/sessionService"
	"casino_common/proto/ddproto"
	"casino_common/common/Error"
	"fmt"
	"casino_common/utils"
)

//金币场的配置
var coinConfig struct {
	boardsCount      int32 //总局数
	capMax           int64 //顶翻
	cardsNum         int32 //牌的张数
	settlement       int32
	baseValue        int64 //低分
	ziMoRadio        int32
	huRadio          int32
	dianGangHuaRadio int32
	totalPlayCount   int32 //玩家人数
	fee              int64
	userCountLimit   int32
	fangCountLimit   int32
}

func init() {
	coinConfig.boardsCount = 1000
	coinConfig.capMax = 8
	coinConfig.cardsNum = 13
	coinConfig.settlement = 0
	coinConfig.baseValue = 100                                  //低分默认是
	coinConfig.ziMoRadio = int32(mjproto.MJOption_ZIMO_JIA_FAN) //默认是自摸加翻
	coinConfig.huRadio = 1
	coinConfig.dianGangHuaRadio = 1
	coinConfig.totalPlayCount = 4 //玩家人数
	coinConfig.fee = 0
	coinConfig.userCountLimit = 4
	coinConfig.fangCountLimit = 3

}

//创建金币场的desk
func (r *MjRoom) CreateCoinPlayDesk() *MjDesk {
	log.T("开始创建一个新的金币场desk...")

	//创建金币场，创建时，需要更新mjroom的配置来创建
	desk := r.CreateDesk( //创建金币场的房间
		0,                                               //房主
		ROOMTYPE_COINPLAY,                               //房间类型
		int32(mjproto.MJRoomType_roomType_xueZhanDaoDi), //麻将类型 血战到底、三人两房、四人两房、德阳麻将、倒倒胡、血流成河
		coinConfig.boardsCount,                          //总局数
		coinConfig.capMax,                               // 顶翻
		coinConfig.cardsNum,                             //牌的张数
		coinConfig.settlement,
		r.BaseValue, //低分
		coinConfig.ziMoRadio,
		nil,
		coinConfig.huRadio,
		coinConfig.dianGangHuaRadio, //电杠花
		coinConfig.totalPlayCount,
		r.EnterCoinFee,
		coinConfig.userCountLimit, //人数
		coinConfig.fangCountLimit, //房数(牌的花色数)
	)
	if desk == nil {
		log.E("创建 desk 的时候出错...")
	}

	desk.CoinLimit = proto.Int64(r.CoinLimit) //最低的金币限制
	desk.CoinLimitUL = proto.Int64(r.CoinLimitUL)

	//新创建的房间，玩家进入之后是可以准备的
	desk.SetStatus(MJDESK_STATUS_READY) //可以准备的状态

	return desk
}

//通过session 找到对应的金币场的游戏desk
func (r *MjRoom) GetAbleCoinDeskBySession(userId uint32) (*MjDesk, error) {
	var retDesk *MjDesk
	//得到金币场的session
	csession := sessionService.GetSession(userId, ROOMTYPE_COINPLAY)
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
		return nil, Error.NewError(int32(ddproto.COMMON_ENUM_ERROR_TYPE_ENTERCOINROOM_ERROR_EC_OTHER_LV_DESK_GAMING), fmt.Sprintf(errMsg, gameName, ""))
	}

	//
	if csession.GetRoomId() != r.GetRoomId() {
		log.T("玩家[%v]当前的游戏状态room(id==%v)没有在当前room(id==%v)中..返回错误", userId, csession.GetRoomId(), r.GetRoomId())
		return nil, Error.NewError(int32(ddproto.COMMON_ENUM_ERROR_TYPE_ENTERCOINROOM_ERROR_EC_OTHER_LV_DESK_GAMING),
			fmt.Sprintf(errMsg, gameName, MjroomManagerIns.GetFMJRoom().RoomName))
	}

	//进入房间的error
	retDesk = r.GetDeskByDeskId(csession.GetDeskId())
	if retDesk == nil {
		//玩家的session信息已经过时，直接删除
		log.T("通过玩家[%v]的session[%v]没有找到对应的desk，表示session过期，需要删除对应的session", userId, csession)
		sessionService.DelSessionByKey(csession.GetUserId(), csession.GetRoomType(), csession.GetGameId(), csession.GetDeskId()) //通过session没有找到desk，证明session失效,此时删除session
		return nil, nil
	}

	//判断桌子的等级是否一样
	log.T("开始判断找到的desk的level[%v] 和 room 的level[%v] 是否相等", retDesk.GetRoomLevel(), r.GetRoomLevel())
	if retDesk.GetRoomLevel() != r.GetRoomLevel() {
		return nil, Error.NewError(int32(ddproto.COMMON_ENUM_ERROR_TYPE_ENTERCOINROOM_ERROR_EC_OTHER_LV_DESK_GAMING),
			fmt.Sprintf(errMsg, gameName, MjroomManagerIns.GetFMJRoom().RoomName))
	}

	//返回最终的数据
	return retDesk, nil
}

//找到金币场可以玩的桌子
func (r *MjRoom) GetAbleCoinDesk(userId uint32) (*MjDesk, error) {
	var retDesk *MjDesk
	//1,通过session来找

	var err error
	retDesk, err = r.GetAbleCoinDeskBySession(userId)
	if err != nil {
		return nil, err
	}

	////返回从session中找到的desk
	if retDesk != nil {
		log.T("找玩家%v对应的到老的desk %v ,", userId, retDesk)
		return retDesk, nil
	}

	log.T("么有找到玩家%v 对应老的desk，现在开中找合适的desk。", userId)
	//2,session中没有找到原来的desk，重新分配桌子
	for _, d := range r.Desks {
		if d != nil {
			//log.T("%v 金币场 d.GetUserCount():%v d.GetUserCountLimit() ", d.DlogDes(), d.GetUserCount(), d.GetUserCountLimit())
			if !d.IsPlayerEnough() {
				retDesk = d
			}
		}
	}

	//创建一个金币场的桌子
	if retDesk == nil {
		retDesk = r.CreateCoinPlayDesk()
	}
	log.T("获取到的可以使用level[%v]的金币场的desk[%v].", r.GetRoomLevel(), retDesk.DlogDes())

	return retDesk, nil

}
