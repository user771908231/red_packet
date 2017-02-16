package skeleton

import (
	"casino_common/common/log"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/protogo"
	"casino_majiang/service/majiang"
	"casino_majiang/msg/funcsInit"
)

//常见的get set 方法 需要放置在这里
func (f *SkeletonMJDesk) GetMJConfig() *data.SkeletonMJConfig {
	log.Debug("玩家[%v]进入fdesk")
	return f.config
}

//得到麻将的Status
func (r *SkeletonMJDesk) GetStatus() *data.MjDeskStatus {
	return r.status
}

//todo 日志信息
func (r *SkeletonMJDesk) DlogDes() string {
	return "todo"
}

//todo 通过userId 找到对应的User
func (r *SkeletonMJDesk) GetUserByUserId(userId uint32) api.MjUser {
	return nil
}

//todo 广播
func (d *SkeletonMJDesk) BroadCastProto(p proto.Message) {

}

//todo 是否有牌可以摸
func (d *SkeletonMJDesk) HandPaiCanMo() bool {
	return false
}
func (d *SkeletonMJDesk) GetCheckCase() *data.CheckCase {
	return d.CheckCase
}

// todo 游戏中玩家的人数
func (d *SkeletonMJDesk) GetGamingCount() int32 {
	return 0
}

func (d *SkeletonMJDesk) GetUsers() []api.MjUser {
	return d.Users
}

func (d *SkeletonMJDesk) GetBankerUser() api.MjUser {
	return d.GetUserByUserId(d.GetMJConfig().Banker)
}

//是否需要自摸加底
func (d *SkeletonMJDesk) IsNeedZiMoJiaDi() bool {
	if mjproto.MJOption(d.GetMJConfig().ZiMoRadio) == mjproto.MJOption_ZIMO_JIA_DI {
		return true
	}
	return false
}

//是否需要自摸加番
func (d *SkeletonMJDesk) IsNeedZiMoJiaFan() bool {
	if mjproto.MJOption(d.GetMJConfig().ZiMoRadio) == mjproto.MJOption_ZIMO_JIA_FAN {
		return true
	}
	return false
}

//判断是否开启房间的某个选
func (d *SkeletonMJDesk) IsOpenOption(option mjproto.MJOption) bool {
	for _, opt := range d.GetMJConfig().OthersCheckBox {
		if opt == int32(option) {
			return true
		}
	}
	return false
}

//是否需要换三张
func (d *SkeletonMJDesk) IsNeedExchange3zhang() bool {
	return d.IsOpenOption(mjproto.MJOption_EXCHANGE_CARDS)
}

//是否需要天地胡
func (d *SkeletonMJDesk) IsNeedTianDiHu() bool {
	return d.IsOpenOption(mjproto.MJOption_TIAN_DI_HU)
}

//是否需要幺九将对
func (d *SkeletonMJDesk) IsNeedYaojiuJiangdui() bool {
	return d.IsOpenOption(mjproto.MJOption_YAOJIU_JIANGDUI)
}

//是否需要门清中张
func (d *SkeletonMJDesk) IsNeedMenqingZhongzhang() bool {
	return d.IsOpenOption(mjproto.MJOption_MENQING_MID_CARD)
}

//判断是否是倒倒胡
func (d *SkeletonMJDesk) IsDaodaohu() bool {
	//倒倒胡，长沙麻将默认为倒倒胡
	if mjproto.MJRoomType(d.GetMJConfig().MjRoomType) == mjproto.MJRoomType_roomType_daoDaoHu ||
		mjproto.MJRoomType(d.GetMJConfig().MjRoomType) == mjproto.MJRoomType_roomType_changSha {
		return true
	}
	return false
}

//判断下一个庄是否已经确定
func (d *SkeletonMJDesk) IsNextBankerExist() bool {
	if d.GetMJConfig().NextBanker > 0 {
		return true
	} else {
		return false
	}
}

//设置下一个庄
func (d *SkeletonMJDesk) SetNextBanker(userId uint32) {
	d.GetMJConfig().NextBanker = userId
}


