package agentProModel

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
)

//代理申请记录
type AgentProRecordRow struct {
	ObjId bson.ObjectId  //id
	Name string  //姓名
	Telphone string  //电话
	Comment string  //留言
	Wxid string  //微信id
	Ip string  //当前ip
	AddTime  time.Time  //添加时间
}

//插入记录
func (t *AgentProRecordRow) Insert() error {
	return db.C(tableName.DBT_APPLY_AGENTPRO_RECORD).Insert(t)
}
