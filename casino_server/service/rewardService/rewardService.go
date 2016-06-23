package rewardService

import (
	"casino_server/msg/bbproto"
	"github.com/name5566/leaf/gate"
	"casino_server/conf/intCons"
)
/**
处理奖励的路由函数
呆改进,可以在这里获取UserId, 而不用在每个子函数中获取
 */
func HandlerRewards(m *bbproto.GetRewards,a gate.Agent) error{
	rewardsType := m.GetRewardsType()
	switch rewardsType {
	case intCons.REWARDS_TYPE_ONLINE:
		handlerRewardsOnline(m,a)
	case intCons.REWARDS_TYPE_RELIEF:
		handlerRewardsRelief(m,a)
	case intCons.REWARDS_TYPE_SIGNIN:
		handlerRewardsSignin(m,a)
	case intCons.REWARDS_TYPE_TIMING:
		handlerRewardsTiming(m,a)
	case intCons.REWARDS_TYPE_TURNTABLE:
		handlerRewardsTurntable(m,a)
	}

	return nil

}

//处理在线奖励
func handlerRewardsOnline(m *bbproto.GetRewards,a gate.Agent){

}

//处理签到奖励
func handlerRewardsSignin(m *bbproto.GetRewards,a gate.Agent){

}

//处理定时奖励
func handlerRewardsTiming(m *bbproto.GetRewards,a gate.Agent){

}

//处理救济金奖励
func handlerRewardsRelief(m *bbproto.GetRewards,a gate.Agent){

}

//处理转盘奖励
func handlerRewardsTurntable(m *bbproto.GetRewards,a gate.Agent){

}