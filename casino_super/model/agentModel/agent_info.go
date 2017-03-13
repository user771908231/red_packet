package agentModel

import (
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/consts/tableName"
	"casino_common/utils/db"
)

type AgentType int

const (
	AGENT_TYPE_1 AgentType = 1  //一级总代理
	AGENT_TYPE_2 AgentType = 2  //一级普通代理
	AGENT_TYPE_3 AgentType = 3  //普通代理
)

//代理信息表
type AgentInfo struct {
	Id       bson.ObjectId `bson:"_id"`
	UserId   uint32  //游戏内user_id
	NickName string  //昵称
	RealName string  //真实姓名
	Phone    string  //手机号
	OpenId   string
	UnionId  string
	RootId   uint32  //所属根代理id
	Pid      uint32  //所属上级代理id
	Level    int32  //代理级别：一级 二级 三级
	Type     AgentType  //代理类型： 总代 普通带理
}

//插入表
func (r *AgentInfo) Insert() error {
	r.Id = bson.NewObjectId()
	return db.C(tableName.DBT_AGENT_INFO).Insert(r)
}

//保存
func (r *AgentInfo) Save() error {
	return db.C(tableName.DBT_AGENT_INFO).Update(bson.M{
		"_id": r.Id,
	}, r)
}

//通过Agent_id获取agent_info
func GetAgentInfoById(agent_id uint32) *AgentInfo {
	info := new(AgentInfo)
	err := db.C(tableName.DBT_AGENT_INFO).Find(bson.M{
		"userid": agent_id,
	}, info)
	if err != nil {
		return nil
	}
	return info
}

//判断一个用户是否为代理商
func IsAgent(user_id uint32) bool {
	row := new(ApplyRecord)
	err := db.C(tableName.DBT_AGENT_INFO).Find(bson.M{
		"userid": user_id,
	}, row)
	if err == nil {
		return true
	}
	return false
}

//获取代理下面的子代理数量
func GetAgentChildNum(agent_id uint32) int {
	num, _ := db.C(tableName.DBT_AGENT_INFO).Count(bson.M{
		"pid": agent_id,
	})
	return num
}
