package config

import (
	"casino_server/common/Error"
	"casino_server/common/cfg"
	"casino_server/common/consts"
	"casino_server/common/log"
	"fmt"
	"math"
	"os"
	"strconv"
	"casino_server/utils/redis"
	"casino_server/conf/casinoConf"
	"casino_server/utils/db"
)

var (
	TableDevourCostCoin   []int32
	G_GAME_LANGUAGE string                          // 系统语言
	G_LoginPresentConf    *BonusConfList                  // 登录奖励
	G_EventCityList       *EventCityInfoList              // 活动关卡
	G_VersionMap          map[string]*VersionInfo         // 版本控制
	G_VipConfList         *VipConfList                    // vip配置
	G_VipLevelPay         map[int32]float32               // vip等级支付
	G_VipLevelStaminaAdd  map[int32]int32                 // vip等级体力增加
	G_VipLevelCostAdd     map[int32]int32                 // vip等级COST增加
	G_VipLevelResumeQuest map[int32]int32                 // vip等级复活打折次数
	G_VipLevelBuyStamina  map[int32]int32                 // vip等级购买体力次数
	G_VipLevelBuyMoney    map[int32]int32                 // vip等级购买金币次数
)

func init() {
	log.T("config.go >>> step 1...")

	TableDevourCostCoin = make([]int32, 99)
	//系统语言
	G_GAME_LANGUAGE = Language_EN
	//公告
	//登录奖励
	G_LoginPresentConf = new(BonusConfList)
	//活动关卡
	G_EventCityList = new(EventCityInfoList)
	//版本控制
	G_VersionMap = make(map[string]*VersionInfo)
	//vip系统
	G_VipConfList = new(VipConfList)
	G_VipLevelPay = make(map[int32]float32)
	G_VipLevelStaminaAdd = make(map[int32]int32)
	G_VipLevelCostAdd = make(map[int32]int32)
	G_VipLevelResumeQuest = make(map[int32]int32)
	G_VipLevelBuyStamina = make(map[int32]int32)
	G_VipLevelBuyMoney = make(map[int32]int32)
	//任务系统
	//副本系统

	log.T("config.go >>> step end...")
}

// 服务器初始化
func InitConfig(reloadFlg bool) (e Error.Error) {
	fmt.Println(">>>>>>>>>>InitConfig..")
	//conf/conf.ini
	cfg.Reload()
	//初始化日志处理
	if reloadFlg {
		reloadLogger()
	} else {
		InitLogger()
		log.T("initlogger...")
	}
	fmt.Println("InitLanguage...")

	//初始化系统语言
	InitLanguage()



	//初始化redis连接配置
	fmt.Println("InitRedis...")
	InitRedis()


	//初始化mongoDb
	fmt.Println("InitMongoDb...")
	errInitMongoDb := InitMongoDb()
	if errInitMongoDb != nil {
		panic(errInitMongoDb)
	}

	fmt.Println("InitVersion...")
	//强化时需要花费的金币
	for i := int32(0); i < 99; i++ {
		TableDevourCostCoin[i] = 100 * (i + 1)
	}

	fmt.Println("LoadNotice...")

	fmt.Println(">>>>>>>>>end...")
	return Error.OK()
}

func InitLanguage() {
	//配置文件:conf.ini
	config := cfg.Get()

	if len(config["game_language"]) != 0 {
		G_GAME_LANGUAGE = config["game_language"]
	}

	log.T("当前系统语言: %v", G_GAME_LANGUAGE)
}

/**
初始化数据库
1,建立自增主键
2,建立索引


 */
func InitMongoDb() error {
	//return errors.New("初始化数据库失败")

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


//只返回rank这一级的exp（not total）
func GetUserRankExp(rank int32) int32 {
	if rank > consts.N_MAX_USER_RANK {
		rank = consts.N_MAX_USER_RANK
	}

	return TableUserRankExp[rank - 1]
	//	return CommonCurveValue(rank, consts.N_MAX_USER_RANK, 0, MAX_USER_RANK_EXP, USER_RANK_GROW_CURVE)
}

//return total exp by rank
func GetTotalRankExp(rank int32) (totalExp int32) {
	totalExp = int32(0)
	for r := int32(1); r <= rank; r++ {
		//log.T("GetUserRankExp(%v)=%v",rank, GetUserRankExp(r))
		totalExp += GetUserRankExp(r)
	}
	return totalExp
}

//return user rank by exp (in total)
func GetRankByExp(exp int32) (rank int32) {
	totalExp := int32(0)
	for rank := int32(1); rank < consts.N_MAX_USER_RANK; rank++ {
		totalExp += GetUserRankExp(rank)
		if totalExp >= exp {
			log.T("GetRankByExp( exp=%v) return new rank=%v.", exp, rank)
			return rank
		}
	}
	return 1
}

func GetCostMax(rank int32) int32 {
	if rank < 1 || rank > consts.N_MAX_USER_RANK {
		log.Error("GetCostMax :: invalid rank:%v", rank)
		return -1
	}
	return TableCostMax[rank - 1]
}

func GetUnitMax(rank int32) int32 {
	if rank > 16 {
		rank = 16
	}
	return TableUnitMax[rank - 1]
}

func GetFriendMax(rank int32) int32 {
	if rank > 67 {
		rank = 67
	}
	return TableFriendMax[rank - 1]
}

func GetStaminaMax(rank int32) int32 {
	if rank > consts.N_MAX_USER_RANK {
		rank = consts.N_MAX_USER_RANK
	}
	return TableStaminaMax[rank - 1]
}

func GetCurveValue(lv, maxLv, minValue, maxValue int32, growCurve float32) (result int32) {
	// growCurve= 0.7, 1.0, 1.5

	V := float64(minValue) + float64(maxValue - minValue) * math.Pow(float64(lv - 1) / float64(maxLv - 1), float64(growCurve))

	log.T("GetCurveValue=%v :: level:[%v,%v] min:%v max:%v grow:%v", int32(math.Floor(V)), lv, maxLv, minValue, maxValue, growCurve)

	return int32(math.Floor(V))
}
