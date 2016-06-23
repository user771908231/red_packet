package fruitService

import (
	"casino_server/msg/bbproto"
	"casino_server/utils"
	"casino_server/common/log"
)


/**
水果几的类型,是押注还是赌大小
 */
var SGJ_TYPE_Bet int32 = 1
var SGJ_TYPE_Hilomp int32 = 2


//奖励的类型
var SGJ_WIN_TYPE_LUCK0 int32 = 1
var SGJ_WIN_TYPE_LUCK1 int32 = 2
var SGJ_WIN_TYPE_LUCK2 int32 = 3
var SGJ_WIN_TYPE_LUCK3 int32 = 4
var SGJ_WIN_TYPE_DASIXI int32 = 5
var SGJ_WIN_TYPE_DASANYUAN int32 = 6
var SGJ_WIN_TYPE_XIAOSANYUAN int32 = 7
var SGJ_WIN_TYPE_RAND2 int32 = 8
var SGJ_WIN_TYPE_RAND3 int32 = 9
var SGJ_WIN_TYPE_ZONGHENGSIHAI int32 = 10

//奖励的个数
var NUMBER_INT_1 int32 = 1

/**
配置每个水果对应的积分
 */
var SGJV SGJValue = SGJValue{
	Apple                :                        5,
	AppleLittle        :                        3,
	Orange                :                        10,
	OrangeLittle        :                        3,
	Mango                :                        15,
	MangoLittle        :                        3,
	Bell                :                        20,
	BellLittle        :                        3,
	Watermelon        :                        20,
	WatermelonLittle:                        3,
	Star                :                        30,
	StarLittle        :                        3,
	S77                :                        40,
	S77Little        :                        3,
	Bar                :                        120,
	BarLittle        :                        50,
	Litter                :                        3,
	Lucky                :                        0,
}

var SGJP ShuiGuoPro = ShuiGuoPro{
	Apple                        :        5,
	AppleLittle                :        15,
	Orange                        :        20,
	OrangeLittle                :        25,
	Mango                        :        30,
	MangoLittle                :        35,
	Bell                        :        40,
	BellLittle                :        42,
	Watermelon                :        43,
	WatermelonLittle        :        46,
	Star                        :        47,
	StarLittle                :        48,
	S77                        :        49,
	S77Little                :        50,
	Bar                        :        60,
	BarLittle                :        68,
	Litter                        :        69,
	Lucky0                        :        70,
	Lucky1                        :        71,
	Lucky2                        :        72,
	Lucky3                        :        73,

	DaSiXi                        :        74, //大四喜
	DaSanYuan                :        75, //大三元
	XiaoSanYuan                :        76, //小三元
	Wu                        :        77, //空,吃掉
	Rand2                        :        78, //随机两炮
	Rand3                        :        80, //随机三炮
	rand4                        :        90, //Bar
	ZongHengSiHai                :        100, //纵横四海
	MAX                        :        100,

}


/**

水果机押注对应的积分的数据
 */
type SGJValue struct {
	Apple            int32
	AppleLittle      int32
	Orange           int32
	OrangeLittle     int32
	Mango            int32
	MangoLittle      int32
	Bell             int32
	BellLittle       int32
	Watermelon       int32
	WatermelonLittle int32
	Star             int32
	StarLittle       int32
	S77              int32
	S77Little        int32
	Bar              int32
	BarLittle        int32
	Litter           int32
	Lucky            int32
}


/**
概率设计
 */
type ShuiGuoPro struct {
	Apple            int32
	AppleLittle      int32
	Orange           int32
	OrangeLittle     int32
	Mango            int32
	MangoLittle      int32
	Bell             int32
	BellLittle       int32
	Watermelon       int32
	WatermelonLittle int32
	Star             int32
	StarLittle       int32
	S77              int32
	S77Little        int32
	Bar              int32
	BarLittle        int32
	Litter           int32
	Lucky0           int32
	Lucky1           int32
	Lucky2           int32
	Lucky3           int32
	TypeBet          int32
	TypeHilomp       int32
	DaSiXi           int32 //大四喜
	DaSanYuan        int32 //大三元
	XiaoSanYuan      int32 //小三元
	Wu               int32 //空,吃掉
	Rand2            int32 //随机两炮
	Rand3            int32 //随机三炮
	rand4            int32 //Bar
	ZongHengSiHai    int32 //纵横四海

	MAX              int32 //设置一个峰值
}


