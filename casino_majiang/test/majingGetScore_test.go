package test

import (
	"testing"
	"casino_majiang/service/majiang"
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"github.com/name5566/leaf/log"
)

/**

	handPai *MJHandPai, 手牌的信息
	isZimo bool,	是否是自摸
	extraAct HuPaiType, 胡牌的类型
	roomInfo RoomTypeInfo	roomTypeInfo
	return :
		fan int32,	返回翻数
		score int64,	返回分数
		huCardStr[] string	返回胡牌的描述
 */

func TestGetScore(t *testing.T) {
	//测试自摸
	tZimo1fan()
}

func tZimo1fan() {
	isZimo := false
	var hupaiType mjproto.HuPaiType = 1;
	rfan, rscore, rhuCardStr := majiang.GetHuScore(getMjHandPai(), isZimo, hupaiType, *getRoomInfo())
	log.Debug("tZimo1fan，   番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)

}

//返回一个特性的roomTypeInfo 方便测试..
func getRoomInfo() *mjproto.RoomTypeInfo {
	info := newProto.NewRoomTypeInfo()
	*info.MjRoomType = getMjRoomType()
	*info.BoardsCout = 1        //局数
	*info.CapMax = 4        // 封顶，如：2番、3番、4番
	info.PlayOptions = getPlayOptions()
	*info.CardsNum = 13        // 牌张，如：7张、10张、13张
	*info.Settlement = 1 // 结算，如：1拖1、1拖2、3拖5
	*info.BaseValue = 1
	return info
}

func getMjRoomType() mjproto.MJRoomType {
	/**
	enum MJRoomType {
	    roomType_xueZhanDaoDi = 0; // 血战到底
	    roomType_sanRenLiangFang = 1; // 三人两房
	    roomType_siRenLiangFang = 2; // 四人两房
	    roomType_deYangMaJiang = 3; // 德阳麻将
	    roomType_daoDaoHu = 4; // 倒倒胡
	    roomType_xueLiuChengHe = 5; // 血流成河
	}
	 */
	return mjproto.MJRoomType(0)

}

func getPlayOptions() *mjproto.PlayOptions {
	/**
	message PlayOptions {
	    optional int32 ziMoRadio = 1; // 单选，自摸类型，如：自摸加底、自摸加番等
	    optional int32 dianGangHuaRadio = 2; // 单选，点炮类型，如：点杠花（点炮）、点杠花（自摸）等
	    repeated int32 othersCheckBox = 3; // 其他可复选的玩法，如：换三张、幺九将对、门清中张、天地胡、卡二条、可胡七对等
	    optional int32 huRadio = 4; // 胡法，如：自摸胡、点炮胡（可抢杠）
	}
	 */
	op := newProto.NewPlayOptions()
	return op

}


//通过设置不同的牌，来得到不同过得翻数

/**

Pais             []*MJPai `
PengPais         []*MJPai `
GangPais         []*MJPai `
HuPais           []*MJPai `
InPai            *MJPai   `
QueFlower        *int32   `
OutPais          []*MJPai `
XXX_unrecognized []byte   `
 */
func getMjHandPai() *majiang.MJHandPai {
	hand := majiang.NewMJHandPai()
	hand.GangPais = nil
	hand.HuPais = nil
	hand.InPai = nil
	hand.OutPais = nil
	hand.PengPais = nil
	*hand.QueFlower = majiang.W        //定缺的花色
	return hand
}




