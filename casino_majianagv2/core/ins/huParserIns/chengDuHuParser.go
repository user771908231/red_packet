package huParserIns

import (
	"fmt"
	"casino_majianagv2/core/ins/skeleton"
	"casino_majiang/service/majiang"
	"casino_majiang/msg/protogo"
	"casino_majianagv2/core/majiangv2"
	"casino_common/common/log"
)

//成都麻将解析器
type HuParserChengDu struct {
	*skeleton.HuParserSkeleton
	BaseValue               int64 //基本分数
	IsNeedZiMoJiaDi         bool  //自摸加低
	IsNeedYaojiuJiangdui    bool  //19将对
	IsDaodaohu              bool  //倒倒胡
	IsNeedMenqingZhongzhang bool  //门清中张
	IsNeedZiMoJiaFan        bool  //自摸加翻
	CapMax                  int64 //局数
}

//得到一个成都的牌型解析器
func NewChengDuHuParser(BaseValue int64,
	IsNeedZiMoJiaDi bool,
	IsNeedYaojiuJiangdui bool,
	IsDaodaohu bool,
	IsNeedMenqingZhongzhang bool,
	IsNeedZiMoJiaFan bool,
	CapMax int64) *HuParserChengDu {
	return &HuParserChengDu{
		HuParserSkeleton:        skeleton.NewHuParserSkeleton(),
		BaseValue:               BaseValue,
		IsNeedZiMoJiaDi:         IsNeedZiMoJiaDi,
		IsNeedYaojiuJiangdui:    IsNeedYaojiuJiangdui,    //19将对
		IsDaodaohu:              IsDaodaohu,              //倒倒胡
		IsNeedMenqingZhongzhang: IsNeedMenqingZhongzhang, //门清中张
		IsNeedZiMoJiaFan:        IsNeedZiMoJiaFan,        //自摸加翻
		CapMax:                  CapMax,                  //局数
	}
}

/**
	huType 在成都麻将中暂时没有用
	return  hu,is19
 */
func (p *HuParserChengDu) GetCanHu(handPai *majiang.MJHandPai, hupai *majiang.MJPai, isZimo bool, huType mjproto.HuType) (bool, int32, int64, []string, []mjproto.PaiType, bool) {
	if handPai.IsContainQue() {
		return false, 0, 0, nil, nil, false
	}
	//在所有的牌中增加 pai,判断此牌是否能和
	pais := []*majiang.MJPai{}
	if hupai != nil {
		pais = append(pais, hupai)
	}
	if handPai.Pais != nil {
		pais = append(pais, handPai.Pais...)
	}

	counts := majiang.GettPaiStats(pais)

	var canHu bool
	var isAll19 bool
	var fan int32
	var score int64
	var huCardStr [] string
	var paiType mjproto.PaiType
	//七对 龙七对牌型 不带幺九

	canHu = p.TryHU7(counts)
	if !canHu {
		canHu, isAll19, _ = p.TryHU(counts, len(pais))
	}

	//普通33332牌型
	if canHu {
		fan, score, huCardStr, paiType = p.GetHuScore(handPai, hupai, isZimo, isAll19, huType)
	}
	paiTypeSlice := []mjproto.PaiType{}
	paiTypeSlice = append(paiTypeSlice, paiType)
	return canHu, fan, score, huCardStr, paiTypeSlice, false
}

//得到分数
func (p *HuParserChengDu) HuScore() (fan int32, score int64) {
	return 0, 0
}

//通过手牌，返回叫牌
func (p *HuParserChengDu) GetJiaoPais(pais []*majiang.MJPai) []*majiang.MJPai {
	var jiaoPais []*majiang.MJPai
	for i := 0; i < len(majiangv2.MjpaiMap); {
		tempPai := majiangv2.InitMjPaiByIndex(i)
		if p.CanHuByPais(pais, tempPai) {
			jiaoPais = append(jiaoPais, tempPai)
		}
		i += 4
	}
	return jiaoPais
}

func (p *HuParserChengDu) CanHuByPais(handPais []*majiang.MJPai, huPai *majiang.MJPai) bool {
	var pais []*majiang.MJPai
	pais = append(pais, handPais...)
	pais = append(pais, huPai)

	counts := majiang.GettPaiStats(pais)
	//七对 龙七对牌型 不带幺九
	var canHu bool ///胡牌的结果

	canHu = p.TryHU7(counts)
	if canHu {
		return canHu
	}

	//普通33332牌型
	canHu, _, _ = p.TryHU(counts, len(pais))

	return canHu
}

