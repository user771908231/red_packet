package room

import (
	"casino_server/msg/bbprotogo"
	"strings"
	"casino_server/utils/numUtils"
	"casino_server/utils/redisUtils"
	"casino_server/mode/dao/TCsThRecordDao"
	"casino_server/service/userService"
)


//查询锦标赛排名相关的信息
var TOUNAMENTRANK_PRE string = "GetGame_TounamentRank_pre"

func getRounamentRankKey(matchId int32) string {
	matchIdStr, _ := numUtils.Int2String(matchId)
	return strings.Join([]string{TOUNAMENTRANK_PRE, matchIdStr}, "_")
}

//得到redis中的排名的数据
func GetGame_TounamentRank(matchId int32) *bbproto.Game_TounamentRank {

	//
	key := getRounamentRankKey(matchId)
	result := bbproto.NewGame_TounamentRank()

	//取值
	data := redisUtils.GetObj(key, &bbproto.Game_TounamentRank{})
	if data != nil {
		result = data.(*bbproto.Game_TounamentRank)
		return result
	} else {
		//数据库中查找
		record := TCsThRecordDao.GetByMatchId(matchId)
		//如果没有取到,则不做处理
		if record == nil {
			return nil
		}

		for i, rank := range record.Ranks {
			bean := bbproto.NewGame_TounamentRankBean()
			*bean.Coin = rank.WinCoin
			*bean.Place = int32(i + 1)
			*bean.PlayerName = userService.GetUserById(rank.UserId).GetNickName()
			*bean.PlayerImage = userService.GetUserById(rank.UserId).GetHeadUrl()
			result.Data = append(result.Data, bean)
		}

		redisUtils.SaveObj(key, result)
		return result
	}
}