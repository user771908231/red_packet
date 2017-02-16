package skeleton

import (
	"casino_common/common/log"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
	"github.com/golang/protobuf/proto"
	"casino_majiang/msg/protogo"
	"casino_majiang/service/majiang"
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

func (d *SkeletonMJDesk) GetBankerUser() *api.MjUser {
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
	if d.GetStatus(). == majiang.MJDESK_STATUS_DINGQUE {
		return true
	} else {
		return false
	}
}

func (d *SkeletonMJDesk) IsNotDingQue() bool {
	return !d.IsDingQue()
}
