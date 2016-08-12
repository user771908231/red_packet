package noticeServer

import (
	"casino_server/mode"
	"strings"
	"casino_server/conf/casinoConf"
	"casino_server/msg/bbprotogo"
	"casino_server/utils/redis"
	"github.com/name5566/leaf/db/mongodb"
	"gopkg.in/mgo.v2/bson"
	"casino_server/common/log"
	"errors"
	"casino_server/utils/numUtils"
)


var NOTICE_TYPE_GUNDONG  int32  = 1	//滚动
var NOTICE_TYPE_CHONGZHI int32  = 2	//充值信息
var NOTICE_TYPE_GONGGAO  int32  = 3	//公告信息

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

		s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_TH_NOTICE).Find(bson.M{"noticetype" : noticeType}).One(notice)
		if notice.Id == 0 {
			log.T("在mongo中没有查询到notice[%v].", noticeType)
			result = nil
		}else{
			log.T("数据库中查询到noticeType[%v]的 notice[%v]",noticeType,notice)
			//把从数据获得的结果填充到redis的model中
			result = tnotice2Rnotice(notice)
			if result!=nil {
				SaveNotice2Redis(result)
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
	result.Fileds  = notice.NoticeFileds

	return result
}


//把公告的数据保存到redis中
func SaveNotice2Redis(notice *bbproto.Game_AckNotice) error {
	rediConn := data.Data{}
	rediConn.Open(casinoConf.REDIS_DB_NAME)
	defer rediConn.Close()
	key := getRedisKey(notice.GetId())
	rediConn.SetObj(key,notice)
	return nil
}

//把公告的数据保存到redis中
func SaveNotice(tnotice *mode.T_th_notice) error {
	saveNotice2mongos(tnotice)	//保存到mongo
	rnotice := tnotice2Rnotice(tnotice)
	SaveNotice2Redis(rnotice)
	return nil
}

//保存公告到mongo数据库
func saveNotice2mongos(tnotice *mode.T_th_notice) error{
	log.T("开始保存公告信息到mogno,notice[%v]",tnotice)
	//从数据库中获取值
	c, err := mongodb.Dial(casinoConf.DB_IP, casinoConf.DB_PORT)
	if err != nil {
		return errors.New("连接数据库失败")
	}
	defer c.Close()
	s := c.Ref()
	defer c.UnRef(s)

	//从数据库中查询user
	s.DB(casinoConf.DB_NAME).C(casinoConf.DBT_T_TH_NOTICE).Insert(tnotice)
	return nil
}




