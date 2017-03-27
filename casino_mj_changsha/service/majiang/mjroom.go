package majiang

import (
	"errors"
	"github.com/name5566/leaf/gate"
	"casino_mj_changsha/msg/funcsInit"
	"casino_common/common/Error"
	"casino_common/common/consts"
	"casino_common/common/log"
	"casino_common/common/userService"
	"casino_common/utils/db"
	"casino_common/utils/chessUtils"
	"casino_common/proto/ddproto"
	"casino_common/common/consts/tableName"
	"github.com/name5566/leaf/module"
	mjproto	"casino_mj_changsha/msg/protogo"
	"github.com/golang/protobuf/proto"
	"casino_common/common/sessionService"
)

//普通的麻将房间...
const (
	ROOMTYPE_FRIEND   int32 = int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_FRIEND)
	ROOMTYPE_COINPLAY int32 = int32(ddproto.COMMON_ENUM_ROOMTYPE_DESK_COIN)
)

var (
	ERR_ROOMCARD_INSUFFICIENT error = Error.NewError(consts.ACK_RESULT_ERROR, "房卡不足")
)

type MjRoom struct {
	*PMjRoom               //room父类
	*module.Skeleton       //leaf  骨架
	Desks        []*MjDesk //所有的桌子
	CoinLimit    int64     // 金币准入限制
	CoinLimitUL  int64     // 金币准入限制
	BaseValue    int64     //底fen
	EnterCoinFee int64     //房费
	RoomName     string    //room的名字
}

//更具条件计算创建房间的费用

func (r *MjRoom) CalcCreateFee(boardsCout int32, mjtype int32) (int64, error) {

	log.T("麻将朋友桌开始计算房费:boardsCout %v mjtype %v ", boardsCout, mjtype)
	var fee int64 = 0

	//长沙麻将的配置
	if mjtype == int32(mjproto.MJRoomType_roomType_changSha) {
		if boardsCout == 8 {
			fee = 1
		} else if boardsCout == 16 {
			fee = 2
		}
		return fee, nil
	}

	//默认返回成都麻将的
	if r.IsFriend() {
		if boardsCout == 4 {
			fee = 2
		} else if boardsCout == 8 {
			fee = 3
		} else if boardsCout == 12 {
			fee = 5
		} else {
			return 0, errors.New("创建房间失败，局数有误")
		}
	}
	return fee, nil
}

//创建朋友桌desk
/**
	创建朋友桌 扣除房卡的逻辑..
	1,开始创建房间的时候需要判断房卡是否足够，但是并不会扣除，目的是在游戏没有开始就解散的时候，不用扣除房卡
	2,游戏开始之后，并且是第一局的时候扣除房卡

 */
func (r *MjRoom) CreateFriendDesk(
	userId uint32, mjRoomType int32,
	boardsCout int32, capMax int64,
	cardsNum int32,
	settlement int32, baseValue int64,
	ziMoRadio int32, othersCheckBox []int32,
	huRadio int32, dianGangHuaRadio int32, totalPlayCount int32, options *mjproto.ChangShaPlayOptions,
	userCountLimit int32, fangCountLimit int32) (*MjDesk, error) {

	//log.T("玩家[%v]创建朋友桌的房间...", userId)
	//找到是否有已经创建的房间
	oldDesk := MjroomManagerIns.GetFMjDeskBySession(userId)
	if oldDesk != nil && oldDesk.GetOwner() == userId && oldDesk.IsNotGaming() {
		//如果房间没有开始游戏..则返回老的房间,否则创建新的房间...
		return oldDesk, nil
	}

	//计算扣除的房卡数量
	createFee, err := r.CalcCreateFee(boardsCout, mjRoomType)
	if err != nil {
		log.E("玩家[%v]创建房间的时候出错..传入的局数[%v]有误...", userId, boardsCout)
		return nil, Error.ERR_SYS
	}

	rc := userService.GetUserRoomCard(userId) //创建朋友桌的时候，获取房主的房卡
	if rc < createFee {
		log.W("玩家[%v]创建房间的时候出错..房卡[%v]不足...", userId, rc)
		return nil, Error.NewError(consts.ACK_RESULT_ERROR, "房卡不足，创建房间失败")
	}

	//创建朋友桌的房间
	desk := r.CreateDesk(userId, ROOMTYPE_FRIEND, mjRoomType, boardsCout, capMax, cardsNum, settlement,
		baseValue, ziMoRadio, othersCheckBox, huRadio, dianGangHuaRadio, totalPlayCount, createFee, userCountLimit, fangCountLimit)
	if desk == nil {
		return nil, Error.ERR_SYS
	} else {
		desk.deductCreateFeeRoomCard()     //创建房间成功，扣除房卡
		desk.ChangShaPlayOptions = options //长沙麻将的配置
	}
	return desk, nil
}

