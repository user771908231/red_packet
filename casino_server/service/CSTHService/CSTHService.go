//锦标赛德州铺可的service
package CSTHService

import "casino_server/msg/bbprotogo"



//返回存放竞标赛列表的redis-key
func getMatchListRedisKey() string {
	return "game_match_list_redis_key"
}



//竞标赛列表主要存放在标 t_cs_th_record中,但是一般都是在redis中取,如果redis中没有再从数据库中取
func GetGameMatchList() *bbproto.Game_MatchList {


	return nil
}
