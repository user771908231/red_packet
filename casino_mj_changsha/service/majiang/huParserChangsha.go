package majiang

import (
	mjproto "casino_mj_changsha/msg/protogo"
	"casino_common/common/log"
	"strings"
	"github.com/name5566/leaf/util"
)

var CHANGSHA_SCORE_XIAOHU int64 = 1 //小胡接炮 点炮输1分
var CHANGSHA_SCORE_DAHU int64 = 6   //大胡接炮 点炮输6分

type HuParserChangSha struct {
	*HuParserSkeleton
}

//得到一个长沙的牌型解析器
func NewHuParserChangSha() *HuParserChangSha {
	return &HuParserChangSha{
		HuParserSkeleton: NewHuParserSkeleton(),
	}
}

/**
	是否能胡 能胡的番与得分
	pais:手牌
	p:判断胡的牌
	huType:胡的类型（起手,自摸，杠上炮，杠上花，海底，天胡，抢杠等等）
	return  (bool, int32, int64, [] string, mjproto.PaiType
 */
func (p *HuParserChangSha) GetCanHu(handPai *MJHandPai, hupai *MJPai, zimo bool, huType mjproto.HuType, isBanker bool) (hu bool, fan int32, score int64, cardStr []string, pt []mjproto.PaiType, is258Jiang bool) {
	//起手胡 单独判断
	if huType == mjproto.HuType_H_changsha_qishouhu {
		hu, _, score, cardStr, pt = p.ChangShaCanQiShouHu(handPai, hupai)
		return hu, fan, score, cardStr, pt, is258Jiang
	}

	pingHu := false
	hu, _, score, cardStr, pt, is258Jiang, pingHu = p.ChangShaGetHuScore(handPai, hupai, zimo, isBanker)

	if !hu {
		//胡牌的牌型不满足的话 直接返回
		return false, 0, 0, nil, nil, false
	}

	//todo 抢杠胡 全求人
	//其他胡的方式

	attached := false
	switch huType {

	case mjproto.HuType_H_QiangGang://抢杠胡
		attached = true
		score = p.addDaHuScore(score)
		cardStr = append(cardStr, "抢杠")

	case mjproto.HuType_H_GangShangPao: //杠上炮
		attached = true
		score = p.addDaHuScore(score)
		cardStr = append(cardStr, "杠上炮")

	case mjproto.HuType_H_GangShangHua: //杠上花
		attached = true
		score = p.addDaHuScore(score)
		cardStr = append(cardStr, "杠上花")

	case mjproto.HuType_H_HaiDiLao: //海底捞
		attached = true
		score = p.addDaHuScore(score)
		cardStr = append(cardStr, "海底捞月")

	case mjproto.HuType_H_HaiDiPao: //海底炮
		attached = true
		score = p.addDaHuScore(score)
		cardStr = append(cardStr, "海底炮")

	case mjproto.HuType_H_TianHu: //天胡
		if is258Jiang {
			//天胡需要258的将
			attached = true
			score = p.addDaHuScore(score)
			cardStr = append(cardStr, "天胡")
		}
	case mjproto.HuType_H_DiHu: //地胡
		if is258Jiang {
			//地胡需要258的将
			attached = true
			score = p.addDaHuScore(score)
			cardStr = append(cardStr, "地胡")
		}
	default:
	}

	//如果有特殊胡牌累加 并且有平胡 不累计平胡的得分
	if attached && pingHu {
		score = p.subXiaoHuScore(score, zimo)
	}

	return hu, fan, score, cardStr, pt, is258Jiang
}

//得到叫牌的信息
func (p *HuParserChangSha) GetJiaoPais(handPai2 *MJHandPai) []*MJPai {
	log.T("长沙麻将解析器开始获取牌[%v]的叫牌", ServerPais2string(handPai2.Pais))
	var jiaoPais []*MJPai
	handPai := util.DeepClone(handPai2).(*MJHandPai)
	for i := 0; i < len(mjpaiMap); i += 4 {
		tmpPai := InitMjPaiByIndex(i)

		canHu, _, _, cardString, _, _, _ := p.ChangShaGetHuScore(handPai, tmpPai, false, false)

		if canHu {
			log.T("[%v]可胡并加入到叫牌数组, 胡类型[%v]", tmpPai.GetDes(), strings.Join(cardString, ""))
			jiaoPais = append(jiaoPais, tmpPai)
		}
	}
	log.T("长沙麻将解析器完成获取叫牌, 叫牌为 %v", ServerPais2string(jiaoPais))
	return jiaoPais
}

