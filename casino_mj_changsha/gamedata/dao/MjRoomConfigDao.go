package dao

import "casino_mj_changsha/gamedata/model"

//获取麻将的room配置数据
func GetMJRoomConfigData() []*model.TMjRoomConfig {
	return []*model.TMjRoomConfig{
		&model.TMjRoomConfig{
			RoomName:        "5元红包场", //房间名字
			RoomLevel:       1,       //房间等级
			RoomBaseValue:   800,     //玩牌的底分
			RoomLimitCoin:   2000,    //最近准入金币余额
			RoomLimitCoinUL: 50000,   //最近准入金币余额
			EnterCoinFee:    450,     //准入的门票
		},
		&model.TMjRoomConfig{
			RoomName:        "10元红包场",
			RoomLevel:       2, //房间等级
			RoomBaseValue:   2000,
			RoomLimitCoin:   10000,
			RoomLimitCoinUL: 250000, //最近准入金币余额
			EnterCoinFee:    1000,

		},
		&model.TMjRoomConfig{
			RoomName:        "50元红包场",
			RoomLevel:       3, //房间等级
			RoomBaseValue:   5000,
			RoomLimitCoin:   20000,
			RoomLimitCoinUL: 99999999999999, //最近准入金币余额
			EnterCoinFee:    300,

		},
		&model.TMjRoomConfig{
			RoomName:        "100元红包场",
			RoomLevel:       4, //房间等级
			RoomBaseValue:   1000,
			RoomLimitCoin:   60000,
			RoomLimitCoinUL: 99999999999999, //最近准入金币余额
			EnterCoinFee:    300,

		},
	}
}
