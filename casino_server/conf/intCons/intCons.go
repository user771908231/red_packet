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




)

var  SGJV ShuiGuoJiValue = ShuiGuoJiValue{
	Apple		:			5,
	AppleLittle	:			3,
	Orange		:			10,
	OrangeLittle	:			3,
	Mango		:			15,
	MangoLittle	:			3,
	Bell		:			20,
	BellLittle	:			3,
	Watermelon	:			20,
	WatermelonLittle:			3,
	Star		:			30,
	StarLittle	:			3,
	S77		:			40,
	S77Little	:			3,
	Bar		:			120,
	BarLittle	:			50,
	Litter		:			3,
	Lucky		:			0,

	TypeBet		:			1,
	TypeHilomp	:			2,

}




/**

水果机押注对应的积分的数据
 */
type ShuiGuoJiValue struct {
	Apple 			int32
	AppleLittle		int32
	Orange 			int32
	OrangeLittle		int32
	Mango 			int32
	MangoLittle 		int32
	Bell 			int32
	BellLittle 		int32
	Watermelon 		int32
	WatermelonLittle	int32
	Star 			int32
	StarLittle 		int32
	S77 			int32
	S77Little		int32
	Bar 			int32
	BarLittle		int32
	Litter			int32
	Lucky			int32
	TypeBet			int32
	TypeHilomp		int32

}