//金币场 属性赋值
func (r *MjRoom) CreateDesk(userId uint32, roomType, mjRoomType int32, boardsCout int32, capMax int64, cardsNum int32,
	settlement int32, baseValue int64, ziMoRadio int32, othersCheckBox []int32, huRadio int32, dianGangHuaRadio int32, totalPlayCount int32, fee int64, userCountLimit int32, fangCountLimit int32) *MjDesk {
	//create 的时候，是否需要通过type 来判断,怎么样创建房间
	desk := NewMjDesk()
	desk.Skeleton = r.Skeleton
	desk.UserCountLimit = proto.Int32(userCountLimit)      //玩家的数量
	desk.Users = make([]*MjUser, desk.GetUserCountLimit()) //初始化玩家
	desk.FangCountLimit = proto.Int32(fangCountLimit)      //牌的花色的数量控制
	desk.AllUsers = desk.GetUsersApi
	*desk.DeskId, _ = db.GetNextSeq(tableName.DBT_MJ_DESK)
	*desk.RoomType = roomType
	*desk.RoomId = r.GetRoomId()
	desk.SetStatus(MJDESK_STATUS_READY) //设置为刚刚创建的状态
	*desk.MjRoomType = mjRoomType       // 房间类型，如：血战到底、三人两房、四人两房、德阳麻将、倒倒胡、血流成河，长沙麻将
	*desk.Password = r.RandRoomKey()
	*desk.Owner = userId          //设置房主
	*desk.CreateFee = fee         //创建房间的费用
	*desk.BoardsCout = boardsCout //局数，如：4局（房卡 × 2）、8局（房卡 × 3）
	*desk.CapMax = capMax
	*desk.CardsNum = cardsNum
	*desk.Settlement = settlement
	*desk.BaseValue = baseValue
	*desk.ZiMoRadio = ziMoRadio
	desk.OthersCheckBox = othersCheckBox
	*desk.HuRadio = huRadio
	*desk.DianGangHuaRadio = dianGangHuaRadio
	*desk.MJPaiCursor = 0
	*desk.TotalPlayCount = totalPlayCount
	*desk.CurrPlayCount = 0
	*desk.Banker = userId
	*desk.NextBanker = 0 //游戏过程中动态的计算
	desk.CheckCase = nil //判断case
	*desk.ActiveUser = userId
	*desk.GameNumber = 0
	*desk.ActUser = userId
	*desk.ActType = MJDESK_ACT_TYPE_MOPAI //初始化状态
	//*desk.BeginTime	//游戏开始时间...
	//*desk.EndTime		//游戏结束时间...
	*desk.NInitActionTime = 30 // 游戏操作的时间
	*desk.RoomLevel = r.GetRoomLevel()

	//把创建的desk加入到room中
	if desk.IsChangShaMaJiang() {
		//长沙麻将的解析器
		desk.HuParser = NewHuParserChangSha()
	} else {
		//默认使用成都麻将的解析器
		desk.HuParser = NewHuParserChengdu(
			desk.GetBaseValue(),
			desk.IsNeedZiMoJiaDi(),
			desk.IsNeedYaojiuJiangdui(),
			desk.IsDaodaohu(),
			desk.IsNeedMenqingZhongzhang(),
			desk.IsNeedZiMoJiaFan(),
			desk.GetRoomTypeInfo().GetCapMax())
	}
	r.AddDesk(desk)
	return desk
}

func (r *MjRoom) RandRoomKey() string {
	//金币场没有房间号码
	if r.GetRoomType() == ROOMTYPE_COINPLAY {
		return ""
	}
	roomKey := chessUtils.GetRoomPass(int32(ddproto.CommonEnumGame_GID_MAHJONG))
	//1,判断roomKey是否已经存在
	if r.IsRoomKeyExist(roomKey) {
		//log.E("房间密钥[%v]已经存在,创建房间失败,重新创建", roomKey)
		return r.RandRoomKey()
	} else {
		//log.T("最终得到的密钥是[%v]", roomKey)
		return roomKey
	}
	return ""
}

//判断roomkey是否已经存在了
func (r *MjRoom) IsRoomKeyExist(roomkey string) bool {
	ret := false
	for i := 0; i < len(r.Desks); i++ {
		d := r.Desks[i]
		if d != nil && d.GetPassword() == roomkey {
			ret = true
			break
		}
	}

	return ret
}

