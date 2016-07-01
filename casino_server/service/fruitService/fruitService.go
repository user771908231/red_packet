package fruitService

import (
	"casino_server/msg/bbprotogo"
	"casino_server/utils"
	"casino_server/common/log"
	"reflect"
	"github.com/name5566/leaf/gate"
	"errors"
)



/**

水果机的 index
 */

const(
	INDEX_ORANGE_1 		= 	0
	INDEX_BELL_1		=	1
	INDEX_BAR_LITTLE	=	2
	INDEX_BAR		=	3
	INDEX_APPLE_1		=	4
	INDEX_APPLE_LITTLE	=	5
	INDEX_MANGO_1		=	6
	INDEX_WATERMELON	=	7
	INDEX_WATERMELON_LITTLE =	8
	INDEX_LUCK_1		=	9
	INDEX_APPLE_2		=	10
	INDEX_ORANGE_LITTLE	=	11
	INDEX_ORANGE_2		=	12
	INDEX_BELL_2		=	13
	INDEX_77_LITTLE		=	14
	INDEX_77		=	15
	INDEX_APPLE_3		=	16
	INDEX_MANGO_LITTLE	=	17
	INDEX_MANGO_2		=	18
	INDEX_STAR		=	19
	INDEX_STAR_LITTLE	=	20
	INDEX_LUCK_2		=	21
	INDEX_APPLE_4	 	=	22
	INDEX_BELL_LITTLE 	=	23
)




/**
水果几的类型,是押注还是赌大小
 */
var SGJ_TYPE_Bet int32 = 1
var SGJ_TYPE_Hilomp int32 = 2


//

//奖励的类型

var SGJ_WIN_TYPE_N 		int32 = 1		//常见的一炮
var SGJ_WIN_TYPE_LUCK0 		int32 = 2
var SGJ_WIN_TYPE_LUCK1 		int32 = 3
var SGJ_WIN_TYPE_LUCK2 		int32 = 4
var SGJ_WIN_TYPE_LUCK3 		int32 = 5
var SGJ_WIN_TYPE_DASIXI 	int32 = 6
var SGJ_WIN_TYPE_DASANYUAN 	int32 = 7
var SGJ_WIN_TYPE_XIAOSANYUAN 	int32 = 8
var SGJ_WIN_TYPE_RAND1 		int32 = 9
var SGJ_WIN_TYPE_RAND2 		int32 = 10
var SGJ_WIN_TYPE_RAND3 		int32 = 11
var SGJ_WIN_TYPE_ZONGHENGSIHAI 	int32 = 12
var SGJ_WIN_TYPE_WU	 	int32 = 13



//奖励的个数
var NUMBER_INT_1 int32 = 1

/**
配置每个水果对应的积分
 */
var SGJV SGJValue = SGJValue{
	Apple                	:                        5,
	AppleLittle        	:                        3,
	Orange                	:                        10,
	OrangeLittle        	:                        3,
	Mango                	:                        15,
	MangoLittle        	:                        3,
	Bell                	:                        20,
	BellLittle        	:                        3,
	Watermelon        	:                        20,
	WatermelonLittle	:                        3,
	Star                	:                        30,
	StarLittle        	:                        3,
	S77                	:                        40,
	S77Little        	:                        3,
	Bar                	:                        120,
	BarLittle        	:                        50,
	Litter                	:                        3,
	Lucky                	:                        0,
}


/**
概率设计
所有的概率都需要以此结构体来做配置
 */
type ShuiGuoPro struct {
	ORANGE_1 	  int32
	BELL_1	          int32
	BAR_LITTLE        int32
	BAR	          int32
	APPLE_1	          int32
	APPLE_LITTLE      int32
	MANGO_1	          int32
	WATERMELON        int32
	WATERMELON_LITTLE int32
	LUCK_1	          int32
	APPLE_2	          int32
	ORANGE_LITTLE     int32
	ORANGE_2	  int32
	BELL_2	          int32
	S77_LITTLE	  int32
	S77	          int32
	APPLE_3	          int32
	MANGO_LITTLE      int32
	MANGO_2	          int32
	STAR	          int32
	STAR_LITTLE       int32
	LUCK_2	          int32
	APPLE_4	          int32
	BELL_LITTLE	  int32

	DASIXI            int32 //大四喜
	DASANYUAN         int32 //大三元
	XIAOSANYUAN       int32 //小三元
	RAND1		  int32 //随机赠送一炮
	RAND2             int32 //随机赠送两炮
	RAND3             int32 //随机赠送三炮
	RAND4             int32 //随机赠送四炮
	ZONGHENGSIHAI     int32 //纵横四海
	WU		  int32	//没有跑到的概率(注意:只有先跑到luck之后才有这个概率)
	MAX 		  int32	//产生随机数的区间 (0,max]
}

