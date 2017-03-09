package agentModel

import (

)
import (
	"casino_common/proto/ddproto"
	"github.com/golang/protobuf/proto"
)

func init() {
	Rebates = []RebateItem{
		//卖满10张房卡则奖励1张房卡
		RebateItem{
			Num:10,
			Reward: []*ddproto.HallBagItem{
				&ddproto.HallBagItem{
					Type: ddproto.HallEnumTradeType_PROPS_FANGKA.Enum(),
					Amount:proto.Float64(1),
				},
			},
		},
		//卖满50张房卡则奖励10张房卡
		RebateItem{
			Num:50,
			Reward: []*ddproto.HallBagItem{
				&ddproto.HallBagItem{
					Type: ddproto.HallEnumTradeType_PROPS_FANGKA.Enum(),
					Amount:proto.Float64(10),
				},
			},
		},
	}
}

type RebateItem struct {
	Num int64  //数量
	Reward []*ddproto.HallBagItem
}

var Rebates []RebateItem

//获取可得的返利房卡数
func GetRebateRoomCardNum(num int64) (count int64) {
	for i,_ := range Rebates {
		item := Rebates[-i]
		if item.Num <= num && len(item.Reward) > 0 {
			return int64(item.Reward[0].GetAmount())
		}
	}
	return 0
}
