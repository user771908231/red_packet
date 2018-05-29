package manage

import (
	"casino_redpack/modules"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"casino_common/common/service/exchangeService"
	"time"
	"math"
	"casino_redpack/model/agentProModel"

	"casino_redpack/model/weixinModel"
	"encoding/json"
	"fmt"
	"casino_redpack/model/userModel"
)

//红包与实物兑换

//列表
func ExchangeListHandler(ctx *modules.Context) {
	page := ctx.QueryInt("page")
	status := ctx.QueryInt("status")
	if page == 0 {
		page = 1
	}
	if status == 0 {
		status = 1
	}

	query := bson.M{
		"$and": []bson.M{
			bson.M{"status": status,},
		},
	}

	start_time := ctx.Query("start")
	if start_time != "" {
		start,_ := time.Parse("2006-01-02", start_time)
		ctx.Data["start_time"] = start_time
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"requesttime": bson.M{"$gte": start},
		})
	}
	end_time := ctx.Query("end")
	if end_time != "" {
		end,_ := time.Parse("2006-01-02", end_time)
		ctx.Data["end_time"] = end_time
		query["$and"] = append(query["$and"].([]bson.M), bson.M{
			"requesttime": bson.M{"$lt": end.AddDate(0,0,1)},
		})
	}

	list := []*exchangeService.ExchangeRecord{}
	_,count := db.C(tableName.DBT_ADMIN_EXCHANGE_RECORD).Page(query, &list, "-requesttime", page, 10)
	ctx.Data["status"] = status
	ctx.Data["list"] = list
	ctx.Data["page"] = bson.M{
		"count":      count,
		"list_count": len(list),
		"limit":      10,
		"page":       page,
		"page_count": math.Ceil(float64(count) / float64(10)),
	}
	ctx.HTML(200, "admin/manage/exchange/index")
}

//切换状态
func ExchangeSwitchState(ctx *modules.Context) {
	id := ctx.Query("id")
	status := ctx.QueryInt("status")
	if len(id) != 24 || status <= 0 || status > 5 {
		ctx.Ajax(-1, "参数错误", nil)
		return
	}
	row := exchangeService.GetRecordById(id)
	if row == nil {
		ctx.Ajax(-2, "切换状态失败！", nil)
		return
	}
	new_status := exchangeService.ExchangeState(status)
	row.Status = new_status
	if row.Save() != nil {
		ctx.Ajax(-3, "切换状态失败！", nil)
		return
	}
	ctx.Ajax(1, "切换状态成功！", nil)
}

func PostalHandle(ctx *modules.Context)  {
	status := ctx.QueryInt("status")
	page := ctx.QueryInt("page")
	fmt.Println("status:",status,"page:",page)
	switch status {
	case 0:
		query := bson.M{}
		data := agentProModel.GetOrderLists(query,page)
		ctx.Data["Postal"] = data
	default:
	}

	ctx.HTML(200, "admin/block/postal/index")
}
//提现申请
func WithdrawalsHandle(ctx *modules.Context)  {

	status := ctx.QueryInt("status")
	pages := ctx.QueryInt("page")
	switch status {
	case 1:
		query := bson.M{
			"acceptanceid" : 0,
			"status":0,
		}
		data := agentProModel.GetWithdrawalsList(query,pages)
		ctx.Data["withdrawals"] = data


	case 2:
		query := bson.M{
			"acceptanceid" : bson.M{"$gt":0},
			"status":2,
		}
		data := agentProModel.GetWithdrawalsList(query,pages)
		ctx.Data["withdrawals"] = data

	case 3:
		query := bson.M{
			"acceptanceid" :bson.M{"$gt":0},
			"status":1,
		}
		data := agentProModel.GetWithdrawalsList(query,pages)
		ctx.Data["withdrawals"] = data

	default:
		query := bson.M{
			"deletestatus" : 0,
			"status":bson.M{
				"$gte" : 0,
				"$lte" : 2,
			},
		}
		data := agentProModel.GetWithdrawalsList(query,pages)
		ctx.Data["withdrawals"] = data

	}
	ctx.HTML(200, "admin/block/withdrawals/index")
}
//后台申请提现方法
func WithdrawalsOperationHandle(ctx *modules.Context){
	Types := ctx.Query("types")
	Id := ctx.Query("id")
	list := bson.M{
		"code":0,
		"massage":"faild",
		"msg":"错误！",
	}
	switch Types {
	case "ok":
		fmt.Println("ok" )
		val := weixinModel.GetWithdrawalsId(bson.ObjectIdHex(Id))
		if val == nil{
			data,_ := json.Marshal(list)
			ctx.Write([]byte(data))
			return
		}
		////减去用户金币方法 准备
		//weixinModel.GetReady(val)
		//修改申请状态
		err := val.UpdateStatus(1,ctx.IsLogin().Id)
		if err != nil {
			list["msg"] ="修改失败！"
			data,_ := json.Marshal(list)
			ctx.Write([]byte(data))
			return
		}
		//结束
		//weixinModel.Implement("申请下分")
		list["msg"] ="修改成功！"
		list["code"] = 1
		list["massage"] = "success"
		data,_ := json.Marshal(list)
		ctx.Write([]byte(data))
		return
	case "no":
		fmt.Println("no" )
		val := weixinModel.GetWithdrawalsId(bson.ObjectIdHex(Id))
		if val == nil{
			data,_ := json.Marshal(list)
			ctx.Write([]byte(data))
		}
		user := userModel.GetUserById(val.UserId)
		user.CapitalUplete("+",val.Number," 提现拒绝")
		err := val.UpdateStatus(2,ctx.IsLogin().Id)
		if err != nil {
			list["msg"] ="修改失败！"
			data,_ := json.Marshal(list)
			ctx.Write([]byte(data))
			return
		}
		list["msg"] ="修改成功！"
		list["code"] = 1
		list["massage"] = "success"
		data,_ := json.Marshal(list)
		ctx.Write([]byte(data))
		return
	case "del":
		val := weixinModel.GetWithdrawalsId(bson.ObjectIdHex(Id))
		if val == nil{
			data,_ := json.Marshal(list)
			ctx.Write([]byte(data))
			return
		}
		err := val.Delete(1,ctx.IsLogin().Id)
		if err != nil {
			list["msg"] ="删除失败！"
			data,_ := json.Marshal(list)
			ctx.Write([]byte(data))
			return
		}
		list["msg"] ="删除成功！"
		list["code"] = 1
		list["massage"] = "success"
		data,_ := json.Marshal(list)
		ctx.Write([]byte(data))
		return
	default:
		list["msg"] =fmt.Sprint("意料之外的参数[types:%d]",Types)
		data,_ := json.Marshal(list)
		ctx.Write([]byte(data))
		return
	}

}
