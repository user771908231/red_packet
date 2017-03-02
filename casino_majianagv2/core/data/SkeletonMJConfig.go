package data

//麻将的核心配置
/**
	比如,血战，血流等
 */
type SkeletonMJConfig struct {
	Owner            uint32
	Password         string
	DeskId           int32
	RoomType         int32
	Status           int32
	MjRoomType       int32
	RoomId           int32
	CreateFee        int64
	BoardsCout       int32 //局数，如：4局（房卡 × 2）、8局（房卡 × 3）
	CapMax           int64
	CardsNum         int32
	Settlement       int32
	BaseValue        int64
	ZiMoRadio        int32
	OthersCheckBox   []int32
	HuRadio          int32
	DianGangHuaRadio int32
	MJPaiCursor      int32
	TotalPlayCount   int32
	CurrPlayCount    int32
	Banker           uint32
	NextBanker       uint32
	CheckCase        *CheckCase
	ActiveUser       uint32
	GameNumber       int32
	ActUser          uint32
	ActType          int32
	NInitActionTime  int32
	RoomLevel        int32
	FangCount        int32
	PlayerCountLimit int32 //玩家的人数
	BeginTime        string
	CoinLimit        int64 //金币限制
	CanChi           bool  //有些麻将可以吃，有些不可以吃(已知的长沙麻将可以吃牌)
	XueLiuChengHe    bool  // 血流成河
}
