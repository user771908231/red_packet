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
		val.Visssts()
	}()
	if val == nil {
		return ""
	}
	return fmt.Sprintf(val.Url,val.Id,RandMap().Keys)
}

func RandMap() *keysModel.Keys {
	lengt := 0
	lengt = len(keysModel.Keyslist)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(lengt)
	return keysModel.Keyslist[i]
}