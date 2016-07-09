package casinoConf

const (

	//------------------------------mongo-数据库相关--------------------------
	DB_IP 			= "localhost"		//数据库ip
	DB_PORT			= 51668			//数据库端口
	DB_NAME 		= "test"		//数据库名字
	DB_ENSURECOUNTER_KEY	= "id"			//自增键

	DBT_T_USER 		= "t_user"		//user表名字
	DBT_T_TEST 		= "t_test"		//user表名字
	DBT_T_SUB2 		= "t_test_sub2"		//user表名字
	DBT_T_ZJH_ROUND		= "t_zjh_round"		//每局炸扎金花的数据
	DBT_T_BONUS_TURNTABLE	= "t_bonus_turntable"	//转盘奖励的表

	//数据库的常用设置
	MIN_USER_ID = 10000		//USER_ID 起始
	MAX_USER_ID = 100000000		//USER_ID 限制


	//------------------------------redis-数据库相关--------------------------

	REDIS_IP		=	"127.0.0.1"
	REDIS_PORT		=	"6379"
	REDIS_DB_NAME		=	"test"

	//-------------------------------游戏设置---------------------------------
	SWITCH_ZJH	bool	= false
)
