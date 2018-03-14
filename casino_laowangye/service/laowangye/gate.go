package laowangye

import "github.com/name5566/leaf/gate"

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
			//if user.Desk.GetIsOnDissolve() && user.GetDissolveState() == 0 {
			//	user.DoDissolveBack(true)
			//}
		}
	}
}
