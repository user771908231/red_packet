package skeleton

import (
	"casino_majiang/service/majiang"
	"casino_majianagv2/core/data"
	"time"
	"casino_majianagv2/core/api"
	"fmt"
	"casino_majianagv2/core/majiangv2"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/gate"
)

type SkeletonMJUser struct {
	desk            api.MjDesk
	status          *data.MjUserStatus
	userId          uint32
	readyTimer      *time.Timer
	Bill            *majiang.Bill
	GameData        *data.MJUserGameData
	Coin            int64                  //金币
	Statisc         *majiang.MjUserStatisc //统计信息
	a               gate.Agent
	ActTimeoutCount int32 //
}

//初始化一个user骨架
func NewSkeleconMJUser(userId uint32) *SkeletonMJUser {
	return nil
}

func (user *SkeletonMJUser) Ready() {
	//设置为准备的状态,并且停止准备计时器
	user.status.SetStatus(majiang.MJUSER_STATUS_READY)
	user.status.Ready = true
	if user.readyTimer != nil {
		user.readyTimer.Stop()
		user.readyTimer = nil
	}

}

func (user *SkeletonMJUser) UserPai2String() string {
	result := "玩家[%v]牌的信息,handPais[%v],inpai[%v],pengpais[%v],gangpai[%v]"
	result = fmt.Sprintf(result, user.GetUserId(),
		majiangv2.ServerPais2string(user.GetGameData().HandPai.Pais), user.GetGameData().HandPai.InPai.LogDes(),
		majiangv2.ServerPais2string(user.GetGameData().HandPai.PengPais), majiangv2.ServerPais2string(user.GetGameData().HandPai.GangPais))
	return result
}

//比较杠牌之后的叫牌和杠牌之前的叫牌的信息是否一样
func (u *SkeletonMJUser) AfterGangEqualJiaoPai(beforJiaoPais []*majiang.MJPai, gangPai *majiang.MJPai) bool {

	//1，获得杠牌之后的手牌
	var afterPais []*majiang.MJPai
	for _, p := range u.GameData.HandPai.Pais {
		if p.GetClientId() != gangPai.GetClientId() {
			afterPais = append(afterPais, p)
		}
	}

	//2，通过杠牌之后的手牌 获得此时的叫牌
	afterJiaoPais := u.GetDesk().GetHuParser().GetJiaoPais(afterPais)

	//2,比较beforJiaoPais 和 afterJiaoPais
	if len(afterPais) != len(beforJiaoPais) {
		return false
	}

	for _, aj := range afterJiaoPais {

		forbool := false
		for _, bj := range beforJiaoPais {
			if aj.GetClientId() == bj.GetClientId() {
				forbool = true
				break
			}
		}

		if !forbool {
			return false
		}
	}

	return true;
}

//初始化方法
func (u *SkeletonMJUser) BeginInit(CurrPlayCount int32, banker uint32) {

}

//发送overTrun
/**
	这里需要区分有托管 和没有托管的状态：
	1，有托管的时候，给玩家发送
 */
func (u *SkeletonMJUser) SendOverTurn(p proto.Message) error {
	//如果是金币场有超时的处理...
	u.WriteMsg(p)
	return nil
}