/**
	得到一次的结果
	水果机器的结果有可能有很多种,这里需要什么策略来返回结果?
 */
func HandlerShuiguoji(m *bbproto.Shuiguoji) (*bbproto.ShuiguojiRes, error) {
	//1,检测参数并且根据押注的内容选择处理方式
	if m == nil {
		return nil,nil
	}

	//2,活的返回值
	result, err := BetResultWin(m, nil)
	if err != nil {
		log.E(err.Error())
		log.E("获取水果机结果的时候出错")
	}
	log.N("得到的水果机结果%v", result)
	log.N("苹果%v", m.GetScoresApple() * result.GetWinApple() * SGJV.Apple + m.GetScoresApple() * result.GetWinAppleLittle() * SGJV.AppleLittle)
	log.N("橘子%v", m.GetScoresOrange() * result.GetWinOrange() * SGJV.Orange + m.GetScoresOrange() * result.GetWinOrangeLittle() * SGJV.OrangeLittle)
	log.N("芒果%v", m.GetScoresMango() * result.GetWinMango() * SGJV.Mango + m.GetScoresMango() * result.GetWinMangoLittle() * SGJV.MangoLittle)
	log.N("铃铛%v", m.GetScoresBell() * result.GetWinBell() * SGJV.Bell + m.GetScoresBell() * result.GetWinBellLittle() * SGJV.BellLittle)
	log.N("西瓜%v", m.GetScoresWatermelon() * result.GetWinWatermelon() * SGJV.Watermelon + m.GetScoresWatermelon() * result.GetWinWatermelonLittle() * SGJV.WatermelonLittle)
	log.N("星星%v", m.GetScoresStar() * result.GetWinStar() * SGJV.Star + m.GetScoresStar() * result.GetWinStarLittle() * SGJV.StarLittle)
	log.N("77%v", m.GetScores77() * result.GetWin77() * SGJV.S77 + m.GetScores77() * result.GetWin77Little() * SGJV.S77Little)
	log.N("bar%v", m.GetScoresBar() * result.GetWinBar() * SGJV.Bar + m.GetScoresBar() * result.GetWinBarLittle() * SGJV.BarLittle)

	//计算得分结果
	var scoresTotal int32 = (
	m.GetScoresApple() * result.GetWinApple() * SGJV.Apple + m.GetScoresApple() * result.GetWinAppleLittle() * SGJV.AppleLittle +
	m.GetScoresOrange() * result.GetWinOrange() * SGJV.Orange + m.GetScoresOrange() * result.GetWinOrangeLittle() * SGJV.OrangeLittle +
	m.GetScoresMango() * result.GetWinMango() * SGJV.Mango + m.GetScoresMango() * result.GetWinMangoLittle() * SGJV.MangoLittle +
	m.GetScoresBell() * result.GetWinBell() * SGJV.Bell + m.GetScoresBell() * result.GetWinBellLittle() * SGJV.BellLittle +
	m.GetScoresWatermelon() * result.GetWinWatermelon() * SGJV.Watermelon + m.GetScoresWatermelon() * result.GetWinWatermelonLittle() * SGJV.WatermelonLittle +
	m.GetScoresStar() * result.GetWinStar() * SGJV.Star + m.GetScoresStar() * result.GetWinStarLittle() * SGJV.StarLittle +
	m.GetScores77() * result.GetWin77() * SGJV.S77 + m.GetScores77() * result.GetWin77Little() * SGJV.S77Little +
	m.GetScoresBar() * result.GetWinBar() * SGJV.Bar + m.GetScoresBar() * result.GetWinBarLittle() * SGJV.BarLittle)

	log.N("计算得到的总分是%v,地址%v", scoresTotal, &scoresTotal)
	result.ScoresWin = &scoresTotal
	log.N("计算得到的总分是%v", result.ScoresWin)

	//返回值
	return result, err

}


