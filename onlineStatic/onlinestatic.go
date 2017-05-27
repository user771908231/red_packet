package main

import (
	"casino_admin/conf"
	"casino_admin/conf/config"
	"casino_common/common/sys"
	"casino_common/proto/ddproto"
	"os"
	"fmt"
	"time"
	"casino_common/common/service/statisticsService"
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
)

func init() {
	//初始化系统
	err := sys.SysInit(
		int32(ddproto.COMMON_ENUM_RELEASETAG_R_PRO),
		conf.Server.ProdMod,
		conf.Server.RedisAddr,
		"test",
		conf.Server.LogPath,
		"super",
		conf.Server.MongoIp,
		config.SUPER_DBNAM,
		[]string{
			config.DBT_SUPER_LOGS,
			config.DB_USER_SEQ,
		})

	//判断初始化是否成功
	if err != nil {
		os.Exit(-1) //推出系统
	}

	//初始化pushService
	//pushService.PoolInit(conf.Server.HallTcpAddr)
}

func main() {
	go timer()
	timer1()
	//t := time.Now().Unix()
	//fmt.Println(t)
	//fmt.Println(time.Unix(t,0).String())
}
//在线人数统计
const ADMIN_ONLINE_COUNT string = "t_online_count"
func timer() {
	timer1 := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-timer1.C:
			Time := time.Now().Format("2006-01-02 15:04:05")
			a :=statisticsService.OnlineCountAll()
			db.C(ADMIN_ONLINE_COUNT).Insert(bson.M{
				"OnlineCount" : a,
				"Time" : Time,
			})
		}
	}
}

//注册人数统计
const ADMIN_USER_REG_COUNT string = "t_user_reg_count"

func timer1() {
	timer1 := time.NewTicker(1 * time.Hour * 24)
	for {
		select {
		case <-timer1.C:
		//当天的日期
			Time := time.Now().Format("2006-01-02")
			date1,_ := time.Parse("2006-01-02",Time)
		//前一天的日期
			yesTime := time.Now().AddDate(0,0,-1)
			logDay := yesTime.Format("2006-01-02")
			date2,_ := time.Parse("2006-01-02",logDay)

		//长沙
			count1,_ :=db.C(tableName.DBT_T_USER).Count(bson.M{
				"regtime" : bson.M{"$gte": date2.Unix(),"$lte" : date1.Unix()},
				"$or" :[]bson.M{bson.M{"channelid" : 31},bson.M{"channelid" : 32},bson.M{"channelid" : 33}},
			})

		//岳阳
			count2,_ :=db.C(tableName.DBT_T_USER).Count(bson.M{
				"regtime" : bson.M{"$gte": date2.Unix(),"$lte" : date1.Unix()},
				"$or" :[]bson.M{bson.M{"channelid" : 34},bson.M{"channelid" : 35}},
			})

		//四川
			count3,_ :=db.C(tableName.DBT_T_USER).Count(bson.M{
				"regtime" : bson.M{"$gte": date2.Unix(),"$lte" : date1.Unix()},
				"$or" :[]bson.M{bson.M{"channelid" : 1},bson.M{"channelid" : 2},bson.M{"channelid" : 3},bson.M{"channelid" : 11},bson.M{"channelid" : 12},bson.M{"channelid" : 41},bson.M{"channelid" : 21},bson.M{"channelid" : 22}},
			})

		//白山
			count4,_ :=db.C(tableName.DBT_T_USER).Count(bson.M{
				"regtime" : bson.M{"$gte": date2.Unix(),"$lte" : date1.Unix()},
				"$or" :[]bson.M{bson.M{"channelid" : 61},bson.M{"channelid" : 62}},
			})

			db.C(ADMIN_USER_REG_COUNT).Insert(bson.M{
				"RegCount" : count1,
				"Time" : logDay,
				"channelidAll" : "31,32,33",
			})
			db.C(ADMIN_USER_REG_COUNT).Insert(bson.M{
				"RegCount" : count2,
				"Time" : logDay,
				"channelidAll" : "34,35",
			})
			db.C(ADMIN_USER_REG_COUNT).Insert(bson.M{
				"RegCount" : count3,
				"Time" : logDay,
				"channelidAll" : "1,2,3,11,12,21,22,41",
			})
			db.C(ADMIN_USER_REG_COUNT).Insert(bson.M{
				"RegCount" : count4,
				"Time" : logDay,
				"channelidAll" : "61,62",
			})
			fmt.Println("长沙 :",count1)
			fmt.Println("岳阳 :",count2)
			fmt.Println("四川 :",count3)
			fmt.Println("白山 :",count4)
		}
	}
}