//通过房间号码得到desk
func (r *MjRoom) GetDeskByPassword(key string) *MjDesk {
	//如果找到对应的房间，则返回
	for _, d := range r.Desks {
		if d != nil && d.GetPassword() == key {
			return d
		}
	}

	//如果没有找到，则返回nil
	return nil
}

//通过房间号码得到desk
func (r *MjRoom) GetDeskByDeskId(id int32) *MjDesk {
	//log.T("通过deskId【%v】查询desk", id)
	//如果找到对应的房间，则返回
	for _, d := range r.Desks {
		if d != nil && d.GetDeskId() == id {
			//log.T("通过id[%v]找到desk----d.getDeskId()[%v]", id, d.GetDeskId())
			return d
		}
	}
	//如果没有找到，则返回nil
	return nil
}

//进入房间
//进入的时候，需要判断牌房间的类型...
func (r *MjRoom) EnterRoom(key string, userId uint32, a gate.Agent) error {
	var desk *MjDesk
	var err error
	//如果是朋友桌,需要通过房间好来找到desk
	if r.IsFriend() {
		if desk = r.GetDeskByPassword(key); desk == nil {
			//如果玩家的session中的房间号和key相同，那么证明session已经过期了，这个时候可以删除session
			//如果玩家的session中的房间号和key不同，那么证明玩家确实是把房间号输入错误了
			s := sessionService.GetSession(userId, ROOMTYPE_FRIEND)
			if s != nil &&
				s.GetGameId() == int32(ddproto.CommonEnumGame_GID_MAHJONG) &&
				s.GetRoomPassword() == key {
				//证明session过期，这个时候需要删除session
				sessionService.DelSessionByKey(s.GetUserId(), s.GetRoomType(), s.GetGameId(), s.GetDeskId())
				log.T("玩家的session过期，现在开始删除session[%v]", s)
			} else {
				log.T("玩家[%v]在麻将朋友桌输入了错误的房间号码[%v]，进入房间失败..", userId, key)
			}
			return Error.NewError(consts.ACK_RESULT_FAIL, "房间号输入错误")
		}
	} else if r.IsCoinPlay() {
		//金币场加入房间的逻辑
		if desk, err = r.GetAbleCoinDesk(userId); desk == nil {
			//删除玩家对应的session
			return err
		}
	}

	//进入房间的逻辑
	err = desk.enterUser(userId, a) //普通玩家进入房间
	return err
}

func (r *MjRoom) IsFriend() bool {
	return r.GetRoomType() == ROOMTYPE_FRIEND
}

//判断是是否是金币场
func (r *MjRoom) IsCoinPlay() bool {
	return r.GetRoomType() == ROOMTYPE_COINPLAY
}

func (r *MjRoom) AddDesk(desk *MjDesk) error {
	r.Desks = append(r.Desks, desk)
	//加入之后需要更新数据到redis
	desk.updateRedis()
	AddRunningDeskKey(desk.GetDeskId())
	return nil
}

// room 解散房间...解散朋友桌
func (r *MjRoom) DissolveDesk(desk *MjDesk, sendMsg bool) error {
	//清楚数据,1,session相关。2,
	log.T("%v开始解散...", desk.DlogDes())
	for _, user := range desk.GetUsers() {
		desk.RmUserSession(user)
	}

	log.T("开始删除desk[%v]...", desk.GetDeskId())

	//判断是否需要返回房卡
	err := desk.sendBackRoomCrad()
	if err != nil {
		log.E("%v 解散，退还房卡的时候...", desk.DlogDes())
	}

	//删除房间
	rmErr := r.RmDesk(desk)
	if rmErr != nil {
		log.E("删除房间失败,errmsg[%v]", rmErr)
		return rmErr
	}

	//删除reids
	DelMjDeskRedis(desk)

	//发送解散房间的广播
	log.T("删除desk[%v]之后，发送删除的广播...", desk.GetDeskId())
	if sendMsg {
		//发送解散房间的广播
		dissolve := newProto.NewGame_AckDissolveDesk()
		*dissolve.DeskId = desk.GetDeskId()
		*dissolve.PassWord = desk.GetPassword()
		*dissolve.UserId = desk.GetOwner()
		desk.BroadCastProto(dissolve)
	}

	return nil

}

func (r *MjRoom) RmDesk(desk *MjDesk) error {
	index := -1
	for i, d := range r.Desks {
		if d != nil && d.GetDeskId() == desk.GetDeskId() {
			index = i
			break
		}
	}

	if index >= 0 {
		r.Desks = append(r.Desks[:index], r.Desks[index+1:]...)
		return nil
	} else {
		return errors.New("删除失败，没有找到对应的desk")
	}

}