//得到分数
func (p *HuParserChangSha) HuScore() (fan int32, score int64) {
	return 0, 0
}

//累加大胡的得分 庄家要多加一分
func (p *HuParserChangSha) addDaHuScore(score int64) int64 {
	score += CHANGSHA_SCORE_DAHU
	return score
}

//累加小胡的得分 小胡自摸要多加一分
func (p *HuParserChangSha) addXiaoHuScore(score int64, zimo bool) int64 {
	if zimo {
		score++
	}
	score += CHANGSHA_SCORE_XIAOHU
	return score
}

func (p *HuParserChangSha) subXiaoHuScore(score int64, zimo bool) int64 {
	if zimo {
		score--
	}
	score -= CHANGSHA_SCORE_XIAOHU
	if score < 0 {
		score = 0
	}
	return score
}

//起手胡的四种牌型 大四喜 板板胡 缺一色 六六顺
func (p *HuParserChangSha) ChangShaCanQiShouHu(handPai *MJHandPai, hupai *MJPai) (hu bool, fan int32, score int64, cardStr []string, pt []mjproto.PaiType) {
	//log.T("长沙麻将解析器 开始判断是否是起手胡")
	if p.isYouChi(handPai) || p.isYouPeng(handPai) || p.isYouGang(handPai) {
		log.T("碰杠吃牌不为空 起手胡错误")
		return false, fan, score, cardStr, pt
	}

	//小胡
	if p.ChangShaIsSiXi(handPai, hupai) {
		//四喜
		hu = true
		score = p.addXiaoHuScore(score, true)
		cardStr = append(cardStr, "四喜")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_DAXISI)
	}
	if p.ChangShaIsBanBanHu(handPai, hupai) {
		//板板胡
		hu = true
		score = p.addXiaoHuScore(score, true)
		cardStr = append(cardStr, "板板胡")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_BANBANHU)
	}
	if p.ChangShaIsQueYiMen(handPai, hupai) {
		//缺一门
		hu = true
		score = p.addXiaoHuScore(score, true)
		cardStr = append(cardStr, "缺一色")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QUEYISE)
	}
	if p.ChangShaIsLiuLiuHu(handPai, hupai) {
		//六六顺
		hu = true
		score = p.addXiaoHuScore(score, true)
		cardStr = append(cardStr, "六六顺")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_LIULIUSHUN)
	}

	return hu, fan, score, cardStr, pt
}

