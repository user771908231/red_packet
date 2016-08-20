//锦标赛德州铺可的service
package CSTHService

import (
	"casino_server/msg/bbprotogo"
	"casino_server/utils/redisUtils"
	"casino_server/service/room"
	"casino_server/utils/db"
)


//返回存放竞标赛列表的redis-key
func GetMatchListRedisKey() string {
	return "game_match_list_redis_key"
}

//竞标赛列表主要存放在标 t_cs_th_record中,但是一般都是在redis中取,如果redis中没有再从数据库中取
func GetGameMatchList() *bbproto.Game_MatchList {
	data := redisUtils.GetObj(GetMatchListRedisKey(),&bbproto.Game_MatchList{})
	if data == nil{
		return data
	}else{
		return data.(*bbproto.Game_MatchList)
	}
}

//刷新锦标赛的列表信息
func RefreshRedisMatchList(r *room.CSThGameRoom) {
	//1,获取数据库中的近20场次的信息(通过时间来排序)
	db.UpdateMgoData()


}