func (p *HuParserChengDu) GetHuScore(handPai *majiang.MJHandPai, hupai *majiang.MJPai, isZimo bool, is19 bool, extraAct mjproto.HuType) (fan int32, score int64, huCardStr [] string, paiType PaiType) {

	log.T("判断是否能胡牌的牌,手牌[%v],碰牌[%v],杠牌[%v]", handPai.GetDes(), majiangv2.ServerPais2string(handPai.PengPais), majiangv2.ServerPais2string(handPai.GangPais))

	//底分
	score = p.BaseValue

	//取得番数
	huFan, huCardStr, paiType := p.getHuFan(handPai, hupai, isZimo, is19, extraAct)

	for i := int32(0); i < huFan; i++ {
		score *= 2
	}

	//log.T("胡[%d]番 (%v)", huFan, huCardStr)

	if isZimo {
		if p.IsNeedZiMoJiaDi {
			//自摸加底
			score += p.BaseValue
		}
	}

	log.T("判断是否能胡牌的牌,手牌[%v],碰牌[%v],杠牌[%v] 结果hufan[%v],score[%v],huCardStr[%v],payType[%v]",
		handPai.GetDes(), majiangv2.ServerPais2string(handPai.PengPais), majiangv2.ServerPais2string(handPai.GangPais), huFan, score, huCardStr, paiType)

	return huFan, score, huCardStr, paiType
}

//计算带几个"勾"
func (p *HuParserChengDu) getGou(handPai *majiang.MJHandPai, handCounts [] int) (gou int32) {
	// 已杠的牌
	gou = int32(len(handPai.GangPais))
	gou = gou / 4 //杠牌/4才是gou 的数目

	//log.T("杠牌的勾:%v", gou)
	// 计算 碰牌+手牌 的勾数
	for _, pai := range handPai.Pais {
		for _, pengPai := range handPai.PengPais {
			if pai.GetClientId() == pengPai.GetClientId() {
				gou ++
				break
			}
		}
	}

	//log.T("碰牌的勾:%v", gou)

	// 计算手牌中的勾数(未暗杠)
	for _, cnt := range handCounts {
		if cnt == 4 {
			gou ++
		}
	}

	//log.T("手牌的勾:%v", gou)

	return gou
}

