package room

import (
	"casino_server/msg/bbprotogo"
	"casino_server/utils/redisUtils"
	"casino_server/common/log"
	"casino_server/service/pokerService"
	"casino_server/common/Error"
)
////////////////////////////////////////////服务器 数据恢复相关/////////////////////////////////////
var RUNNING_DESKS = "running_desk_keys"
//这里保存正在游戏中的thdesk
func AddRunningDesk(t *ThDesk) {
	tk := getRedisThDeskKey(t.Id, t.GameNumber)
	var keys *bbproto.RUNNING_DESKKEYS
	data := redisUtils.GetObj(RUNNING_DESKS, &bbproto.RUNNING_DESKKEYS{})
	if data == nil {
		keys = bbproto.NewRUNNING_DESKKEYS()
	} else {
		keys = data.(*bbproto.RUNNING_DESKKEYS)
	}

	for _, key := range keys.Desks {
		if key == tk {
			return
		}
	}

	keys.Desks = append(keys.Desks, tk)
	redisUtils.SetObj(RUNNING_DESKS, keys)
}


//删除正在进行的游戏
//一局正常结束的时候,需要把这一局的desk,gameNumber 删除

func RmRunningDesk(t *ThDesk) {
	index := -1
	tk := getRedisThDeskKey(t.Id, t.GameNumber)
	var keys *bbproto.RUNNING_DESKKEYS
	data := redisUtils.GetObj(RUNNING_DESKS, &bbproto.RUNNING_DESKKEYS{})
	if data != nil {
		keys = bbproto.NewRUNNING_DESKKEYS()
		for i, key := range keys.Desks {
			if key == tk {
				index = i
				break
			}
		}
	}
	//删除
	keys.Desks = append(keys.Desks[:index], keys.Desks[index + 1:]...)
	redisUtils.SetObj(RUNNING_DESKS, keys)

}

//查询目前正在游戏中的desk,目的是用来恢复服务器当机之前的数据
func GetRunningDesk() *bbproto.RUNNING_DESKKEYS {
	var keys *bbproto.RUNNING_DESKKEYS
	data := redisUtils.GetObj(RUNNING_DESKS, &bbproto.RUNNING_DESKKEYS{})
	if data != nil {
		keys = data.(*bbproto.RUNNING_DESKKEYS)
		return keys
	} else {
		return nil
	}
}

//恢复游戏数据
func (r *ThGameRoom) Recovery() {
	defer Error.ErrorRecovery("recovery...")
	log.T("开始恢复服务器crash之前的数据...")
	//1,找到对应的key
	keys := GetRunningDesk()
	log.T("找到需要恢复的desk[%v]", keys)
	if keys != nil {
		//2,循环处理每个key
		for _, key := range keys.Desks {
			log.T("开始恢复desk[%v]的数据", key)
			//通过key在数据库中恢复thdesk
			redisThdesk := GetRedisThDeskByKey(key)
			desk := RedisDeskTransThdesk(redisThdesk)
			if desk != nil && !desk.IsOver() {
				for _, userId := range redisThdesk.UserIds {
					log.T("开始恢复desk[%v]的user[%v]", key, userId)
					user := RedisThuserTransThuser(GetRedisThUser(desk.Id, desk.GameNumber, userId))        //依次恢复user
					user.thCards = pokerService.GetTHPoker(user.HandCards, desk.PublicPai, 5)                //重新计算玩家的牌信息
					desk.AddThuserBean(user)        //把desk add 到room
					desk.RecoveryRun()        //重新开始游戏,
				}
			} else {
				log.T("没有找到desk【%v】的数据(desk==nil?[%v]),desk.IsOver()[%v],desk.IsStop()[%v],恢复失败", key, desk == nil, desk.IsOver(), desk.IsStop())
			}
			r.AddThDesk(desk)        //把thdesk 加到room中
		}
	}
	log.T("恢复服务器crash之前的数据...end")
}


//thdesk 恢复之后,开始run
func (t *ThDesk) RecoveryRun() {
	//找到当前押注的人,然后等待押注
	user := t.GetUserByUserId(t.BetUserNow)
	if user != nil && t.IsRun() {
		user.wait()
	} else {
		log.E("恢复thdesk失败,没有找到betUserNow【%v】", t.BetUserNow)
	}
}

