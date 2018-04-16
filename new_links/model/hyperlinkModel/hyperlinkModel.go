package hyperlinkModel

import (
	"new_links/model/groupingModel"
	"new_links/model/linksModel"
	"fmt"
	"new_links/model/keysModel"
	"math/rand"
	"github.com/chanxuehong/time"
)

func GetGroup(string string) string {

	//判断变量等于空
	if(string == ""){
		//随机给处一个分组名
		string = groupingModel.RandGetGroup()
	}
	G := groupingModel.GetGroupHost(string)
	val := linksModel.RandLink(G.ObjId)
	defer func() {
		if val != nil {
			val.Visssts()
		}
	}()
	if val == nil {
		return ""
	}
	return fmt.Sprintf(val.Url,val.Id,RandMap().Keys)
}

func RandMap() *keysModel.Keys {
	lengt := 0
	//读取初始化变量
	lengt = len(keysModel.Keyslist)
	//判断初始化变量的长度
	if lengt == 0 {
		//重新读取变量
		keysModel.UpdateInit()
		//获得变量的长度
		lengt = len(keysModel.Keyslist)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(lengt)
	return keysModel.Keyslist[i]
}