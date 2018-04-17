package hyperlinkModel

import (
	"new_links/model/groupingModel"
	"new_links/model/linksModel"
	"fmt"
	"new_links/model/keysModel"
	"math/rand"
	"github.com/chanxuehong/time"
	"casino_common/common/log"
)

func GetGroup(string string) string {

	//判断变量等于空
	if(string == ""){
		log.T("获取到的Url:为空")
		//随机给处一个分组名
		string = groupingModel.RandGetGroup()
	}
	G := groupingModel.GetGroupHost(string)
	if G != nil {
		log.T("获取到的Group:%s",string)
		val := linksModel.RandLink(G.ObjId)
		defer func() {
			if val != nil {
				val.Visssts()
			}
		}()
		if val == nil {
			log.T("获取到的随机链接:为空")
			return ""
		}
		key := RandMap().Keys
		log.T("获取到的随机链接ID:[%d] 关键词：%s",val.Id,key)
		return fmt.Sprintf(val.Url,val.Id,key)
	}
	return ""
}

func RandMap() *keysModel.Keys {
	lengt := 0
	//读取初始化变量
	lengt = len(keysModel.Keyslist)
	//判断初始化变量的长度
	if lengt == 0 {
		log.T("初始化变量的长度：为零")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(lengt)
	return keysModel.Keyslist[i]
}