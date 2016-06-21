package casinoConf

const (
	//------------------------------登陆注册-------------------------
	LOGIN_WAY_QUICK = 1		//快速登录模式
	LOGIN_WAY_LOGIN = 2		//普通登录模式

	//------------------------------数据库相关--------------------------
	DB_IP 			= "localhost"		//数据库ip
	DB_PORT			= 51668			//数据库端口
	DB_NAME 		= "test"		//数据库名字
	DB_ENSURECOUNTER_KEY	= "id"			//自增键

	DBT_T_USER 		= "t_user"		//user表名字
	DBT_T_TEST 		= "t_test"		//user表名字
	DBT_T_SUB2 		= "t_test_sub2"		//user表名字

	//数据库的常用设置
	MIN_USER_ID = 10000		//USER_ID 起始
	MAX_USER_ID = 100000000		//USER_ID 限制




	//--------------------------------活动奖励-----------------------------

	REWARDS_TYPE_ONLINE		= 1				//在线奖励
	REWARDS_TYPE_TURNTABLE		= 2				//转盘奖励
	REWARDS_TYPE_SIGNIN		= 3				//签到奖励
	REWARDS_TYPE_TIMING		= 4				//定时领奖
	REWARDS_TYPE_RELIEF		= 5				//救济金



)