/**

水果机比大小的业务

 */
func HandlerShuiguojiHilomp(m *bbproto.Shuiguoji) (*bbproto.ShuiguojiHilomp, error) {
	return HilompResult(m.GetProtoHeader().GetUserId())
}



/**
得到跑到的结果
 */
func BetResultWin(m *bbproto.Shuiguoji, res *bbproto.ShuiguojiRes) (*bbproto.ShuiguojiRes, error) {
	result := &bbproto.ShuiguojiRes{}
	r := utils.Randn(SGJP.MAX)

	log.T("水果机 随机数 %v", r)

	if r < SGJP.Apple {
		//大苹果
		result.WinApple = &NUMBER_INT_1
	} else if r < SGJP.AppleLittle {
		result.WinAppleLittle = &NUMBER_INT_1
		//小苹果
	} else if r < SGJP.Orange {
		result.WinOrange = &NUMBER_INT_1
		//橘子
	} else if r < SGJP.OrangeLittle {
		result.WinOrangeLittle = &NUMBER_INT_1
		//小橘子
	} else if r < SGJP.Mango {
		result.WinMango = &NUMBER_INT_1
		//芒果
	} else if r < SGJP.MangoLittle {
		result.WinMangoLittle = &NUMBER_INT_1
		//小芒果
	} else if r < SGJP.Bell {
		//铃铛
		result.WinBell = &NUMBER_INT_1
	} else if r < SGJP.BellLittle {
		//小铃铛
		result.WinBellLittle = &NUMBER_INT_1
	} else if r < SGJP.Watermelon {
		//西瓜
		result.WinWatermelon = &NUMBER_INT_1
	} else if r < SGJP.WatermelonLittle {
		//小西瓜
		result.WinWatermelonLittle = &NUMBER_INT_1
	} else if r < SGJP.Star {
		//星星
		result.WinStar = &NUMBER_INT_1
	} else if r < SGJP.StarLittle {
		//小星星
		result.WinStarLittle = &NUMBER_INT_1
	} else if r < SGJP.S77 {
		//77
		result.Win77 = &NUMBER_INT_1
	} else if r < SGJP.S77Little {
		//小77
		result.Win77Little = &NUMBER_INT_1
	} else if r < SGJP.Bar {
		//bar
		result.WinBar = &NUMBER_INT_1
	} else if r < SGJP.BarLittle {
		//小bar
		result.WinBellLittle = &NUMBER_INT_1
	} else if r < SGJP.Lucky0 {
		//luck
		result.WinType = &SGJ_WIN_TYPE_LUCK0
		result, _ = BetResultWin(m, result)
	} else if r < SGJP.Lucky1 {
		//luck
		result.WinType = &SGJ_WIN_TYPE_LUCK1
		result, _ = BetResultWin(m, result)
	} else if r < SGJP.Lucky2 {
		//luck
		result.WinType = &SGJ_WIN_TYPE_LUCK2
		result, _ = BetResultWin(m, result)
	} else if r < SGJP.Lucky3 {
		//luck
		result.WinType = &SGJ_WIN_TYPE_LUCK3
		result, _ = BetResultWin(m, result)
	} else if r < SGJP.DaSiXi {
		//大四喜
		result.WinType = &NUMBER_INT_1
	} else if r < SGJP.DaSanYuan {
		//大三元
	} else if r < SGJP.XiaoSanYuan {
		//小三元
	} else if r < SGJP.Wu {
		//空
	} else if r < SGJP.Rand2 {
	} else if r < SGJP.Rand3 {
	} else if r < SGJP.rand4 {
	} else if r < SGJP.ZongHengSiHai {
	}

	return result, nil

}

/**
	比大小的结果
 */
func HilompResult(id uint32) (*bbproto.ShuiguojiHilomp, error) {
	result := &bbproto.ShuiguojiHilomp{}
	return result, nil
}
