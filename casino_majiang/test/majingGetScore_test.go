package test

import (
	"testing"
	"casino_majiang/service/majiang"
	"casino_majiang/msg/protogo"
	"casino_majiang/msg/funcsInit"
	"github.com/name5566/leaf/log"
	"time"
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
	//tZimo1fan()

	//测试CanHuPai效率
	//tCanHuPaiPerformance()

	//测试CanHuPai
	//tCanHuPai()

	//测试平胡
	//tPinghu()
	//////
	//////测试对对胡
	//tDuiDuihu()
	////
	//////测试清一色
	//tQingyise()
	//
	//////测试带幺九
	//tDaiyaojiu()

	//
	////测试七对
	//tQidui()
	//
	////测试龙七对
	//tLongqidui()
	////
	////测试清对
	//tQingdui()

	//测试将对
	//tJiangdui()

	//测试将七对
	//tJIangqidui()

	//测试清七对
	//tQingqidui()
	//
	////测试青龙七对
	//tQinglongqidui()
}

//func tZimo1fan() {
//	isZimo := true
//	var hupaiType mjproto.HuPaiType = 1;
//	rfan, rscore, rhuCardStr := majiang.GetHuScore(getMjHandPai(), isZimo, hupaiType, *getRoomInfo())
//	if rfan ! = 1 {
//		fmt.Println(" error :")
//	}
//	log.Debug("tZimo1fan，   番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
//}

func tCanHuPaiPerformance() {
	var canHu bool
	var beginTime, endTime time.Time
	var totalTime time.Duration

	beginTime = time.Now()
	var loop int = 100000
	log.Debug("%v次循环开始", loop)
	for i := 0; i < loop; i++ {
		canHu, _ = majiang.CanHuPai(getPinghu())
	}
	endTime = time.Now()
	totalTime = endTime.Sub(beginTime)

	log.Debug("%v次循环结束, 总耗时%v, 平均耗时%v", loop, totalTime.Seconds(), time.Duration(int(totalTime.Seconds()) / loop))

	log.Debug("－－－－牌型CanHuPai耗时测试开始－－－－")
	endTime = time.Now()
	beginTime = time.Now()
	canHu, _ = majiang.CanHuPai(getQingdui())
	log.Debug("清对的牌型CanHuPai: %v", canHu)
	endTime = time.Now()
	log.Debug("清对牌型CanHuPai耗时:%v", endTime.Sub(beginTime))

	beginTime = time.Now()
	canHu, _ = majiang.CanHuPai(getQidui())
	log.Debug("七对的牌型CanHuPai: %v", canHu)
	endTime = time.Now()
	log.Debug("七对牌型CanHuPai耗时:%v", endTime.Sub(beginTime))

	beginTime = time.Now()
	canHu, _ = majiang.CanHuPai(getJiangdui())
	log.Debug("将对的牌型CanHuPai: %v", canHu)
	endTime = time.Now()
	log.Debug("将对牌型CanHuPai耗时:%v", endTime.Sub(beginTime))

	beginTime = time.Now()
	canHu, _ = majiang.CanHuPai(getLongqidui())
	log.Debug("龙七对的牌型CanHuPai: %v", canHu)
	endTime = time.Now()
	log.Debug("龙七对牌型CanHuPai耗时:%v", endTime.Sub(beginTime))

	beginTime = time.Now()
	canHu, _ = majiang.CanHuPai(getDaiyaojiu())
	log.Debug("带幺九的牌型CanHuPai: %v", canHu)
	endTime = time.Now()
	log.Debug("带幺九牌型CanHuPai耗时:%v", endTime.Sub(beginTime))

	beginTime = time.Now()
	canHu, _ = majiang.CanHuPai(getDuiduihu())
	log.Debug("对对胡的牌型CanHuPai: %v", canHu)
	endTime = time.Now()
	log.Debug("对对胡牌型CanHuPai耗时:%v", endTime.Sub(beginTime))

	beginTime = time.Now()
	canHu, _ = majiang.CanHuPai(getPinghu())
	log.Debug("平胡的牌型CanHuPai: %v", canHu)
	endTime = time.Now()
	log.Debug("平胡牌型CanHuPai耗时:%v", endTime.Sub(beginTime))

	beginTime = time.Now()
	canHu, _ = majiang.CanHuPai(getQinglongqidui())
	log.Debug("青龙七对的牌型CanHuPai: %v", canHu)
	endTime = time.Now()
	log.Debug("清龙七对牌型CanHuPai耗时:%v", endTime.Sub(beginTime))

	beginTime = time.Now()
	canHu, _ = majiang.CanHuPai(getQingqidui())
	log.Debug("清七对的牌型CanHuPai: %v", canHu)
	endTime = time.Now()
	log.Debug("清七对胡牌型CanHuPai耗时:%v", endTime.Sub(beginTime))

	beginTime = time.Now()
	canHu, _ = majiang.CanHuPai(getQingyaojiu())
	log.Debug("清幺九的牌型CanHuPai: %v", canHu)
	endTime = time.Now()
	log.Debug("清幺九胡牌型CanHuPai耗时:%v", endTime.Sub(beginTime))

	beginTime = time.Now()
	canHu, _ = majiang.CanHuPai(getQingyise())
	log.Debug("清一色的牌型CanHuPai: %v", canHu)
	endTime = time.Now()
	log.Debug("清一色胡牌型CanHuPai耗时:%v", endTime.Sub(beginTime))

	log.Debug("－－－－牌型CanHuPai耗时测试结束－－－－")
}

