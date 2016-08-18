package intCons

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



// 常用的数字变量,方便取地址 1--10
var(
	NUM_INT32_0  int32 = 0
	NUM_INT64_0  int64 = 0
)


var GAME_TYPE_TH	int32 =	1		//普通的德州
var GAME_TYPE_TH_CS	int32 =	2		//锦标赛德州




//返回值
var(
	ACK_RESULT_ERROR	int32 =	-1
	ACK_RESULT_SUCC		int32 =	0
	ACK_RESULT_FAIL		int32 =	-2
)




//