// 返回胡牌番数
// extraAct:指定HuPaiType.H_GangShangHua(杠上花/炮,海底等)
//
func (p *HuParserChengDu) getHuFan(handPai *majiang.MJHandPai, hupai *majiang.MJPai, isZimo bool, is19 bool, extraAct mjproto.HuType) (fan int32, huCardStr [] string, paiType mjproto.PaiType) {
	fan = int32(0)

	fanXingStr := ""
	jiaFanStr := ""
	gouStr := ""
	fengdingStr := ""
	//基本番型 勾数 总番型-加番类型(加番数)x出现次数

	pais := []*majiang.MJPai{}
	pais = append(pais, handPai.Pais...)
	pais = append(pais, hupai)

	handCounts := majiang.GettPaiStats(pais) //计算手牌的每张牌数量

	pais = append(pais, handPai.PengPais...)
	pais = append(pais, handPai.GangPais...)

	isQingYiSe := p.IsQingYiSe(pais) //清一色
	//log.T("判断是否是清一色: %v", isQingYiSe)

	isCountGou := true //是否计算勾 七对 龙七对 清龙七对 将七对 不算勾

	switch {
	case p.IsLongQiDui(handCounts): //case 清龙七对 龙七对
		//log.T("是龙七对")
		isCountGou = false
		if isQingYiSe {
			//清龙七对
			//log.T("是清龙七对")
			fan = majiangv2.FAN_QINGLONGQIDUI
			paiType = mjproto.PaiType_H_QingLongQiDui
			fanXingStr = "清龙七对"
		} else {
			//龙七对
			//log.T("是龙七对")
			fan = majiangv2.FAN_LONGQIDUI
			paiType = mjproto.PaiType_H_LongQiDui
			fanXingStr = "龙七对"
		}
	case p.IsQiDui(handCounts): //case 清七对 将七对 七对
		//log.T("是七对")
		isCountGou = false
		if isQingYiSe {
			//清七对
			//log.T("是清七对")
			fan = majiangv2.FAN_QINGQIDUI
			paiType = mjproto.PaiType_H_QingQiDui
			fanXingStr = "清七对"
		} else {
			//七对
			//log.T("是七对")
			fan = majiangv2.FAN_QIDUI
			fanXingStr = "七对"
		}
	case p.IsDaDuiZi(pais): //case 清对 将对 大对子
		//log.T("是大对子")
		if isQingYiSe {
			//清对
			//log.T("是清对")
			fan = majiang.FAN_QINGDUI
			paiType = mjproto.PaiType_H_DuiDuiHu
			fanXingStr = "清对"
		} else if p.IsNeedYaojiuJiangdui {
			//将对选项开启
			//log.T("是将对")
			if p.IsJiangDui(handPai) {
				fan = majiangv2.FAN_DADUIZI
				paiType = mjproto.PaiType_H_JiangDui
				fanXingStr = "将对"
			}
		} else {
			//大对子
			//log.T("是大对子")
			fan = majiangv2.FAN_DADUIZI
			paiType = mjproto.PaiType_H_DuiDuiHu
			fanXingStr = "大对子"
		}
	default: //default 清一色 平胡
		if isQingYiSe {
			//平胡清一色
			//log.T("是清一色")
			fan = majiangv2.FAN_QINGYISE
			paiType = mjproto.PaiType_H_QingYiSe
			fanXingStr = "清一色"
		} else {
			//平胡
			//log.T("是平胡")
			fan = majiangv2.FAN_PINGHU
			paiType = mjproto.PaiType_H_PingHu
			fanXingStr = "平胡"
		}
	}

	//倒倒胡 所有番型均算平胡
	if p.IsDaodaohu {
		//log.T("getHuFan: 倒倒胡 平胡")
		fan = majiangv2.FAN_PINGHU
		fanXingStr = "平胡"
	}

	//TODO MJOption_JIANGOUDIAO
	//if IsOpenRoomOption(roomInfo.PlayOptions.OthersCheckBox, MJOption_JINGOUDIAO) { //金钩钓
	//	if IsJingGouDiao(handCounts) {
	//		fan += FAN_JINGOUDIAO
	//		huCardStr = append(huCardStr, "金钩钓")
	//	}
	//}

	//附加选项
	if p.IsNeedYaojiuJiangdui {
		//带幺九选项开启
		if is19 && p.IsPengGang19(handPai) {
			//手牌带幺九 且 碰杠牌带幺九
			fan += majiangv2.FAN_DAIYAOJIU
			paiType = mjproto.PaiType_H_DaiYaoJiu
			jiaFanStr = "带幺九"
		}
	}

	if p.IsNeedMenqingZhongzhang {
		//门清中张选项开启
		if p.IsMenqing(handPai) {
			fan += majiangv2.FAN_MENQ_ZHONGZ
			paiType = mjproto.PaiType_H_MenQing
			jiaFanStr = "门清"
		}
		if p.IsZhongzhang(handPai, hupai) {
			fan += majiangv2.FAN_MENQ_ZHONGZ
			paiType = mjproto.PaiType_H_ZhongZhang
			jiaFanStr = "中张"
		}
	}
	isTianDiHuFlag := false //天地胡选项 避免多次搜索
	if p.IsNeedMenqingZhongzhang {
		//天地胡选项开启
		isTianDiHuFlag = true
	}

	//自摸
	if isZimo {
		if p.IsNeedZiMoJiaFan {
			fan += majiangv2.FAN_ZIMO
			jiaFanStr = fmt.Sprintf("自摸(+%d番)", majiangv2.FAN_ZIMO)
		} else if p.IsNeedZiMoJiaDi {
			//result += di
			jiaFanStr = "自摸"
		}

	}

	switch mjproto.HuType(extraAct) {

	//天地胡为牌型番数，非加番
	case mjproto.HuType_H_TianHu:
		if isTianDiHuFlag {
			//天地胡选项开启
			fan = majiangv2.FAN_TIAN_DI_HU
			fanXingStr = "天胡"
		}
	case mjproto.HuType_H_DiHu:
		if isTianDiHuFlag {
			//天地胡选项开启
			fan = majiangv2.FAN_TIAN_DI_HU
			fanXingStr = "地胡"
		}

	case mjproto.HuType_H_GangShangHua:
		fan += majiangv2.FAN_GANGSHANGHUA
		jiaFanStr = fmt.Sprintf("杠上花(+%d番)", majiangv2.FAN_GANGSHANGHUA)

	case mjproto.HuType_H_GangShangPao:
		fan += majiang.FAN_GANGSHANGPAO
		jiaFanStr = fmt.Sprintf("杠上炮(+%d番)", majiangv2.FAN_GANGSHANGPAO)

	case mjproto.HuType_H_HaiDiLao:
		fan += majiangv2.FAN_HD_LAO
		jiaFanStr = fmt.Sprintf("海底捞(+%d番)", majiangv2.FAN_HD_LAO)

	case mjproto.HuType_H_HaiDiPao:
		fan += majiangv2.FAN_HD_PAO
		jiaFanStr = fmt.Sprintf("海底炮(+%d番)", majiangv2.FAN_HD_PAO)

	case mjproto.HuType_H_QiangGang:
		fan += majiangv2.FAN_QIANGGANG
		jiaFanStr = fmt.Sprintf("抢杠(+%d番)", majiangv2.FAN_QIANGGANG)

	case mjproto.HuType_H_HaidiGangShangHua:
		fan += majiangv2.FAN_HD_GANGSHANGHUA
		jiaFanStr = fmt.Sprintf("海底杠上花(+%d番)", majiangv2.FAN_HD_GANGSHANGHUA)

	case mjproto.HuType_H_HaidiGangShangPao:
		fan += majiangv2.FAN_HD_GANGSHANGPAO
		jiaFanStr = fmt.Sprintf("海底杠上炮(+%d番)", majiangv2.FAN_HD_GANGSHANGPAO)
	default:
	}

	// 计算有几个"勾"
	if isCountGou {
		//log.T("加勾")
		gou := p.getGou(handPai, handCounts)

		fan += gou
		if gou > 0 {
			//str, _ := numUtils.Int2String(gou)
			//gouStr = append(gouStr, "勾X" + str)
			gouStr = fmt.Sprintf("勾x%d", gou)
		}
	}

	//封顶
	fanTop := p.CapMax
	if fan > int32(fanTop) {
		fan = int32(fanTop)
		fengdingStr = "封顶"
	}

	fanStr := fmt.Sprintf("%d番", fan)

	huCardStr = append(huCardStr, fanXingStr)
	huCardStr = append(huCardStr, gouStr)
	huCardStr = append(huCardStr, fanStr)
	huCardStr = append(huCardStr, jiaFanStr)
	huCardStr = append(huCardStr, fengdingStr)

	return fan, huCardStr, paiType
}

