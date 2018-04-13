package hyperlink

import (
	"new_links/modules"
	"new_links/model/hyperlinkModel"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

func Indexhandler(ctx *modules.Context) {
	host := ctx.Query("host")
	fmt.Println(host)
	url := hyperlinkModel.GetGroup(host)
	res := bson.M{
		"code":0,
		"url":"https://m.baidu.com/",
	}
	//defer func() {
	//	ctx.Resp.Header().Add()
	//	data,_ := json.Marshal(res)
	//	ctx.Write([]byte(data))
	//}()
	if url != "" {
		res["code"] = 1
		res["url"] = url
	}
	ctx.JSON(200, res)

}