func tCanHuPai() {
	var canHu bool
	canHu, _ = majiang.CanHuPai(getQingdui())
	log.Debug("清对的牌型CanHuPai: %v", canHu)

	canHu, _ = majiang.CanHuPai(getQidui())
	log.Debug("七对的牌型CanHuPai: %v", canHu)

	canHu, _ = majiang.CanHuPai(getJiangdui())
	log.Debug("将对的牌型CanHuPai: %v", canHu)

	canHu, _ = majiang.CanHuPai(getLongqidui())
	log.Debug("龙七对的牌型CanHuPai: %v", canHu)

	canHu, _ = majiang.CanHuPai(getDaiyaojiu())
	log.Debug("带幺九的牌型CanHuPai: %v", canHu)

	canHu, _ = majiang.CanHuPai(getDuiduihu())
	log.Debug("对对胡的牌型CanHuPai: %v", canHu)

	canHu, _ = majiang.CanHuPai(getPinghu())
	log.Debug("平胡的牌型CanHuPai: %v", canHu)

	canHu, _ = majiang.CanHuPai(getQinglongqidui())
	log.Debug("青龙七对的牌型CanHuPai: %v", canHu)

	canHu, _ = majiang.CanHuPai(getQingqidui())
	log.Debug("清七对的牌型CanHuPai: %v", canHu)

	canHu, _ = majiang.CanHuPai(getQingyaojiu())
	log.Debug("清幺九的牌型CanHuPai: %v", canHu)

	canHu, _ = majiang.CanHuPai(getQingyise())
	log.Debug("清幺九的牌型CanHuPai: %v", canHu)

}

