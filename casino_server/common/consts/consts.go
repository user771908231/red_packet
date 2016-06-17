package consts

const (
	DEFAULT_USER_NAME    = ""
	N_DEFAULT_COST_MAX   = 50
	N_DEFAULT_UNIT_MAX   = 50
	N_DEFAULT_FRIEND_MAX = 50
)

//redis table name list
const (
	TABLE_USER      = "1"
	TABLE_UNIT      = "2"
	TABLE_QUEST     = "3"
	TABLE_SESSION   = "4"
	TABLE_FRIEND    = "5"
	TABLE_QUEST_LOG = "6"
)

//key prefix
const (
	X_UUID      = "X_UUID_"
	X_HELP_INFO = "X_HELP_INFO_"

	X_FRIEND_DATA        = "X_FRI_DATA_"
	X_HELPER_RECORD      = "X_HELPER_REC_"         // A user's helpers
	X_USER_RANK          = "X_USER_RANK_"          // ZSET: uid - rank
	X_PRO_HELPER_EVO     = "X_PRO_HELPER_EVO_"     // for get premium helpers
	X_PRO_HELPER_LEVELUP = "X_PRO_HELPER_LEVELUP_" // for get premium helpers

	//quest
	X_QUEST_CITY   = "X_CITY_"
	X_QUEST_STAGE  = "X_STAGE_"
	X_QUEST_CONFIG = "X_CONFIG_"
	X_QUEST_LOG    = "X_Q_LOG_"
	X_QUEST_RECORD = "X_Q_REC_"

	//unit
	X_UNIT_INFO        = "X_UNIT_"
	X_SKILL_INFO       = "X_SKILL_"
	X_EVOLVE_SESSION   = "X_EVOLVE_"
	X_GACHA_UNIT       = "X_GACHA_"
	X_GACHA_CONF       = "X_GACHA_CONF"
	X_FRESH_GACHA_CONF = "X_FRESH_GACHA_CONF"
	X_SKILL_CONF       = "X_SKILL_CONF"
)

const (
	KEY_MAX_USER_ID  = "K_MAX_USER_ID"
	KEY_MAX_UNIT_ID  = "K_MAX_UNIT_ID"
	KEY_QUEST_PREFIX = "K_QUEST_INFO_"
)

const (
	N_MAX_RARE       = 6
	N_MAX_USER_RANK  = 500
	N_MAX_UNIT_NUM   = 400
	N_MAX_FRIEND_NUM = 200

	N_DUNGEON_GRID_COUNT    = 25
	N_USER_SPACE_PARTS      = 10
	N_FRI_HELPER_COUNT      = 30 //Friend & helper count
	N_FRIEND_COUNT          = 20
	N_HELPER_COUNT          = 10 //Helper count
	N_HELPER_RANK_RANGE     = 5
	N_STAMINA_TIME          = 600 // seconds
	N_QUEST_COLOR_BLOCK_NUM = 2400
	N_GACHA_MAX_COUNT       = 9

	N_UNITMAX_EXPAND_COUNT   = 5
	N_FRIENDMAX_EXPAND_COUNT = 5
	N_BUY_MONEY_COUNT        = 5

	N_FRIEND_HELPER_POINT  = 10
	N_SUPPORT_HELPER_POINT = 5

	N_NORMAL_QUEST_STONE = 20
	N_ELITE_QUEST_STONE  = 50
)

// consume cost
const (
	N_GACHA_FRIEND_COST = 200 // cost 200 friend points
	N_GACHA_BUY_COST    = 300 // cost 300 stone for a gacha

	N_RECOVER_STAMINA_COST  = 60  // cost 60 stone
	N_UNITMAX_EXPAND_COST   = 60  // cost 60 stone
	N_FRIENDMAX_EXPAND_COST = 60  // cost 60 stone
	N_REDO_QUEST_COST       = 60  // cost 60 stone
	N_RESUME_QUEST_COST     = 60  // cost 60 stone
	N_BUY_MONEY_COST        = 60  // cost 60 stone
	N_EVOLVE_COST           = 500 // cost 500 money
	N_FRAG_FUSION_COST      = 1000
)

// vip
const (
	V_RESUME_QUEST_DISCOUNT = 0.5
)

// shop buy
const (
	B_MONTHCARD_DAYBACK = 100 //月卡每日返还
	B_WEEKCARD_DAYBACK  = 100 //周卡每日返还
)

const (
	STEP_USER_GUIDE_GACHA int32 = 18
)
