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

	//timer()
	timer2()
}
//在线人数统计
const ADMIN_ONLINE_COUNT string = "t_online_count"
func timer() {
	timer1 := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-timer1.C:
			a :=statisticsService.OnlineCountAll()
			fmt.Println(a)
			fmt.Println(time.Now().Unix())
			db.C(ADMIN_ONLINE_COUNT).Insert(bson.M{
				"OnlineCount" : a,
				"Time" : time.Now().Unix(),
			})
		}
	}
}
//注册人数统计
func timer2() {
	timer1 := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer1.C:
			Time := time.Now().Format("2006-01-02")
			date1,_ := time.Parse("2006-01-02",Time)
			count,_ :=db.C(tableName.DBT_T_USER).Count(bson.M{
				"regtime" : bson.M{"$gte": date1.Unix()},
			})
			nTime := time.Now()
			yesTime := nTime.AddDate(0,0,-1)
			logDay := yesTime.Format("2006-01-02")
			date2,_ := time.Parse("2006-01-02",logDay)
			fmt.Println(count)
			fmt.Println(date1.Unix())
			fmt.Println(date2)
		}
	}
}