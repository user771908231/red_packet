package agentPro

import (
	"new_links/modules"
	"new_links/model/agentProModel"
	"time"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/log"
)

//申请处理
func ApplyPostHandler(ctx *modules.Context, form agentProModel.AgentProRecordRow) {
	form.ObjId = bson.NewObjectId()
	form.Ip = ctx.Req.RemoteAddr
	form.AddTime = time.Now()
	form.Insert()

	log.T("[apply_agentpro post] %v", form)
}