//将int32数组转为paiType数组
func (d *SkeletonMJDesk) IntArry2PaiTypeEnum(ia []int32) []mjproto.PaiType {
	var result []mjproto.PaiType
	for _, i := range ia {
		result = append(result, mjproto.PaiType(i))
	}
	return result
}

//得到当前桌子的人数..
func (d *SkeletonMJDesk) GetUserCount() int32 {
	var count int32 = 0
	for _, user := range d.GetUsers() {
		if user != nil {
			count ++
		}
	}
	//log.T("当前桌子的玩家数量是count[%v]", count)
	return count
}

//玩家是否足够
func (d *SkeletonMJDesk) IsPlayerEnough() bool {
	return d.GetUserCount() == d.GetMJConfig().PlayerCountLimit
}

//是不是所有人都准备
func (d *SkeletonMJDesk) IsAllReady() bool {
	for _, u := range d.Users {
		if u != nil && !u.GetStatus().IsReady() {
			return false
		}
	}
	return true
}

//是否在定缺中
func (d *SkeletonMJDesk) IsDingQue() bool {
	if d.GetStatus().S() == majiang.MJDESK_STATUS_DINGQUE {
		return true
	} else {
		return false
	}
}

func (d *SkeletonMJDesk) IsNotDingQue() bool {
	return !d.IsDingQue()
}

//得到下一张牌...
func (d *SkeletonMJDesk) GetNextPai() *majiang.MJPai {
	d.GetMJConfig().MJPaiCursor ++
	if d.GetMJConfig().MJPaiCursor >= 36*d.GetMJConfig().FangCount {
		log.E("服务器错误:要找的牌的坐标[%v]已经超过整副麻将的坐标了... ", d.GetMJConfig().MJPaiCursor)
		d.GetMJConfig().MJPaiCursor --
		return nil
	} else {
		p := d.AllMJPais[d.GetMJConfig().MJPaiCursor]
		pai := majiang.NewMjpai()
		*pai.Des = p.GetDes()
		*pai.Flower = p.GetFlower()
		*pai.Index = p.GetIndex()
		*pai.Value = p.GetValue()
		return pai
	}
}

func (d *SkeletonMJDesk) GetTotalMjPaiCount() int32 {
	return 36 * d.GetMJConfig().FangCount
}

func (d *SkeletonMJDesk) GetRemainPaiCount() int32 {
	remainCount := d.GetTotalMjPaiCount() - d.GetMJConfig().MJPaiCursor - 1
	return remainCount
}

