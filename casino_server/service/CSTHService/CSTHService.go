//锦标赛德州铺可的service
package CSTHService

import (
	"casino_server/msg/bbprotogo"
	"casino_server/utils/redisUtils"
	"casino_server/utils/db"
	"casino_server/mode"
	"gopkg.in/mgo.v2"
	"casino_server/conf/casinoConf"
	"gopkg.in/mgo.v2/bson"
	"casino_server/common/log"
	"casino_server/utils/timeUtils"
)


//返回存放竞标赛列表的redis-key
func GetMatchListRedisKey() string {
	return "game_match_list_redis_key"
}

//竞标赛列表主要存放在标 t_cs_th_record中,但是一般都是在redis中取,如果redis中没有再从数据库中取
func GetGameMatchList() *bbproto.Game_MatchList {
	data := redisUtils.GetObj(GetMatchListRedisKey(), &bbproto.Game_MatchList{})
	if data == nil {
		log.T("redis中没有找到竞标赛列表,需要在数据库中查找")
		RefreshRedisMatchList()
		return nil
	} else {
		return data.(*bbproto.Game_MatchList)
	}
}

//刷新锦标赛的列表信息
func RefreshRedisMatchList() {
	//1,获取数据库中的近20场次的信息(通过时间来排序)
	data := []mode.T_cs_th_record{}
	db.Query(func(d *mgo.Database) {
		d.C(casinoConf.DBT_T_CS_TH_RECORD).Find(bson.M{}).Sort("-id").Limit(20).All(&data)
	})

	//把得到的数据保存在数据库中
	if len(data) > 0 {
		//表示有数据,需要存储在redis中
		sdata := &bbproto.Game_MatchList{}
		for i := 0; i < len(data); i++ {
			d := data[i]
			sd := bbproto.NewGame_MatchItem()
			*sd.CostFee = d.CostFee
			*sd.Title = "神经德州锦标赛"
			*sd.Status = d.Status
			*sd.Type = d.GameType
			*sd.Time = timeUtils.Format(d.BeginTime)
			sdata.Items = append(sdata.Items, sd)
		}

		//存储
		redisUtils.SaveObj(GetMatchListRedisKey(), sdata)
	}

}
