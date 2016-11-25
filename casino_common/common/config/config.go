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

	return nil
}

func InitRedis() {
	data.InitRedis(conf.Server.RedisAddr)
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
