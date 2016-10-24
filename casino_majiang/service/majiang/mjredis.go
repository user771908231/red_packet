package majiang

import (
	"casino_server/utils/redisUtils"
	"casino_server/common/log"
	"casino_server/utils/numUtils"
	"strings"
)

var REDIS_KEY_MJ_DESK = "redis_key_mj_desk"
var REDIS_KEY_MJ_RUNNING = "redis_key_mj_running"        //运行中的麻将


func GetDeskRedisKey(id int32) string {
	idStr, _ := numUtils.Int2String(id)
	return strings.Join([]string{REDIS_KEY_MJ_DESK, idStr}, "_")

}

//麻将redis 缓存相关的 业务逻辑代码都在这里....
func UpdateMjDeskRedis(desk *MjDesk) error {

	return nil
}


//通过id得到redis中的mjdesk
func GetDeskReids(key int32) *MjDesk {
	data := redisUtils.GetObj(GetDeskRedisKey(key), &MjDesk{})
	if data == nil {
		return nil
	}

	//得到的数据不为nil，则返回对应的mjdesk数据
	desk := data.(*MjDesk)
	return desk
}

//删除redis中麻将的数据 应该通过游戏编号来删除????
func DelMjDeskRedis(desk *MjDesk) error {

	return nil
}


//恢复麻将的数据
func RecoverFMJ() error {
	//1,得到运行中的keys
	keys := redisUtils.GetObj(REDIS_KEY_MJ_RUNNING, &RunningDeskKeys{})
	if keys == nil {
		log.T("没有需要回复的desk的数据")
		return nil
	}

	log.T("找到需要回复的麻将的数据keys[%v]", keys)
	runningKeys := keys.(*RunningDeskKeys)
	//开始回复数据
	for _, key := range runningKeys.Keys {
		log.T("开始回复key为[%v]的desk..", key)
		deskBak := GetDeskReids(key)
		if deskBak == nil {
			log.E("没有找到key[%v]对应的redis中的数据...继续下一个任务...")
			continue
		}

		//desk 不为空,为desk恢复锁

		//把deskBak 增加到room中去...
		FMJRoomIns.AddDesk(deskBak)
	}

	return nil
}