//暂时用不了了  这个方法
func ( p *ShuiGuoPro) GetMax() int32{
	v := reflect.ValueOf(*p)
	count := v.NumField()
	var result int32 = int32(v.Field(0).Interface().(int32))
	for i:=0;i<count ;i++  {
		f :=v.Field(i)
		if f.Interface().(int32) > result {
			result = f.Interface().(int32)
		}
	}
	log.T("概率最大的值,",result)
	return result
}





/**
跑中luck 时,其他的概率
 */
var SGJLuckP ShuiGuoPro = ShuiGuoPro{
	ORANGE_1 	     	:        5,
	BELL_1	             	:        6,
	BAR_LITTLE           	:        7,
	BAR	             	:        8,
	APPLE_1	             	:        15,
	APPLE_LITTLE         	:        20,
	MANGO_1			:        25,
	WATERMELON            	:        30,
	WATERMELON_LITTLE       :        35,
	LUCK_1	           	:        0,
	APPLE_2	         	:        42,
	ORANGE_LITTLE           :        43,
	ORANGE_2	       	:        46,
	BELL_2	           	:        47,
	S77_LITTLE	        :        48,
	S77	          	:        49,
	APPLE_3	         	:        50,
	MANGO_LITTLE          	:        60,
	MANGO_2	         	:        68,
	STAR	             	:        70,
	STAR_LITTLE             :        71,
	LUCK_2	             	:        0,
	APPLE_4	             	:        73,
	BELL_LITTLE	 	:	 73,

	DASIXI                  :        74, //大四喜
	DASANYUAN               :        75, //大三元
	XIAOSANYUAN             :        76, //小三元
	RAND1		        :        78, //随机两炮
	RAND2                   :        80, //随机三炮
	RAND3                   :        90, //Bar
	RAND4                   :        100, //纵横四海
	ZONGHENGSIHAI           :        100,

}


/**
水果机单次中的概率
 */
var SGJP ShuiGuoPro = ShuiGuoPro{
	ORANGE_1 	        :        0,
	BELL_1	               	:        10,
	BAR_LITTLE              :        20,
	BAR	                :        30,
	APPLE_1	                :        40,
	APPLE_LITTLE           	:        50,
	MANGO_1	                :        60,
	WATERMELON            	:        70,
	WATERMELON_LITTLE     	:        80,
	LUCK_1	            	:        0,
	APPLE_2	                :        100,
	ORANGE_LITTLE         	:        110,
	ORANGE_2	       	:        120,
	BELL_2	             	:        130,
	S77_LITTLE	       	:        140,
	S77	             	:        150,
	APPLE_3	                :        160,
	MANGO_LITTLE            :        170,
	MANGO_2	                :        180,
	STAR	                :        190,
	STAR_LITTLE             :        200,
	LUCK_2	            	:        0,
	APPLE_4	              	:        220,
	BELL_LITTLE	        :        230,

	DASIXI                  :        0,
	DASANYUAN               :        0, //Bar
	XIAOSANYUAN             :        0, //纵横四海
	RAND1		      	:        0,
	RAND2 			:   	 0,
	RAND3                   :        0, //Bar
	RAND4                   :        0, //纵横四海
	ZONGHENGSIHAI           :        0,
	MAX 			:	 230,
}


/**
	随机赠送的时候的概率
 */
