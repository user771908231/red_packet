package config

import "casino_server/conf/casinoConf"

const (
	MJ_DBNAM = casinoConf.DB_NAME
	DB_ENSURECOUNTER_KEY = "id"                        //自增键
	DBT_MJ_DESK = "t_mj_desk"			   //桌子的信息
	DBT_MJ_DESK_ROUND = "dbT_mj_desk_round"		   //一把麻将结束

	DBT_T_TH_GAMENUMBER_SEQ = "t_th_gamenumber_seq"    //麻将 编号
)

