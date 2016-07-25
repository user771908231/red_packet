package intCons

import "casino_server/common/log"

var (
	//--------------------------------统一返回码------------------------------------------------------------------
	CODE_SUCC			int32		=	0		//成功时候的代码
	CODE_FAIL			int32		=	-1		//失败时候的代码

	//--------------------------------活动奖励--------------------------------------------------------------------
	REWARDS_TYPE_ONLINE		int32		= 1				//在线奖励
	REWARDS_TYPE_TURNTABLE		int32		= 2				//转盘奖励
	REWARDS_TYPE_SIGNIN		int32		= 3				//签到奖励
	REWARDS_TYPE_TIMING		int32		= 4				//定时领奖
	REWARDS_TYPE_RELIEF		int32		= 5				//救济金

	//------------------------------登陆注册----------------------------------------------------------------------
	LOGIN_WAY_QUICK			int8		= 1		//快速登录模式
	LOGIN_WAY_LOGIN 		int8		= 2		//普通登录模式



	//------------------------------请求房间----------------------------------------------------------------------

	REQ_TYPE_IN			int32		= 1
	REQ_TYPE_OUT			int32 		= 0
)

var (
	//纸牌的面值
	POKER_HEART_2  		int32	= 0
	POKER_DIAMOND_2  	int32	= 1
	POKER_CLUB_2  		int32	= 2
	POKER_SPADE_2  		int32	= 3
	POKER_HEART_3  		int32	= 4
	POKER_DIAMOND_3  	int32	= 5
	POKER_CLUB_3  		int32	= 6
	POKER_SPADE_3  		int32	= 7
	POKER_HEART_4		int32	= 8
	POKER_DIAMOND_4		int32	= 9
	POKER_CLUB_4		int32	= 10
	POKER_SPADE_4		int32	= 11
	POKER_HEART_5		int32	= 12
	POKER_DIAMOND_5		int32	= 13
	POKER_CLUB_5		int32	= 14
	POKER_SPADE_5		int32	= 15
	POKER_HEART_6		int32	= 16
	POKER_DIAMOND_6		int32	= 17
	POKER_CLUB_6		int32	= 18
	POKER_SPADE_6		int32	= 19
	POKER_HEART_7		int32	= 20
	POKER_DIAMOND_7		int32	= 21
	POKER_CLUB_7		int32	= 22
	POKER_SPADE_7		int32	= 23
	POKER_HEART_8		int32	= 24
	POKER_DIAMOND_8		int32	= 25
	POKER_CLUB_8		int32	= 26
	POKER_SPADE_8		int32	= 27
	POKER_HEART_9		int32	= 28
	POKER_DIAMOND_9		int32	= 29
	POKER_CLUB_9		int32	= 30
	POKER_SPADE_9		int32	= 31
	POKER_HEART_10		int32	= 32
	POKER_DIAMOND_10	int32	= 33
	POKER_CLUB_10		int32	= 34
	POKER_SPADE_10		int32	= 35
	POKER_HEART_11_J	int32	= 36
	POKER_DIAMOND_11_J	int32	= 37
	POKER_CLUB_11_J		int32	= 38
	POKER_SPADE_11_J	int32	= 39
	POKER_HEART_12_Q	int32	= 40
	POKER_DIAMOND_12_Q	int32	= 41
	POKER_CLUB_12_Q		int32	= 42
	POKER_SPADE_12_Q	int32	= 43
	POKER_HEART_13_K	int32	= 44
	POKER_DIAMOND_13_K	int32	= 45
	POKER_CLUB_13_K		int32	= 46
	POKER_SPADE_13_K	int32	= 47
	POKER_HEART_14_A	int32	= 48
	POKER_DIAMOND_14_A	int32	= 49
	POKER_CLUB_14_A		int32	= 50
	POKER_SPADE_14_A	int32	= 51
	POKER_RED_JOKER		int32	= 52
	POKER_BLACK_JOKER	int32	= 53
)


//可以模拟一下纸牌的数据
type Pai struct {
	Pai1 int32
	Pai2 int32
	Pai3 int32
}

/**
	同花顺的
 */
type TongHuaShun Pai
type DuiZi Pai
type SanPai Pai
//模拟的数据

var TongHuaShun1 TongHuaShun = TongHuaShun{};

func init(){
	log.T("",TongHuaShun1)
	log.T("",TongHuaShun1.Pai1)
	log.T("",POKER_CLUB_2)

	TongHuaShun1.Pai1 = POKER_CLUB_2
	TongHuaShun1.Pai2 = POKER_CLUB_3
	TongHuaShun1.Pai3 = POKER_CLUB_4
}


// 常用的数字变量,方便取地址 1--10
var(
	NUM_INT32_0  int32 = 0
	NUM_INT64_0  int64 = 0
)