var SGJPRand ShuiGuoPro = ShuiGuoPro{
	ORANGE_1 	        :        5,
	BELL_1	               	:        15,
	BAR_LITTLE              :        20,
	BAR	                :        25,
	APPLE_1	                :        30,
	APPLE_LITTLE           	:        35,
	MANGO_1	                :        40,
	WATERMELON            	:        42,
	WATERMELON_LITTLE     	:        43,
	LUCK_1	            	:        0,
	APPLE_2	                :        47,
	ORANGE_LITTLE         	:        48,
	ORANGE_2	       	:        49,
	BELL_2	             	:        50,
	S77_LITTLE	       	:        60,
	S77	             	:        68,
	APPLE_3	                :        70,
	MANGO_LITTLE            :        71,
	MANGO_2	                :        72,
	STAR	                :        73,
	STAR_LITTLE             :        74,
	LUCK_2	            	:        0,
	APPLE_4	              	:        76,
	BELL_LITTLE	        :        100,
	DASIXI                  :        0,
	DASANYUAN               :        0,
	XIAOSANYUAN             :        0,
	RAND1		      	:        0,
	RAND2 			:	 0,
	RAND3                   :        0, //Bar
	RAND4                   :        0, //纵横四海
	ZONGHENGSIHAI           :        0,
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
	得到一次的结果
	水果机器的结果有可能有很多种,这里需要什么策略来返回结果?
 */
func HandlerShuiguoji(m *bbproto.Shuiguoji,a gate.Agent) (*bbproto.ShuiguojiRes, error) {
	//1,检测参数并且根据押注的内容选择处理方式
	if m == nil {
		return nil,nil
	}

	//2,活的返回值
	result := &bbproto.ShuiguojiRes{}
	result.Result = make([]int32,24)
	p := SGJP
	var seq int32 = 1		//初始顺序号
	err := BetResultWin(m,result, &p,&seq)
	if err != nil {
		log.E(err.Error())
		log.E("获取水果机结果的时候出错")
	}

	//跑到的结棍
	log.T("押注:%v",m)
	log.T("跑到的结果:%v",result.Result)

	//计算得分结果
	var scoresTotal int32 = 0

	if result.Result[INDEX_ORANGE_1] > 0 {
		scoresTotal += (m.GetScoresOrange() * SGJV.Orange)
	}

	if result.Result[INDEX_BELL_1] > 0 {
		scoresTotal += (m.GetScoresBell() * SGJV.Bell)
	}

	if result.Result[INDEX_BAR_LITTLE] > 0 {
		scoresTotal += (m.GetScoresBar() * SGJV.BarLittle)
	}

	if result.Result[INDEX_BAR] > 0 {
		scoresTotal += (m.GetScoresBar()  * SGJV.Bar)
	}

	if result.Result[INDEX_APPLE_1] > 0 {
		scoresTotal += (m.GetScoresApple()  * SGJV.Apple)

	}

	if result.Result[INDEX_APPLE_LITTLE] > 0 {
		scoresTotal += (m.GetScoresApple()  * SGJV.AppleLittle)
	}

	if result.Result[INDEX_MANGO_1] > 0 {
		scoresTotal += (m.GetScoresMango()  * SGJV.Mango)
	}

	if result.Result[INDEX_WATERMELON] > 0 {
		scoresTotal += (m.GetScoresWatermelon()  * SGJV.Watermelon)
	}

	if result.Result[INDEX_WATERMELON_LITTLE] > 0 {
		scoresTotal += (m.GetScoresWatermelon()  * SGJV.WatermelonLittle)
	}

	if result.Result[INDEX_APPLE_2] > 0 {
		scoresTotal += (m.GetScoresApple()  * SGJV.Apple)
		log.N("计算得到的总分是INDEX_APPLE_2%v", scoresTotal)
	}

	if result.Result[INDEX_ORANGE_LITTLE] > 0 {
		scoresTotal += (m.GetScoresOrange()  * SGJV.OrangeLittle)
	}

	if result.Result[INDEX_ORANGE_2] > 0 {
		scoresTotal += (m.GetScoresOrange()  * SGJV.Orange)
	}

	if result.Result[INDEX_BELL_2] > 0 {
		scoresTotal += (m.GetScoresBell()  * SGJV.Bell)
	}

	if result.Result[INDEX_77_LITTLE] > 0 {
		scoresTotal += (m.GetScores77()  * SGJV.S77Little)
	}
	if result.Result[INDEX_77] > 0 {
		scoresTotal += (m.GetScores77()  * SGJV.S77)
	}
	if result.Result[INDEX_APPLE_3] > 0 {
		scoresTotal += (m.GetScoresApple()  * SGJV.Apple)
		log.N("计算得到的总分是INDEX_APPLE_3%v", scoresTotal)
	}
	if result.Result[INDEX_MANGO_LITTLE] > 0 {
		scoresTotal += (m.GetScoresMango()  * SGJV.MangoLittle)
		log.N("计算得到的总分是INDEX_MANGO_LITTLE%v", scoresTotal)
	}
	if result.Result[INDEX_MANGO_2] > 0 {
		scoresTotal += (m.GetScoresMango()  * SGJV.Mango)
		log.N("计算得到的总分是INDEX_MANGO_2%v", scoresTotal)
	}
	if result.Result[INDEX_STAR] > 0 {
		scoresTotal += (m.GetScoresStar()  * SGJV.Star)
		log.N("计算得到的总分是INDEX_STAR%v", scoresTotal)
	}

	if result.Result[INDEX_STAR_LITTLE] > 0 {
		scoresTotal += (m.GetScoresStar()  * SGJV.StarLittle)
		log.N("计算得到的总分是INDEX_STAR_LITTLE%v", scoresTotal)
	}

	if result.Result[INDEX_APPLE_4] > 0 {
		scoresTotal += (m.GetScoresApple()  * SGJV.Apple)
		log.N("计算得到的总分是INDEX_APPLE_4%v", scoresTotal)
	}
	if result.Result[INDEX_BELL_LITTLE] > 0 {
		scoresTotal += (m.GetScoresBell()  * SGJV.BellLittle)
	}

	result.ScoresWin = &scoresTotal
	log.N("计算得到的总分是%v", result.GetScoresWin())

	//更新用户的余额信息

	updateUserData(result,a);
	//返回值
	return result, err

}


/**

水果机比大小的业务

 */
func HandlerShuiguojiHilomp(m *bbproto.ShuiguojiHilomp) (*bbproto.ShuiguojiHilomp, error) {
	return HilompResult(m.GetProtoHeader().GetUserId())
}


/**
得到跑到的结果
m: 请求的参数
 */
func BetResultWin(m *bbproto.Shuiguoji, res *bbproto.ShuiguojiRes,p *ShuiGuoPro, seq *int32) (error) {
	r := utils.Randn(p.GetMax())
	log.T("水果机 随机数 %v", r)
	log.T("水果机 概率 %v", *p)

	if r < p.ORANGE_1 {
		setResult(m,res,p,INDEX_ORANGE_1,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//

	} else if r < p.BELL_1 {
		setResult(m,res,p,INDEX_BELL_1,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//

	} else if r < p.BAR_LITTLE {
		setResult(m,res,p,INDEX_BAR_LITTLE,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//

	} else if r < p.BAR {
		setResult(m,res,p,INDEX_BAR,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.APPLE_1 {
		setResult(m,res,p,INDEX_APPLE_1,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.WATERMELON {
		setResult(m,res,p,INDEX_WATERMELON,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.WATERMELON_LITTLE {
		setResult(m,res,p,INDEX_WATERMELON_LITTLE,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.LUCK_1 {
		//todo  这里需要处理luck  时候的中奖信息(两个luck都是一样的处理方法)
		pluck := SGJLuckP
		setResult(m,res,&pluck,INDEX_LUCK_1,seq)
		BetResultWin(m,res,&pluck,seq)

	} else if r < p.APPLE_2 {
		setResult(m,res,p,INDEX_APPLE_2,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.ORANGE_LITTLE {
		setResult(m,res,p,INDEX_ORANGE_LITTLE,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.ORANGE_2 {
		setResult(m,res,p,INDEX_ORANGE_2,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.BELL_2 {
		setResult(m,res,p,INDEX_BELL_2,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.S77_LITTLE {
		setResult(m,res,p,INDEX_77_LITTLE,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.S77 {
		setResult(m,res,p,INDEX_77,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.APPLE_3 {
		setResult(m,res,p,INDEX_APPLE_3,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.MANGO_LITTLE {
		setResult(m,res,p,INDEX_MANGO_LITTLE,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.MANGO_2 {
		setResult(m,res,p,INDEX_MANGO_2,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.STAR {
		setResult(m,res,p,INDEX_STAR,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.STAR_LITTLE {
		setResult(m,res,p,INDEX_STAR_LITTLE,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.LUCK_2 {
		pluck := SGJLuckP
		setResult(m,res,&pluck,INDEX_LUCK_2,seq)
		BetResultWin(m,res,&pluck,seq)
	} else if r < p.APPLE_4 {
		setResult(m,res,p,INDEX_APPLE_4,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.BELL_LITTLE {
		setResult(m,res,p,INDEX_BELL_LITTLE,seq)
		setWinType(res,&SGJ_WIN_TYPE_N)//
	} else if r < p.DASIXI {	//大四喜

		setResult(m,res,p,INDEX_APPLE_1,seq)
		setResult(m,res,p,INDEX_APPLE_2,seq)
		setResult(m,res,p,INDEX_APPLE_3,seq)
		setResult(m,res,p,INDEX_APPLE_4,seq)

	} else if r < p.DASANYUAN {	//大三元
		setResult(m,res,p,INDEX_WATERMELON,seq)
		setResult(m,res,p,INDEX_STAR,seq)
		setResult(m,res,p,INDEX_77,seq)
	} else if r < p.XIAOSANYUAN {	//小三元
		setResult(m,res,p,INDEX_ORANGE_1,seq)
		setResult(m,res,p,INDEX_MANGO_1,seq)
		setResult(m,res,p,INDEX_BELL_1,seq)

	} else if r < p.RAND1 {	//随机赠送一炮,由于随机炮的概率是不同的,所以不用考虑随机炮再次打到随机炮的问题
		log.T("彩蛋,随机赠送一炮")
		setWinType(res,&SGJ_WIN_TYPE_RAND1)//没有中奖
		pr1 := SGJPRand
		BetResultWin(m,res,&pr1,seq)
		BetResultWin(m,res,&pr1,seq)
	} else if r < p.RAND2 {
		log.T("彩蛋,随机赠送二炮")
		pr2 := SGJPRand
		BetResultWin(m,res,&pr2,seq)
		BetResultWin(m,res,&pr2,seq)
		BetResultWin(m,res,&pr2,seq)
	} else if r < p.RAND3 {
		log.T("彩蛋,随机赠送三炮")
		pr3 := SGJPRand
		BetResultWin(m, res, &pr3,seq)
		BetResultWin(m, res, &pr3,seq)
		BetResultWin(m, res, &pr3,seq)
		BetResultWin(m, res, &pr3,seq)
	}else if r < p.ZONGHENGSIHAI {
	} else if r < p.WU {
		setWinType(res,&SGJ_WIN_TYPE_WU)//没有中奖
		res.Result[INDEX_LUCK_2]   = 1
		setResult(m,res,p,INDEX_ORANGE_1,seq)
	}

	return  nil

}

/**
跑中之后设置值
 */
func setResult(m *bbproto.Shuiguoji, res *bbproto.ShuiguojiRes,p *ShuiGuoPro,index int,seq *int32) error{
	if res.Result[index] > 0 {
		BetResultWin(m,res,p,seq)
	}else{
		res.Result[index] = *seq
		*seq ++		//顺序号加1
	}
	return nil
}

/**
设置奖励类型
 */
func setWinType(res *bbproto.ShuiguojiRes,value *int32){
	if res.GetWinType() <= 0 {
		res.WinType = value
	}

}

/**
	比大小的结果
	一般情况下
 */
func HilompResult(id uint32) (*bbproto.ShuiguojiHilomp, error) {
	result := &bbproto.ShuiguojiHilomp{}
	return result, nil
}


/**
更新用户的信息
	1,本次等分
	2,本次水果机剩余得分

 */
func updateUserData(res *bbproto.ShuiguojiRes,m gate.Agent) error{
	//检查参数是否正确
	user := m.UserData()	//存放的是指针
	if user == nil  {
		log.E("agent 中取User 的时候出错")
		return errors.New("没有找到用户")
	}

	//2做更新操作

	return nil

}