//
func (p *HuParserChangSha) ChangShaGetHuScore(handPai *MJHandPai, hupai *MJPai, zimo bool, isBanker bool) (hu bool, fan int32, score int64, cardStr []string, pt []mjproto.PaiType, is258Jiang bool, pingHu bool) {

	//log.T("开始判断胡牌类型 判定牌[%v]", hupai.GetDes())

	pais := p.JoinHandPaiPaisAndHuPai(handPai, hupai)

	handCounts := GettPaiStats(pais)

	//判断各牌型是否满足
	canHu, _, jiang := p.tryHU(handCounts, len(pais))
	is258Jiang = p.is258(jiang)
	//log.T("长沙平胡 tryHu[%v] ,jiang [%v],258jiang[%v]", canHu, jiang, is258Jiang)

	qingYiSe := p.ChangShaIsQingYiSe(handPai, hupai)
	isQingYiSeCounted := false

	ppHu := p.ChangShaIsPengPengHu(handPai, hupai)
	jjHu := p.ChangShaIsJiangJiangHu(handPai, hupai)

	qiDui := p.ChangShaIsQiXiaoDui(handPai, hupai)
	qiDuiHaoHua := p.ChangShaIsQiDuiHaoHua(handPai, hupai)
	qiDuiHaoHuaDouble := p.ChangShaIsQiDuiHaoHuaDouble(handPai, hupai)

	quanQiuRen := p.ChangShaIsQuanQiuRen(handPai, hupai)

	//统计分数/牌型
	if ppHu {
		//碰碰胡
		//log.T("碰碰胡")
		hu = true
		score = p.addDaHuScore(score)
		cardStr = append(cardStr, "碰碰胡")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_PENGPENGHU)
	}
	if jjHu {
		//将将胡
		//log.T("将将胡")
		hu = true
		score = p.addDaHuScore(score)
		cardStr = append(cardStr, "将将胡")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_JIANGJIANGHU)
	}

	switch {
	case qiDuiHaoHuaDouble:
		//双豪华七小对
		//log.T("双豪华七小对")
		hu = true
		score = p.addDaHuScore(score)
		score = p.addDaHuScore(score)
		score = p.addDaHuScore(score)
		score = p.addDaHuScore(score)

		cardStr = append(cardStr, "双豪华七小对")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QIXIAODUI_HAOHUA_DOUBLE)
	case qiDuiHaoHua:
		//豪华七小对
		//log.T("豪华七小对")
		hu = true
		score = p.addDaHuScore(score)
		score = p.addDaHuScore(score)

		cardStr = append(cardStr, "豪华七小对")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QIXIAODUI_HAOHUA)
	case qiDui:
		//七对
		//log.T("七对")
		hu = true
		score = p.addDaHuScore(score)
		cardStr = append(cardStr, "七小对")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QIXIAODUI)
	default:

	}

	switch {
	case quanQiuRen:
		//全求人
		//log.T("全求人")
		hu = true
		score = p.addDaHuScore(score)

		cardStr = append(cardStr, "全求人")

		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QUANQIUREN)
	case canHu && is258Jiang && !qingYiSe:
		pingHu = true
		//长沙平胡
		//log.T("平胡")
		hu = true
		score = p.addXiaoHuScore(score, zimo)
		cardStr = append(cardStr, "平胡")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_PINGHU)

	case canHu && is258Jiang && qingYiSe: //如果是清一色的平胡 就算清一色
		//log.T("清一色")
		hu = true
		isQingYiSeCounted = true
		score = p.addDaHuScore(score)
		cardStr = append(cardStr, "清一色")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QINGYISE)
	default:
	}

	//清一色可以和大胡叠加 可以不用258做将胡
	//注: 清一色平胡的情况上面已覆盖 这里需要排出
	if (canHu || ppHu || jjHu || qiDui || qiDuiHaoHua || qiDuiHaoHuaDouble || quanQiuRen) && qingYiSe && !isQingYiSeCounted {
		//清一色 满足任意一个牌型即可
		//log.T("清一色")
		hu = true
		score = p.addDaHuScore(score)
		cardStr = append(cardStr, "清一色")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QINGYISE)
	}

	return hu, fan, score, cardStr, pt, is258Jiang, pingHu
}

/******************** 长沙麻将 小胡 ********************/

//是否是大四喜 4张一样的牌
func (p *HuParserChangSha) ChangShaIsSiXi(handPai *MJHandPai, hupai *MJPai) bool {
	joinedPais := p.JoinHandPaiPaisAndHuPai(handPai, hupai)
	handCounts := GettPaiStats(joinedPais)
	for i := 0; i < len(handCounts); i++ {
		if handCounts[i] == 4 {
			return true
		}
	}
	return false
}

//是否是板板胡 没有一张 2、5、8
func (p *HuParserChangSha) ChangShaIsBanBanHu(handPai *MJHandPai, hupai *MJPai) bool {
	joinedPais := p.JoinAllHandPaiAndHuPai(handPai, hupai)
	for i := 0; i < len(joinedPais); i++ {
		if *joinedPais[i].Value == 2 ||
			*joinedPais[i].Value == 5 ||
			*joinedPais[i].Value == 8 {
			//只要包含2 5 8 就不是板板胡
			return false
		}
	}
	return true
}

//是否是缺一门 筒索万任缺一门
func (p *HuParserChangSha) ChangShaIsQueYiMen(handPai *MJHandPai, hupai *MJPai) bool {
	joinedPais := p.JoinAllHandPaiAndHuPai(handPai, hupai)
	t, s, w := 0, 0, 0
	for i := 0; i < len(joinedPais); i++ {
		switch *joinedPais[i].Flower {
		case T:
			t++
		case S:
			s++
		case W:
			w++
		}
	}
	return t == 0 || s == 0 || w == 0
}

