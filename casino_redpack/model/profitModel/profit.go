package profitModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_common/utils/db"
	"casino_redpack/model/agentModel"
	"casino_common/common/consts/tableName"
	"casino_redpack/model/userModel"
	"fmt"
	"casino_common/common/log"
	"math"
)
//总的收益
type Profit struct {
	Id bson.ObjectId `bson:"_id"`
	UserId uint32 //玩家ID
	Name string
	ProfitMoney float64
}
//每天的收益
type ProfitDetailed struct {
	Id bson.ObjectId `bson:"_id"`
	UserId uint32 //玩家ID
	Name string
	ProfitMoney float64
	Time time.Time
}

func (P *Profit) Insert() error {
	P.Id = bson.NewObjectId()
	err := db.C(tableName.TABLE_REDPACK_PROFIT_LOG).Upsert(bson.M{"_id":P.Id},P)
	return err
}

func (P *ProfitDetailed) Insert() error {
	P.Id = bson.NewObjectId()
	P.Time = agentModel.TimeObject()
	err := db.C(tableName.TABLE_REDPACK_PROFIT_ONE_DAY_LOG).Upsert(bson.M{"_id":P.Id},P)
	return err
}

func (P *Profit) Inc(val float64) error {
	err := db.C(tableName.TABLE_REDPACK_PROFIT_LOG).Update(bson.M{"_id":P.Id},bson.M{"$inc":bson.M{"profitmoney":val}})
	return err
}

func (P *ProfitDetailed) Inc(val float64) error {
	err := db.C(tableName.TABLE_REDPACK_PROFIT_ONE_DAY_LOG).Update(bson.M{"_id":P.Id,"time":agentModel.TimeObject()},bson.M{"$inc":bson.M{"profitmoney":val}})
	return err
}

func GetUserIdSelectProfit(id uint32) *Profit {
	P := new(Profit)
	err := db.C(tableName.TABLE_REDPACK_PROFIT_LOG).Find(bson.M{"userid":id},P)
	if err != nil {
		return nil
	}
	return P
}

func GetSelectProfit(query bson.M,page int) bson.M {
	list := bson.M{
		"page":bson.M{},
		"data":[]bson.M{},
		"user":0,
		"zong":0,
	}
	data := []bson.M{}
	P := []*Profit{}
	_,count := db.C(tableName.TABLE_REDPACK_PROFIT_LOG).Page(query,&P,"-_id",page,10)
	var zong float64
	for _,item := range P {
		row := bson.M{
			"id":item.Id.Hex(),
			"user":item.UserId,
			"name":userModel.GetUsernicknameById(item.UserId),
			"money":fmt.Sprintf("%.2f",item.ProfitMoney),
		}
		zong += item.ProfitMoney
		data=append(data,row)
	}
	list["page"] = bson.M{
		"count":      count,
		"list_count": len(data),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}
	list["data"] = data
	list["zong"] = fmt.Sprintf("%.2f",zong)
	return list
}

func GetUserIdSelectProfitDetailed(id uint32) *ProfitDetailed {
	P := new(ProfitDetailed)
	err := db.C(tableName.TABLE_REDPACK_PROFIT_ONE_DAY_LOG).Find(bson.M{"userid":id,"time":agentModel.TimeObject()},P)
	if err != nil {
		return nil
	}
	return P
}

func GetUserIdSelectProfitDetailedAll(query bson.M,page int) bson.M {
	P := []*ProfitDetailed{}
	list := bson.M{
		"page":bson.M{},
		"data":[]bson.M{},
		"user":1,
		"zong":0,
	}
	data := []bson.M{}
	var zong float64
	_,count := db.C(tableName.TABLE_REDPACK_PROFIT_ONE_DAY_LOG).Page(query,&P,"-_id",page,10)
	for _,item := range P {
		row := bson.M{
			"id":item.Id.Hex(),
			"user":item.UserId,
			"name":item.Name,
			"money":fmt.Sprintf("%.2f",item.ProfitMoney),
			"time":item.Time.Format("2006-01-02"),
		}
		zong += item.ProfitMoney
		data=append(data,row)
	}
	list["page"] = bson.M{
		"count":      count,
		"list_count": len(data),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}
	list["data"] = data
	list["zong"] = fmt.Sprintf("%.2f",zong)
	return list
}

func IncProfitLog(U *userModel.User,M float64) error{
	var err error = nil
	P := GetUserIdSelectProfit(U.Id)
	if P != nil {
		err = P.Inc(M)
		return err
	}
	P = new(Profit)
	P.UserId = U.Id
	P.Name	= IsAccountNumber(U)
	P.ProfitMoney = M
	err = P.Insert()
	return err
}

func IsAccountNumber(U *userModel.User) string {
	if U.AccountNumber != "" {
		return U.AccountNumber
	}
	return U.NickName
}

func IncProfitDetailedLog(U *userModel.User,M float64) error{
	var err error = nil
	P := GetUserIdSelectProfitDetailed(U.Id)
	if P != nil {
		err = P.Inc(M)
		return err
	}
	P = new(ProfitDetailed)
	P.UserId = U.Id
	P.Name	= IsAccountNumber(U)
	P.ProfitMoney = M
	err = P.Insert()
	return err
}

func ProfitLog(U *userModel.User,val float64) error {
	var err error = nil
	errP1 := IncProfitLog(U,val)
	errP2 := IncProfitDetailedLog(U,val)
	if errP1 != nil {
		log.T("收益总记录更新失败！ 错误原因：%s，用户：%s,值：%。2f",errP1,U,val)
		return errP1
	}
	if errP2 != nil {
		log.T("收益每日记录更新失败！ 错误原因：%s，用户：%s,值：%。2f",errP2,U,val)
		return errP1
	}
	return err
}





