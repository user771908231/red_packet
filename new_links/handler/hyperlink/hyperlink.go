package hyperlink

import (
	"new_links/modules"
	"new_links/model/hyperlinkModel"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/log"
	"time"
)

func Indexhandler(ctx *modules.Context) {
	start := time.Now()
	defer func() {
		cost := time.Since(start)
		log.T("请求用时：%s", cost)
	}()
	host := ctx.Query("host")
	log.T("获取到的Url:%s",host)
	url := hyperlinkModel.GetGroup(host)
	res := bson.M{
		"code":1,
		"url":"https://m.baidu.com/",
	}
	if url != "" {
		res["code"] = 1
		res["url"] = url
	}
	ctx.JSON(200, res)
}



