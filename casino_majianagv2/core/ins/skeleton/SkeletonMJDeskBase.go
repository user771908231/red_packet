package skeleton

import (
	"casino_common/common/log"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/protogo"
	"casino_majiang/service/majiang"
	"fmt"
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

//日志信息
func (r *SkeletonMJDesk) DlogDes() string {
	s := fmt.Sprintf("[desk-%v-r-%v]", r.GetMJConfig().DeskId, r.GetMJConfig().CurrPlayCount)
	return s
}

//通过userId 找到对应的User
func (r *SkeletonMJDesk) GetUserByUserId(userId uint32) api.MjUser {
	for _, u := range r.GetUsers() {
		if u != nil && u.GetUserId() == userId {
			return u
		}
	}

	return nil
}

//广播
func (d *SkeletonMJDesk) BroadCastProto(p proto.Message) {
	for _, u := range d.Users {
		if u != nil {
			u.WriteMsg(p)
		}
	}
}

//是否有牌可以摸
func (d *SkeletonMJDesk) HandPaiCanMo() bool {
	if d.GetRemainPaiCount() == 0 {
		return false
	} else {
		return true
	}
}
func (d *SkeletonMJDesk) GetCheckCase() *data.CheckCase {
	return d.CheckCase
}

//游戏中玩家的人数
func (d *SkeletonMJDesk) GetGamingCount() int32 {
	var gamingCount int32 = 0 //正在游戏中的玩家数量
	for _, user := range d.GetUsers() {
		if user != nil && user.GetStatus().IsGaming() {
			gamingCount ++
		}
	}
	return gamingCount
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

//判断是否是血流成河
func (d *SkeletonMJDesk) IsXueLiuChengHe() bool {
	return d.GetMJConfig().MjRoomType == int32(mjproto.MJRoomType_roomType_xueLiuChengHe)
}

//返回desk 骨架
func (d *SkeletonMJDesk) GetSkeletonMJDesk() *SkeletonMJDesk {
	return d
}

func (d *SkeletonMJDesk) GetSkeletonMJUser(user api.MjUser) *SkeletonMJUser {
	return user.(*SkeletonMJUser)
}

//得到骨架User
func (d *SkeletonMJDesk) GetSkeletonMJUsers() []*SkeletonMJUser {
	ret := make([]*SkeletonMJUser, len(d.GetUsers()))
	for i, u := range d.GetUsers() {
		if u != nil {
			ret[i] = u.(*SkeletonMJUser)
		}
	}
	return ret
}

func (d *SkeletonMJDesk) BroadCastProtoExclusive(msg proto.Message, userId uint32) {
	for _, u := range d.Users {
		if u != nil && u.GetUserId() != userId {
			u.WriteMsg(msg)
		}
	}
}
