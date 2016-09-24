package config

import (
	"casino_server/common/Error"
	"casino_server/common/cfg"
	"casino_server/common/log"
	"fmt"
	"os"
	"strconv"
	"casino_server/utils/redis"
	"casino_server/conf/casinoConf"
	"casino_server/utils/db"
	"casino_server/conf"
)



// 服务器初始化
func InitConfig(reloadFlg bool) (e Error.Error) {
	//初始化日志处理
	if reloadFlg {
		reloadLogger()
	} else {
		InitLogger()
		log.T("initlogger...")
	}

	//初始化redis连接配置
	fmt.Println("InitRedis...")
	InitRedis()

	//初始化mongoDb
	fmt.Println("InitMongoDb...")
	errInitMongoDb := InitMongoDb()
	if errInitMongoDb != nil {
		panic(errInitMongoDb)
	}
	return Error.OK()
}
/**
初始化数据库
1,建立自增主键
2,建立索引

 */
func InitMongoDb() error {
	//初始化地址...
	db.Oninit(conf.Server.MongoIp, conf.Server.MongoPort)
	//0,活的数据库连接
	c, err := db.GetMongoConn()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer c.Close()

	//1,建立自增主键,t_user表
	err = c.EnsureCounter(casinoConf.DB_NAME, casinoConf.DBT_T_USER, casinoConf.DB_ENSURECOUNTER_KEY)
	if err != nil {
		return err
	}

	//,建立自增主键,t_zjh_round 数据库表
	err = c.EnsureCounter(casinoConf.DB_NAME, casinoConf.DBT_T_ZJH_ROUND, casinoConf.DB_ENSURECOUNTER_KEY)
	if err != nil {
		return err
	}

	//为转盘奖励设置自增主键
	err = c.EnsureCounter(casinoConf.DB_NAME, casinoConf.DBT_T_BONUS_TURNTABLE, casinoConf.DB_ENSURECOUNTER_KEY)
	if err != nil {
		return err
	}

	//游戏编号的虚列号
	err = c.EnsureCounter(casinoConf.DB_NAME, casinoConf.DBT_T_TH_GAMENUMBER_SEQ, casinoConf.DB_ENSURECOUNTER_KEY)
	if err != nil {
		return err
	}

	//为德州桌子 thdesk
	err = c.EnsureCounter(casinoConf.DB_NAME, casinoConf.DBT_T_TH_DESK, casinoConf.DB_ENSURECOUNTER_KEY)
	if err != nil {
		return err
	}
	//德州比赛记录
	err = c.EnsureCounter(casinoConf.DB_NAME, casinoConf.DBT_T_TH_RECORD, casinoConf.DB_ENSURECOUNTER_KEY)
	if err != nil {
		return err
	}

	//为钻石交易记录创建自增键
	err = c.EnsureCounter(casinoConf.DB_NAME, casinoConf.DBT_T_USER_DIAMOND_DETAILS, casinoConf.DB_ENSURECOUNTER_KEY)
	if err != nil {
		return err
	}

	//为钻石交易记录创建自增键
	err = c.EnsureCounter(casinoConf.DB_NAME, casinoConf.DBT_T_CS_TH_RECORD, casinoConf.DB_ENSURECOUNTER_KEY)
	if err != nil {
		return err
	}

	return nil
}

func InitRedis() {
	data.InitRedis()
}

func InitLogger() {

	//配置文件:conf.ini
	config := cfg.Get()

	logMaxFile := 20
	if len(config["log_max_file"]) != 0 {
		logMaxFile, _ = strconv.Atoi(config["log_max_file"])
	}
	log.SetMaxFileCount(int32(logMaxFile))

	logMaxSize := 50
	if len(config["log_max_size"]) != 0 {
		logMaxSize, _ = strconv.Atoi(config["log_max_size"])
	}
	log.SetMaxFileSize(int64(logMaxSize), log.MB)

	logLevelStr := "ALL"
	if len(config["log_level"]) != 0 {
		logLevelStr = config["log_level"]
	}
	logLevel, ok := log.LevelMapping[logLevelStr]
	if !ok {
		logLevel = log.ALL
	}
	log.SetLevel(logLevel)

	logOut := "both"
	if len(config["log_output"]) != 0 {
		logOut = config["log_output"]
	}
	switch {
	case logOut == "both" || logOut == "BOTH" || logOut == "Both":
		log.SetConsole(true)
	case logOut == "file" || logOut == "FILE" || logOut == "File":
		log.SetConsole(false)
	}

	logDebugS := "FALSE"
	logDebug := false
	if len(config["log_debug"]) != 0 {
		logDebugS = config["log_debug"]
	}
	if logDebugS == "TRUE" || logDebugS == "True" || logDebugS == "true" {
		logDebug = true
	}
	log.SetDebug(logDebug)

	logMaxSeq := 500
	if len(config["log_max_seq"]) != 0 {
		logMaxSeq, _ = strconv.Atoi(config["log_max_seq"])
	}
	log.SetMaxLogSeq(int64(logMaxSeq))

	logPath := os.Getenv("GOPATH") + "/log"
	if len(config["log_path"]) != 0 {
		logPath = config["log_path"]
	}

	logName := "bbsvr"
	if len(config["log_name"]) != 0 {
		logName = config["log_name"]
	}

	if !(logOut == "none" || logOut == "NONE" || logOut == "None") {
		log.InitLoggers(logPath, logName)
	}
}

func reloadLogger() {
	//配置文件:conf.ini
	config := cfg.Get()

	logMaxFile := 20
	if len(config["log_max_file"]) != 0 {
		logMaxFile, _ = strconv.Atoi(config["log_max_file"])
	}
	log.SetMaxFileCount(int32(logMaxFile))

	logMaxSize := 50
	if len(config["log_max_size"]) != 0 {
		logMaxSize, _ = strconv.Atoi(config["log_max_size"])
	}
	log.SetMaxFileSize(int64(logMaxSize), log.MB)

	logLevelStr := "ALL"
	if len(config["log_level"]) != 0 {
		logLevelStr = config["log_level"]
	}
	logLevel, ok := log.LevelMapping[logLevelStr]
	if !ok {
		logLevel = log.ALL
	}
	log.SetLevel(logLevel)

	logOut := "both"
	if len(config["log_output"]) != 0 {
		logOut = config["log_output"]
	}
	switch {
	case logOut == "both" || logOut == "BOTH" || logOut == "Both":
		log.SetConsole(true)
	case logOut == "file" || logOut == "FILE" || logOut == "File":
		log.SetConsole(false)
	}

	logDebugS := "FALSE"
	logDebug := false
	if len(config["log_debug"]) != 0 {
		logDebugS = config["log_debug"]
	}
	if logDebugS == "TRUE" || logDebugS == "True" || logDebugS == "true" {
		logDebug = true
	}
	log.SetDebug(logDebug)
}
