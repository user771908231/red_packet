package hyperlink

import (
	"new_links/modules"
	"new_links/model/hyperlinkModel"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

func Indexhandler(ctx *modules.Context) {
	host := ctx.Query("host")
	url := hyperlinkModel.GetGroup(host)
	res := bson.M{
		"code":0,
		"url":"",
	}
	defer func() {
		data,_ := json.Marshal(res)
		ctx.Write([]byte(data))
	}()
	if url != "" {
		res["code"] = 1
		res["url"] = url
	}
	//ctx.Redirect(url,302)

}


