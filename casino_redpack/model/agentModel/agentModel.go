package agentModel

import (
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"casino_common/common/log"
	"math"
	"time"
)

type AgentRebateLog struct {
	Id 			bson.ObjectId	`bson:"_id"`
	AgentId		uint32
	RebateId	uint32
	RebateMoeny	float64
}
type AgentRebateLog2 struct {
	Id 			bson.ObjectId	`bson:"_id"`
	AgentId		uint32
	RebateId	uint32
	RebateName string
	RebateMoeny	float64
	Time  time.Time
}


func (a *AgentRebateLog) Insert() error{
	a.Id = bson.NewObjectId()
	err := db.C(tableName.TABLE_REDPACK_AGENT_REBATE_LOG).Insert(a)
	return err
}
func (a *AgentRebateLog2) Insert() error{
	a.Id = bson.NewObjectId()
	err := db.C(tableName.TABLE_REDPACK_AGENT_REBATE_LOG).Insert(a)
	return err
}

func (a *AgentRebateLog) Inc(val float64) error{
	err := db.C(tableName.TABLE_REDPACK_AGENT_REBATE_LOG).Update(bson.M{"_id":a.Id},bson.M{"$inc":bson.M{"rebatemoeny":val}})
	return err
}

func (a *AgentRebateLog2) Inc(val float64) error{
	err := db.C(tableName.TABLE_REDPACK_AGENT_REBATE_LOG).Update(bson.M{"_id":a.Id},bson.M{"$inc":bson.M{"rebatemoeny":val}})
	return err
}

func GetAgentRebateLogPage(query interface{},page int,limit int) (int,[]*AgentRebateLog){

	list := []*AgentRebateLog{}
	_,count :=db.C(tableName.TABLE_REDPACK_AGENT_REBATE_LOG).Page(query, &list, "", page, limit)
	return count,list
}

func GetAgentRebateLogById(AgentId uint32,RebateId uint32) *AgentRebateLog{
	row := new(AgentRebateLog)
	err :=db.C(tableName.TABLE_REDPACK_AGENT_REBATE_LOG).Find(bson.M{"agentid":AgentId,"rebateid":RebateId},row)
	if err != nil {
		return nil
	}
	return row
}
func FloatValue(f float64,n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}
func GetAgentRebateLog(AgentId uint32,RebateId	uint32,moeny float64) error {
	A := GetAgentRebateLogById(AgentId,RebateId)
	//B := GetAgentRebateLogById(AgentId,RebateId)
	log.T("money:%f",moeny)
	if A != nil {
		log.T("找到了",A)
		moeny := FloatValue(moeny * 0.7,2)
		log.T("money:%f",moeny)
		err := A.Inc(moeny)
		return err
	}else {
		log.T("没找到")
		Agent := new(AgentRebateLog)
		Agent.AgentId = AgentId
		Agent.RebateId = RebateId
		Agent.RebateMoeny = moeny
		err := Agent.Insert()
		return err
	}

}
