package room

import (
	"casino_server/msg/bbprotogo"
	"casino_server/utils/redisUtils"
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
		keys = bbproto.NewRUNNING_DESKKEYS()
		return keys
	} else {
		return nil
	}
}
