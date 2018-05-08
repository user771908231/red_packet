package paosangong

import (
	"github.com/name5566/leaf/gate"
	"casino_common/proto/ddproto"
)

//断线回调
func OnOffLine(gate gate.Agent) {
	agentData := gate.UserData()
	if agentData != nil {
		userId := agentData.(uint32)
		user, err := FindUserById(userId)
		if err == nil {
			//切为离线状态
			*user.IsOnline = false
			//发送离线广播
			user.SendOffineBc()
			//如果正在解散房间，则自动同意解散房间
			if user.Desk.GetIsOnDissolve() && user.GetDissolveState() == 0 {
				user.DoDissolveBack(true)
			}
			//自动进入托管模式
			if !user.GetIsTuoguan() && user.GetIndex() != -1 {
				user.DoTuoguan(true, &ddproto.NiuTuoguanOption{
					QiangZhuangOpt: ddproto.NiuEnumTuoguanQzopt_NIU_TG_QZ_BU_QIANG.Enum(),
					YaZhuOpt: ddproto.NiuEnumTuoguanYzopt_NIU_TG_YZ_YA_1.Enum(),
				})
			}
		}
	}
}