//是否是六六胡 至少两副刻子
func (p *HuParserChangSha) ChangShaIsLiuLiuHu(handPai *MJHandPai, hupai *MJPai) bool {
	joinedPais := p.JoinAllHandPaiAndHuPai(handPai, hupai)
	handCounts := GettPaiStats(joinedPais)
	count := 0
	for i := 0; i < len(handCounts); i++ {
		if handCounts[i] == 3 {
			count++
		}
	}

	if count >= 2 {
		return true
	}
	return false
}

/******************** 长沙麻将 大胡 ********************/
//是否是碰碰胡
func (p *HuParserChangSha) ChangShaIsPengPengHu(handPai *MJHandPai, hupai *MJPai) bool {
	if p.isYouChi(handPai) {
		//log.T("碰碰胡不能有吃牌 碰碰胡失败")
		return false
	}

	joinedPais := p.JoinHandPaiPaisAndHuPai(handPai, hupai)
	handCounts := GettPaiStats(joinedPais)
	//log.T("判断是否是碰碰胡的统计数据:%v", handCounts)
	jiangDui := 0

	for i := 0; i < len(handCounts); i++ {
		if handCounts[i] == 0 {
			continue
		}
		if handCounts[i] == 2 {
			jiangDui ++
			if jiangDui > 1 {
				//log.T("牌值为2的牌大于1对 不是碰碰胡")
				return false
			}
		} else if handCounts[i] != 3 {
			//log.T("牌数[%v]不为3 不是碰碰胡", handCounts[i])
			return false
		}
	}
	//log.T("是碰碰胡")
	return true
}

//是否是将将胡 没有一张 2、5、8
func (p *HuParserChangSha) ChangShaIsJiangJiangHu(handPai *MJHandPai, hupai *MJPai) bool {
	//log.T("判断是否是将将胡")
	pais := p.JoinAllHandPaiAndHuPai(handPai, hupai)
	for i := 0; i < len(pais); i++ {
		if *pais[i].Value != 2 &&
			*pais[i].Value != 5 &&
			*pais[i].Value != 8 {
			//只要不包含2 5 8 就不是将将胡
			//log.T("牌值[%v]不是2 5 8, 不是将将胡", *pais[i].Value)
			return false
		}
	}
	return true
}

//是否是清一色 这里用majiang的判断清一色方法
func (p *HuParserChangSha) ChangShaIsQingYiSe(handPai *MJHandPai, hupai *MJPai) bool {
	pais := p.JoinAllHandPaiAndHuPai(handPai, hupai)
	return p.IsQingYiSe(pais)
}

//是否是全求人
func (p *HuParserChangSha) ChangShaIsQuanQiuRen(handPai *MJHandPai, hupai *MJPai) bool {
	joinedPais := p.JoinHandPaiPaisAndHuPai(handPai, hupai)
	if joinedPais != nil && len(joinedPais) != 2 {
		//只有两张牌 才是全求人
		return false
	}

	handCounts := GettPaiStats(joinedPais)
	for i := 0; i < len(handCounts); i++ {
		if handCounts[i] == 0 {
			continue
		}
		if handCounts[i] == 2 {
			//相同的牌有两张 单调 即全求人
			return true
		}
	}
	return false
}

//是否是七小对
func (p *HuParserChangSha) ChangShaIsQiXiaoDui(handPai *MJHandPai, hupai *MJPai) bool {
	isQidui, gou := p.IsQiDui(handPai, hupai)
	return isQidui && gou == 0
}

//是否是豪华七对
func (p *HuParserChangSha) ChangShaIsQiDuiHaoHua(handPai *MJHandPai, hupai *MJPai) bool {
	isQidui, gou := p.IsQiDui(handPai, hupai)
	return isQidui && gou == 1
}

//是否是双豪华七对
func (p *HuParserChangSha) ChangShaIsQiDuiHaoHuaDouble(handPai *MJHandPai, hupai *MJPai) bool {
	isQidui, gou := p.IsQiDui(handPai, hupai)
	return isQidui && gou == 2
}