func tPinghu() {
	isZimo := true
	var hupaiType mjproto.HuPaiType = 1;
	rfan, rscore, rhuCardStr := majiang.GetHuScore(getPinghu(), isZimo,false, hupaiType, *getRoomInfo())
	if rfan != majiang.FAN_PINGHU {
		log.Debug("error: 平胡")
	}
	log.Debug("平胡: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
}

func tDuiDuihu() {
	isZimo := true
	var hupaiType mjproto.HuPaiType = 1;
	rfan, rscore, rhuCardStr := majiang.GetHuScore(getDuiduihu(), isZimo,false, hupaiType, *getRoomInfo())
	if rfan != majiang.FAN_DADUIZI {
		log.Debug("error: 对对胡")
	}
	log.Debug("对对胡: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
}

func tQingyise() {
	isZimo := true
	var hupaiType mjproto.HuPaiType = 1;
	rfan, rscore, rhuCardStr := majiang.GetHuScore(getQingyise(), isZimo,false, hupaiType, *getRoomInfo())
	if rfan != majiang.FAN_QINGYISE {
		log.Debug("error: 清一色")
	}
	log.Debug("清一色: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
}

func tDaiyaojiu() {
	//isZimo := true
	//var hupaiType mjproto.HuPaiType = 1
	var isAll19 bool
	_, isAll19 = majiang.CanHuPai(getDaiyaojiu())
	log.Debug("带幺九的牌型判定带幺九: %v", isAll19)

	//_, isAll19 = majiang.CanHuPai(getPinghu())
	//log.Debug("平胡的牌型判定带幺九: %v", isAll19)
	//
	//_, isAll19 = majiang.CanHuPai(getDuiduihu())
	//log.Debug("对对胡的牌型判定带幺九: %v", isAll19)

	//_, isAll19 = majiang.CanHuPai(getQingyise())
	//log.Debug("清一色的牌型判定带幺九: %v", isAll19)

	//_, isAll19 = majiang.CanHuPai(getJiangdui())
	//log.Debug("将对的牌型判定带幺九: %v", isAll19)

	_, isAll19 = majiang.CanHuPai(getLongqidui())
	log.Debug("龙七对的牌型判定带幺九: %v", isAll19)

	_, isAll19 = majiang.CanHuPai(getQidui())
	log.Debug("七对的牌型判定带幺九: %v", isAll19)

	//_, isAll19 = majiang.CanHuPai(getQingdui())
	//log.Debug("清对的牌型判定带幺九: %v", isAll19)

	//rfan, rscore, rhuCardStr := majiang.GetHuScore(getDaiyaojiu(), isZimo,false, hupaiType, *getRoomInfo())
	//if rfan != majiang.FAN_DAIYAOJIU {
	//	log.Debug("error: 带幺九")
	//}
	//log.Debug("带幺九: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
}

func tQidui() {
	isZimo := true
	var hupaiType mjproto.HuPaiType = 1;
	rfan, rscore, rhuCardStr := majiang.GetHuScore(getQidui(), isZimo,false, hupaiType, *getRoomInfo())
	if rfan != majiang.FAN_QIDUI {
		log.Debug("error: 七对")
	}
	log.Debug("七对: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
}

func tQingdui() {
	isZimo := true
	var hupaiType mjproto.HuPaiType = 1;
	rfan, rscore, rhuCardStr := majiang.GetHuScore(getQingdui(), isZimo,false, hupaiType, *getRoomInfo())
	if rfan != majiang.FAN_QINGDUI {
		log.Debug("error: 清对")
	}
	log.Debug("清对: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
}

//func tJiangqidui() {
//	isZimo := true
//	var hupaiType mjproto.HuPaiType = 1;
//	rfan, rscore, rhuCardStr := majiang.GetHuScore(getJiangqidui(), isZimo,false, hupaiType, *getRoomInfo())
//	if rfan != majiang.FAN_JIANGQIDUI {
//		log.Debug("error: 将七对")
//	}
//	log.Debug("将七对: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
//}

func tJiangdui() {
	isZimo := true
	var hupaiType mjproto.HuPaiType = 1;
	rfan, rscore, rhuCardStr := majiang.GetHuScore(getJiangdui(), isZimo,false, hupaiType, *getRoomInfo())
	if rfan != majiang.FAN_JIANGDUI {
		log.Debug("error: 将对")
	}
	log.Debug("将对: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
}

func tLongqidui() {
	isZimo := true
	var hupaiType mjproto.HuPaiType = 1;
	rfan, rscore, rhuCardStr := majiang.GetHuScore(getLongqidui(), isZimo,false, hupaiType, *getRoomInfo())
	if rfan != majiang.FAN_LONGQIDUI {
		log.Debug("error: 龙七对")
	}
	log.Debug("龙七对: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
}

func tQingqidui() {
	isZimo := true
	var hupaiType mjproto.HuPaiType = 1;
	rfan, rscore, rhuCardStr := majiang.GetHuScore(getQingqidui(), isZimo,false, hupaiType, *getRoomInfo())
	if rfan != majiang.FAN_QINGQIDUI {
		log.Debug("error: 清七对")
	}
	log.Debug("清七对: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
}

func tQingyaojiu() {
	isZimo := true
	var hupaiType mjproto.HuPaiType = 1;
	rfan, rscore, rhuCardStr := majiang.GetHuScore(getQingyaojiu(), isZimo,false, hupaiType, *getRoomInfo())
	if rfan != majiang.FAN_QINGYAOJIU {
		log.Debug("error: 清幺九")
	}
	log.Debug("清幺九: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
}

func tQinglongqidui() {
	isZimo := true
	var hupaiType mjproto.HuPaiType = 1;
	rfan, rscore, rhuCardStr := majiang.GetHuScore(getQinglongqidui(), isZimo,false, hupaiType, *getRoomInfo())
	if rfan != majiang.FAN_QINGLONGQIDUI {
		log.Debug("error: tQinglongqidui")
	}
	log.Debug("青龙七对: 番数[%v],得分[%v],描述[%v]", rfan, rscore, rhuCardStr)
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
func getMjHandPai(inPaiDes string, pengPaisDes []string, gangPaisDes []string, paisDes []string) *majiang.MJHandPai {
	hand := majiang.NewMJHandPai()

	hand.InPai = majiang.InitMjPaiByDes(inPaiDes, hand)

	for i := 0; i < len(paisDes); i++ {
		hand.Pais = append(hand.Pais, majiang.InitMjPaiByDes(paisDes[i], hand))
	}
	for i := 0; i < len(pengPaisDes); i++ {
		hand.PengPais = append(hand.PengPais, majiang.InitMjPaiByDes(pengPaisDes[i], hand))
	}
	for i := 0; i < len(gangPaisDes); i++ {
		hand.GangPais = append(hand.GangPais, majiang.InitMjPaiByDes(gangPaisDes[i], hand))
	}


	//ignore params
	hand.HuPais = nil
	hand.OutPais = nil
	*hand.QueFlower = majiang.W        //定缺的花色
	return hand
}

func getNoHu() *majiang.MJHandPai {
	inPaiDes	:= "T_4"
	paisDes		:= []string{"S_6", "S_7", "S_8", "T_7", "T_7", "T_9", "T_4", "T_4", "T_4"}
	pengPaisDes	:= []string{"T_1", "T_1", "T_1"}
	gangPaisDes	:= []string{"T_3", "T_3", "T_3", "T_3"}

	return getMjHandPai(inPaiDes, pengPaisDes, gangPaisDes, paisDes)
}

//平胡 1番
func getPinghu() *majiang.MJHandPai {
	inPaiDes	:= "T_2"
	paisDes		:= []string{"S_6", "S_7", "S_8", "T_7", "T_8", "T_9", "T_4", "T_4", "T_4", "T_2"}
	pengPaisDes	:= []string{"T_1", "T_1", "T_1"}
	gangPaisDes	:= []string{"T_3", "T_3", "T_3", "T_3"}

	return getMjHandPai(inPaiDes, pengPaisDes, gangPaisDes, paisDes)

}

//对对胡 2番
func getDuiduihu() *majiang.MJHandPai {
	inPaiDes	:= "T_4" //4T
	paisDes		:= []string{"S_6", "S_6", "S_6", "T_7", "T_7", "T_7", "T_4"} //666S 777T 4T
	pengPaisDes	:= []string{"T_2", "T_2", "T_2"} //111T 222T
	gangPaidsDes	:= []string{"T_1", "T_1", "T_1", "T_1"}
	return getMjHandPai(inPaiDes, pengPaisDes, gangPaidsDes, paisDes)
}

//清一色 3番
func getQingyise() *majiang.MJHandPai {
	inPaiDes	:= "T_9" //4T
	paisDes		:= []string{"T_2", "T_3", "T_4", "T_5", "T_6", "T_7", "T_4", "T_4", "T_9", "T_9"} //123T 456T 4T
	gangPaisDes	:= []string{"T_1", "T_1", "T_1", "T_1"}
	return getMjHandPai(inPaiDes, nil, gangPaisDes, paisDes)
}

//带幺九 3番
func getDaiyaojiu() *majiang.MJHandPai {
	//inPaiDes	:= "S_1" //1S
	//paisDes		:= []string{"S_1", "S_2", "S_3", "S_9", "S_9", "S_9", "S_1"} //123S 999S 1S
	//pengPaisDes	:= []string{"T_1", "T_1", "T_1", "T_9", "T_9", "T_9"} //111T 999T

	inPaiDes	:= "S_1" //1S
	paisDes		:= []string{"S_1", "S_1", "S_2", "S_2", "S_2", "S_3", "S_3", "S_3"}
	pengPaisDes	:= []string{"T_1", "T_1", "T_1", "T_9", "T_9", "T_9"}

	//inPaiDes	:= "S_1" //1S
	//paisDes		:= []string{"S_1", "S_2", "S_3", "S_8", "S_8", "S_8", "S_1"} //123S 999S 1S
	//pengPaisDes	:= []string{"T_1", "T_1", "T_1", "T_9", "T_9", "T_9"} //111T 999T

	return getMjHandPai(inPaiDes, pengPaisDes, nil, paisDes)
}

//七对 3番
func getQidui() *majiang.MJHandPai {
	inPaiDes	:= "S_6" //6S
	paisDes		:= []string{"S_1", "S_1", "S_2", "S_2", "S_4", "S_4", "T_9", "T_9", "T_7", "T_7", "S_7", "S_7", "S_6"} //11S 22S 44S 99T 77T 6S
	return getMjHandPai(inPaiDes, nil, nil, paisDes)
}

//将七对 3番
//func getJiangqidui() *majiang.MJHandPai {
//	inPaiDes	:= "S_6" //6S
//	paisDes		:= []string{"T_2", "T_2", "S_2", "S_2", "S_5", "S_5", "T_8", "T_8", "T_5", "T_5", "S_8", "S_8", "S_6"} //11S 22S 44S 99T 77T 6S
//	return getMjHandPai(inPaiDes, nil, nil, paisDes)
//}

//清对 4番
func getQingdui() *majiang.MJHandPai {
	inPaiDes	:= "T_9" //9T
	paisDes		:= []string{"T_5", "T_5", "T_5", "T_7", "T_7", "T_7", "T_9"} //555T 777T 9T
	pengPaisDes	:= []string{"T_1", "T_1", "T_1", "T_3", "T_3", "T_3"} //111T 333T
	return getMjHandPai(inPaiDes, pengPaisDes, nil, paisDes)
}

//将对 4番
func getJiangdui() *majiang.MJHandPai {
	inPaiDes	:= "S_2" //2S
	paisDes		:= []string{"S_5", "S_5", "S_5", "S_8", "S_8", "S_8", "S_2"} //555S 888S 2S
	pengPaisDes	:= []string{"T_2", "T_2", "T_2", "T_5", "T_5", "T_5"} //222T 555T
	gangPaisDes	:= []string{"T_8", "T_8", "T_8", "T_8"}
	return getMjHandPai(inPaiDes, pengPaisDes, gangPaisDes, paisDes)
}

//龙七对 5番
func getLongqidui() *majiang.MJHandPai {
	inPaiDes	:= "T_7" //7T
	paisDes		:= []string{"S_1", "S_1", "S_2", "S_2", "S_4", "S_4", "T_9", "T_9", "T_6", "T_6", "T_7", "T_7", "T_7"} //11S 22S 44S 99T 777T
	return getMjHandPai(inPaiDes, nil, nil, paisDes)
}

//清七对 5番
func getQingqidui() *majiang.MJHandPai {
	inPaiDes	:= "S_7" //7S
	paisDes		:= []string{"S_1", "S_1", "S_2", "S_2", "S_3", "S_3", "S_5", "S_5", "S_6", "S_6", "S_9", "S_9", "S_7"} //11S 22S 55S 66S 7S
	return getMjHandPai(inPaiDes, nil, nil, paisDes)
}


//清幺九 5番
func getQingyaojiu() *majiang.MJHandPai {
	inPaiDes	:= "S_4" //4S
	paisDes	:= []string{"S_1", "S_2", "S_3", "S_1", "S_2", "S_3", "S_1", "S_2", "S_3", "S_4"} //123S 123S 123S 4S
	pengPaisDes	:= []string{"S_9", "S_9", "S_9"} //999S
	return getMjHandPai(inPaiDes, pengPaisDes, nil, paisDes)
}

//清龙七对 6番
func getQinglongqidui() *majiang.MJHandPai {
	inPaiDes	:= "T_7" //7T
	paisDes		:= []string{"T_1", "T_1", "T_2", "T_2", "T_3", "T_3", "T_4", "T_4", "T_5", "T_5", "T_7", "T_7", "T_7"} //11T 22T 44T 777T
	return getMjHandPai(inPaiDes, nil, nil, paisDes)
}