//可以把overturn放在一个地方,目前都是摸牌的时候在用
func (d *SkeletonMJDesk) GetMoPaiOverTurn(user api.MjUser, isOpen bool) *mjproto.Game_OverTurn {

	overTurn := newProto.NewGame_OverTurn()
	*overTurn.UserId = user.GetUserId()                 //这个是摸牌的，所以是广播...
	*overTurn.PaiCount = d.GetRemainPaiCount()          //桌子剩余多少牌
	*overTurn.ActType = majiang.OVER_TURN_ACTTYPE_MOPAI //摸牌
	*overTurn.Time = 30
	if isOpen {
		overTurn.ActCard = user.GetGameData().HandPai.InPai.GetBackPai()
	} else {
		overTurn.ActCard = user.GetGameData().HandPai.InPai.GetCardInfo()
	}

	log.T("[%v]摸牌的时候牌:%v", d.DlogDes(), user.UserPai2String())
	*overTurn.CanHu, _, _, _, _, _ = d.HuParser.GetCanHu(user.GetGameData().GetHandPai(), user.GetGameData().GetHandPai().GetInPai(), true, 0) //是否可以胡牌
	*overTurn.CanPeng = false                                                                                                                  //是否可以碰牌

	//处理杠牌的时候
	/**
		1，血战到底：用户胡牌之后是不会进入到这个方法的
		2，血流成河：用户已经胡牌，那么杠牌之后，胡牌不会改变的情况下，才可以杠 // todo
	 */
	canGangBool, gangPais := user.GetGameData().HandPai.GetCanGang(nil, d.GetRemainPaiCount()) //是否可以杠牌
	*overTurn.CanGang = canGangBool
	if canGangBool && gangPais != nil {
		if user.GetStatus().IsHu() && d.IsXueLiuChengHe() {
			//血流成河，胡牌之后 杠牌的逻辑
			//jiaoPais := user.GetJiaoPaisByHandPais(); //得到杠牌之前的可以胡的叫牌
			jiaoPais := d.HuParser.GetJiaoPais(user.GetGameData().HandPai.Pais)
			for _, g := range gangPais {
				//判断杠牌之后的叫牌是否和杠牌之前一样
				if user.AfterGangEqualJiaoPai(jiaoPais, g) {
					overTurn.GangCards = append(overTurn.GangCards, g.GetCardInfo())
				}
			}
		} else {
			//没有胡牌之前，杠牌的逻辑....
			for _, g := range gangPais {
				overTurn.GangCards = append(overTurn.GangCards, g.GetCardInfo())
			}
		}
	}

	//最后判断是否可以杠牌
	if overTurn.GangCards == nil || len(overTurn.GangCards) <= 0 {
		*overTurn.CanGang = false
	}

	//最后判断是否需要增加过(可以杠，可以胡的时候需要增加可以过的按钮)
	if overTurn.GetCanGang() || overTurn.GetCanHu() {
		overTurn.CanGuo = proto.Bool(true)
	}

	//对长沙麻将做特殊处理
	overTurn.JiaoInfos = d.GetJiaoInfos(user)

	//这里需要对长沙麻将做特殊处理(主要是杠，补的处理)
	if d.IsChangShaMaJiang() {
		if overTurn.GetCanGang() {
			overTurn.CanBu = proto.Bool(true)
			overTurn.CanGang = proto.Bool(false)
			overTurn.BuCards = overTurn.GangCards
			overTurn.GangCards = nil
			//判断长沙麻将能不能杠
			for _, g := range overTurn.BuCards {
				cang := user.GetCanChangShaGang(InitMjPaiByIndex(int(g.GetId())))
				log.T("判断玩家[%v]对牌[%v]是否可以长沙杠[%v]", user.GetUserId(), g.GetId(), cang)
				if cang {
					overTurn.CanGang = proto.Bool(true)
					overTurn.GangCards = append(overTurn.GangCards, g)
				}
			}
		}
	}
	return overTurn
}

//通过checkCase 得到一个OverTurn
func (d *SkeletonMJDesk) GetOverTurnByCaseBean(checkPai *majiang.MJPai, caseBean *majiang.CheckBean, actType int32) *mjproto.Game_OverTurn {
	overTurn := newProto.NewGame_OverTurn()
	*overTurn.UserId = caseBean.GetUserId()
	*overTurn.CanGang = caseBean.GetCanGang()
	*overTurn.CanPeng = caseBean.GetCanPeng()
	*overTurn.CanHu = caseBean.GetCanHu()
	*overTurn.PaiCount = d.GetRemainPaiCount() //剩余多少钱
	overTurn.ActCard = checkPai.GetCardInfo()  //
	*overTurn.ActType = actType
	*overTurn.Time = 30
	overTurn.CanGuo = caseBean.CanGuo //目前默认是能过的
	overTurn.CanGuo = proto.Bool(true)
	overTurn.CanChi = caseBean.CanChi
	for i := 0; i < len(caseBean.ChiCards); i += 3 {
		c := &mjproto.ChiOverTurn{}
		c.ChiCard = append(c.ChiCard, caseBean.ChiCards[i].GetCardInfo())
		c.ChiCard = append(c.ChiCard, caseBean.ChiCards[i+1].GetCardInfo())
		c.ChiCard = append(c.ChiCard, caseBean.ChiCards[i+2].GetCardInfo())
		overTurn.ChiInfo = append(overTurn.ChiInfo, c)
	}

	return overTurn
}

func (d *SkeletonMJDesk) UpdateUserStatus(status int32) {
	for _, user := range d.GetUsers() {
		if user != nil {
			user.GetStatus().SetStatus(status)
		}
	}
}

//当前操作的玩家
func (d *SkeletonMJDesk) SetActUserAndType(userId uint32, actType int32) error {
	d.GetMJConfig().ActUser = userId
	d.GetMJConfig().ActType = actType
	return nil
}
