package noticeServer

import (
	"casino_server/mode"
	"strings"
	"majiang/utils/numUtils"
	"casino_server/conf/casinoConf"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/redis"
	"github.com/name5566/leaf/db/mongodb"
	"gopkg.in/mgo.v2/bson"
	"casino_server/common/log"
)


func getRedisKey(id int32) string{
	keyPre := "public_notice_key"
	idStr ,_ := numUtils.Int2String(id)
	return strings.Join([]string{keyPre,idStr},"_")
}


//通过notice的type 来查找notice
func GetNoticeByType(noticeType int32) *bbproto.Game_AckNotice{
	//1,先从redis中获取notice
	rediConn := data.Data{}
	rediConn.Open(casinoConf.REDIS_DB_NAME)
	defer rediConn.Close()

	key := getRedisKey(noticeType)
	result := &bbproto.Game_AckNotice{}

	rediConn.GetObj(key, result)
	if result == nil ||  result.Id == nil {
		//2,如果notice 不存在则从数据库中获取,并且保存在redis中
		//redis 中去来的值为空
		//从数据库中获取值
		c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
		if err != nil {
			result = nil
		}
		defer c.Close()
		s := c.Ref()
		defer c.UnRef(s)

		//从数据库中查询user
		notice := &mode.T_th_notice{}
		s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_TH_NOTICE).Find(bson.M{"NoticeType": noticeType}).One(notice)
		if notice.Id == 0 {
			log.T("在mongo中没有查询到user[%v].", noticeType)
			result = nil
		}else{
			//把从数据获得的结果填充到redis的model中
			result = tnotice2Rnotice(notice)
			if result!=nil {
				saveNotice2Redis(result)
			}
		}
	}
	//3,返回数据
	return result
}

func tnotice2Rnotice(notice *mode.T_th_notice) *bbproto.Game_AckNotice{
	result := &bbproto.Game_AckNotice{}
	result.Id = new(int32)
	result.NoticeType = new(int32)
	result.NoticeTitle = new(string)
	result.NoticeContent = new(string)
	result.NoticeMemo = new(string)

	*result.Id = notice.Id
	*result.NoticeContent = notice.NoticeContent
	*result.NoticeTitle = notice.NoticeTitle
	*result.NoticeMemo = notice.NoticeMemo
	*result.NoticeType = notice.NoticeType

	return result
}


//把公告的数据保存到redis中
func saveNotice2Redis(notice *bbproto.Game_AckNotice) error {
	rediConn := data.Data{}
	rediConn.Open(casinoConf.REDIS_DB_NAME)
	defer rediConn.Close()
	key := getRedisKey(notice.GetId())
	rediConn.SetObj(key,notice)
	return nil
}