func (p *HuParserChengDu) IsPengGang19(handPai *majiang.MJHandPai) bool {
	pengPais := handPai.PengPais
	gangPais := handPai.GangPais
	if pengPais != nil {
		for i := 0; i < len(pengPais); i++ {
			if *pengPais[i].Value != 1 || *pengPais[i].Value != 9 {
				//
				return false
			}
		}
	}

	if gangPais != nil {
		for i := 0; i < len(gangPais); i++ {
			if *gangPais[i].Value != 1 || *gangPais[i].Value != 9 {
				//
				return false
			}
		}
	}
	return true
}

//大对子
func (p *HuParserChengDu) IsDaDuiZi(pais []*majiang.MJPai) bool {
	counts := majiang.GettPaiStats(pais)
	//log.T("判断是否是大对子的统计数据:%v", counts)

	jiangDui := 0

	for i := 0; i < len(counts); i++ {
		if counts[i] == 0 {
			continue
		}
		if counts[i] == 2 {
			jiangDui ++
			if jiangDui > 1 {
				//log.T("不是大对子")
				return false
			}
		} else if counts[i] == 1 {
			//log.T("不是大对子")
			return false
		}
	}
	//log.T("是大对子")
	return true
}

//将对(全是2,5,8的大对子)
func (p *HuParserChengDu) IsJiangDui(handPai *majiang.MJHandPai) bool {
	pais := handPai.Pais
	for i := 0; i < len(pais); i++ {
		if *pais[i].Value != 2 && *pais[i].Value != 5 && *pais[i].Value != 8 {
			return false
		}
	}
	return p.IsDaDuiZi(pais) //是大对子
}

//门清 没有明杠 碰牌
func (p *HuParserChengDu) IsMenqing(handPai *majiang.MJHandPai) bool {
	//TODO 明杠判断
	if handPai.PengPais != nil {
		//含碰牌
		return false
	}
	return true
}

//中张 没有1、9
func (p *HuParserChengDu) IsZhongzhang(handPai *majiang.MJHandPai, hupai *majiang.MJPai) bool {
	//
	pais := []*majiang.MJPai{}
	if handPai.Pais != nil {
		pais = append(pais, handPai.Pais...)
	}
	if hupai != nil {
		pais = append(pais, hupai)
	}

	if handPai.GangPais != nil {
		pais = append(pais, handPai.GangPais...)
	}

	if handPai.PengPais != nil {
		pais = append(pais, handPai.PengPais...)
	}

	for _, pai := range pais {
		if (pai.GetValue() == 1) || (pai.GetValue() == 9) {
			return false
		}
	}
	return true
}
