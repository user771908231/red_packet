package majiang

import (
	"strings"
	"sync"
	"casino_common/utils/numUtils"
	"casino_common/utils/redisUtils"
	"casino_common/common/log"
)

var REDIS_KEY_MJ_DESK = "redis_key_mj_desk"
var REDIS_KEY_MJ_RUNNING = "redis_key_mj_running"        //运行中的麻将


//删除  和  增加 runningkey 的时候需要保存同步...

var runningKey_lock sync.Mutex = sync.Mutex{}

func GetDeskRedisKey(id int32) string {
	idStr, _ := numUtils.Int2String(id)
	return strings.Join([]string{REDIS_KEY_MJ_DESK, idStr}, "_")

}

//麻将redis 缓存相关的 业务逻辑代码都在这里....
func UpdateMjDeskRedis(desk *MjDesk) error {
	key := GetDeskRedisKey(desk.GetDeskId())
	redisUtils.SetObj(key, desk)
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

	//删除desk
	key := GetDeskRedisKey(desk.GetDeskId())
	redisUtils.Del(key)

	//删除running key
	DelRunningDeskKey(desk.GetDeskId())

	return nil
}

func GetRunningDeskeys() *RunningDeskKeys {
	keys := redisUtils.GetObj(REDIS_KEY_MJ_RUNNING, &RunningDeskKeys{})
	if keys == nil {
		log.T("没有需要回复的desk的数据")
		return nil
	} else {
		return keys.(*RunningDeskKeys)
	}
}

//删除对应的key
func DelRunningDeskKey(deskId int32) error {
	//删除和增加的时候 需要同步
	runningKey_lock.Lock()
	defer runningKey_lock.Unlock()

	runningKeys := GetRunningDeskeys()
	if runningKeys == nil {
		return nil
	}

	//删除key
	delIndex := -1
	for index, k := range runningKeys.Keys {
		if k == deskId {
			delIndex = index
		}
	}

	//删除对应的key
	if delIndex >= 0 {
		runningKeys.Keys = append(runningKeys.Keys[:delIndex], runningKeys.Keys[delIndex + 1:]...)
	}
	SaveRunningDeskKeys(runningKeys)
	return nil
}

//增加对应的key
func AddRunningDeskKey(deskId int32) error {
	//增加和删除的时候 需要同步...
	runningKey_lock.Lock()
	defer runningKey_lock.Unlock()

	runningKeys := GetRunningDeskeys()
	if runningKeys == nil {
		//如果runningkeys为nil新生成一个...
		runningKeys = new(RunningDeskKeys)
	}

	//如果已经存在，则不需要增加...
	for _, k := range runningKeys.Keys {
		if k == deskId {
			return nil
		}
	}

	runningKeys.Keys = append(runningKeys.Keys, deskId)
	SaveRunningDeskKeys(runningKeys)
	return nil
}

func SaveRunningDeskKeys(keys *RunningDeskKeys) {
	redisUtils.SetObj(REDIS_KEY_MJ_RUNNING, keys)
}


//恢复麻将的数据
func RecoverFMJ() error {
	//1,得到运行中的keys
	runningKeys := GetRunningDeskeys()
	if runningKeys == nil {
		log.T("没有需要回复的desk")
		return nil
	}

	//开始回复数据
	for _, key := range runningKeys.Keys {
		log.T("开始回复key为[%v]的desk..", key)
		deskBak := GetDeskReids(key)
		if deskBak == nil {
			log.E("没有找到key[%v]对应的redis中的数据...继续下一个任务...")
			continue
		}

		//desk 不为空,为desk恢复锁
		//把deskBak 增加到room中去...,增加的时候 生成lock
		FMJRoomIns.AddDesk(deskBak)
	}

	return nil
}




