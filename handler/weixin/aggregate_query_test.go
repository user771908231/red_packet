package weixin

import (
	//"casino_common/common/sys"
	//db_init "casino_common/common/db"
	"testing"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_admin/model/agentModel"
)

func init() {
	//db_init.InitMongoDb("192.168.199.200", 27017, "test", "id",[]string{})
	//sys.InitRedis("192.168.199.200:6379","test")
}

func TestAggregateQuery(t *testing.T) {
	now := time.Now()
	y,m,d := now.Date()
	week_one := time.Date(y,m,d-int(now.Weekday())+1,0,0,0,0,time.Local)
	t.Log(week_one)
	resp := struct {
		Sum int32
	}{}
	query := []bson.M{
		bson.M{"$match":bson.M{
			"userid": 10301,
			"agentid": 11755,
			"addtime": bson.M{
				"$gt": week_one,
				"$lt": now,
			},
		}},
		bson.M{"$group":bson.M{
			"_id": nil,
			"sum": bson.M{"$sum": "$num"},
		}},
	}
	db.C(tableName.DBT_AGENT_SALES_LOG).Pipe(query, &resp)
	t.Log(resp.Sum)
}

func TestThinWeekAgentUserSalesCount(t *testing.T) {
	count := agentModel.GetAgentThisWeekOneUserSalesCount(11755, 10301)
	t.Log(count)
}
