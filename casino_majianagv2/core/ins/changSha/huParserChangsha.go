package changSha

import (
	"casino_majiang/msg/protogo"
	"casino_common/common/log"
	"casino_majiang/service/majiang"
	"casino_majianagv2/core/majiangv2"
	"casino_majianagv2/core/ins/skeleton"
)

var CHANGSHA_SCORE_XIAOHU_ZIMO int64 = 2 //小胡自摸 每家2 * 3
var CHANGSHA_SCORE_XIAOHU int64 = 1      //小胡接炮 点炮输1分
var CHANGSHA_SCORE_DAHU_ZIMO int64 = 6   //大胡自摸 每家6 * 3
var CHANGSHA_SCORE_DAHU int64 = 6        //大胡接炮 点炮输6分

type HuParserChangSha struct {
	*skeleton.HuParserSkeleton
}

//得到一个长沙的牌型解析器
func NewHuParserChangSha() *HuParserChangSha {
	return &HuParserChangSha{
		HuParserSkeleton: skeleton.NewHuParserSkeleton(),
	}
}

/**
	是否能胡 能胡的番与得分
	pais:手牌
	p:判断胡的牌
	huType:胡的类型（起手,自摸，杠上炮，杠上花，海底，天胡，抢杠等等）
	return  (bool, int32, int64, [] string, mjproto.PaiType
 */
func (p *HuParserChangSha) GetCanHu(handPai *majiang.MJHandPai, hupai *majiang.MJPai, zimo bool, huType mjproto.HuType) (hu bool, fan int32, score int64, cardStr []string, pt []mjproto.PaiType, is258Jiang bool) {
	//起手胡 单独判断
	if huType == mjproto.HuType_H_changsha_qishouhu {
		hu, _, score, cardStr, pt = p.ChangShaCanQiShouHu(handPai, hupai)
		return hu, fan, score, cardStr, pt, is258Jiang
	}

	pingHu := false
	hu, _, score, cardStr, pt, is258Jiang, pingHu = p.ChangShaGetHuScore(handPai, hupai, zimo)

	if !hu {
		//胡牌的牌型不满足的话 直接返回
		return false, 0, 0, nil, nil, false
	}

	//todo 抢杠胡 全求人
	//其他胡的方式

	attached := false
	switch huType {
	case mjproto.HuType_H_GangShangPao: //杠上炮
		attached = true
		score = p.addDaHuScore(score, zimo)
		cardStr = append(cardStr, "杠上炮")

	case mjproto.HuType_H_GangShangHua: //杠上花
		attached = true
		score = p.addDaHuScore(score, zimo)
		cardStr = append(cardStr, "杠上花")

	case mjproto.HuType_H_HaiDiLao: //海底捞
		attached = true
		score = p.addDaHuScore(score, zimo)
		cardStr = append(cardStr, "海底捞月")

	case mjproto.HuType_H_HaiDiPao: //海底炮
		attached = true
		score = p.addDaHuScore(score, zimo)
		cardStr = append(cardStr, "海底炮")

	case mjproto.HuType_H_TianHu: //天胡
		if is258Jiang {
			//天胡需要258的将
			attached = true
			score = p.addDaHuScore(score, zimo)
			cardStr = append(cardStr, "天胡")
		}
	case mjproto.HuType_H_DiHu: //地胡
		if is258Jiang {
			//地胡需要258的将
			attached = true
			score = p.addDaHuScore(score, zimo)
			cardStr = append(cardStr, "地胡")
		}
	default:
	}

	//如果是小胡 平且有特殊胡牌累加 不累计小胡的得分
	if attached && pingHu {
		score = p.subXiaoHuScore(score, zimo)
	}

	return hu, fan, score, cardStr, pt, is258Jiang
}

//得到叫牌的信息 只要平胡 by 开总
func (p *HuParserChangSha) GetJiaoPais(pais []*majiang.MJPai) []*majiang.MJPai {
	log.T("长沙麻将解析器开始获取叫牌:pais:%v", majiangv2.ServerPais2string(pais))
	var jiaoPais []*majiang.MJPai

	tryHuPais := make([]*majiang.MJPai, len(pais)+1)
	copy(tryHuPais, pais)

	for i := 0; i < len(majiangv2.MjpaiMap); {
		tmpPai := majiangv2.InitMjPaiByIndex(i)

		//组装tryHu的牌数组
		tryHuPais[len(pais)] = tmpPai

		//转换成counts
		tryHuCounts := majiangv2.GettPaiStats(tryHuPais)

		//tryHu 平胡需要258做将
		canHu, _, jiang := p.TryHU(tryHuCounts, len(tryHuPais));
		//log.T("获取叫牌 tryHu canHu[%v] jiang[%v] is258Jiang[%v]", canHu, jiang, p.is258(jiang))
		if canHu && p.Is258(jiang) {
			jiaoPais = append(jiaoPais, tmpPai)
		}
		i += 4
	}
	log.T("长沙麻将解析器完成获取叫牌, 叫牌为 %v", jiaoPais)
	return jiaoPais
}

//得到分数
func (p *HuParserChangSha) HuScore() (fan int32, score int64) {
	return 0, 0
}

//累加大胡的得分
func (p *HuParserChangSha) addDaHuScore(score int64, zimo bool) int64 {
	if zimo {
		score += CHANGSHA_SCORE_DAHU_ZIMO
	} else {
		score += CHANGSHA_SCORE_DAHU
	}
	return score
}

//累加小胡的得分
func (p *HuParserChangSha) addXiaoHuScore(score int64, zimo bool) int64 {
	if zimo {
		score += CHANGSHA_SCORE_XIAOHU_ZIMO
	} else {
		score += CHANGSHA_SCORE_XIAOHU
	}
	return score
}

func (p *HuParserChangSha) subXiaoHuScore(score int64, zimo bool) int64 {
	if zimo {
		score -= CHANGSHA_SCORE_XIAOHU_ZIMO
	} else {
		score -= CHANGSHA_SCORE_XIAOHU
	}
	if score < 0 {
		score = 0
	}
	return score
}

func (p *HuParserChangSha) subDaHuScore(score int64, zimo bool) int64 {
	if zimo {
		score -= CHANGSHA_SCORE_DAHU_ZIMO
	} else {
		score -= CHANGSHA_SCORE_DAHU
	}
	if score < 0 {
		score = 0
	}
	return score
}

//起手胡的四种牌型 大四喜 板板胡 缺一色 六六顺
func (p *HuParserChangSha) ChangShaCanQiShouHu(handPai *majiang.MJHandPai, hupai *majiang.MJPai) (hu bool, fan int32, score int64, cardStr []string, pt []mjproto.PaiType) {
	log.T("长沙麻将解析器 开始判断是否是起手胡")
	if handPai.GangPais != nil && len(handPai.GangPais) > 0 {
		log.T("杠牌不为空 起手胡错误")
		return false, fan, score, cardStr, pt
	}

	if handPai.PengPais != nil && len(handPai.PengPais) > 0 {
		log.T("碰牌不为空 起手胡错误")
		return false, fan, score, cardStr, pt
	}

	//小胡
	if p.ChangShaIsSiXi(handPai, hupai) {
		//四喜
		hu = true
		score = p.addXiaoHuScore(score, false)
		cardStr = append(cardStr, "四喜")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_DAXISI)
	}
	if p.ChangShaIsBanBanHu(handPai, hupai) {
		//板板胡
		hu = true
		score = p.addXiaoHuScore(score, false)
		cardStr = append(cardStr, "板板胡")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_BANBANHU)
	}
	if p.ChangShaIsQueYiMen(handPai, hupai) {
		//缺一门
		hu = true
		score = p.addXiaoHuScore(score, false)
		cardStr = append(cardStr, "缺一色")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QUEYISE)
	}
	if p.ChangShaIsLiuLiuHu(handPai, hupai) {
		//六六顺
		hu = true
		score = p.addXiaoHuScore(score, false)
		cardStr = append(cardStr, "六六顺")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_LIULIUSHUN)
	}

	return hu, fan, score, cardStr, pt
}

//
func (p *HuParserChangSha) ChangShaGetHuScore(handPai *majiang.MJHandPai, hupai *majiang.MJPai, zimo bool) (hu bool, fan int32, score int64, cardStr []string, pt []mjproto.PaiType, is258Jiang bool, pingHu bool) {

	//log.T("开始判断胡牌类型 判定牌[%v]", hupai.GetDes())

	pais := []*majiang.MJPai{}
	if handPai.Pais != nil {
		pais = append(pais, handPai.Pais...)
	}
	if hupai != nil {
		pais = append(pais, hupai)
	}
	handCounts := majiangv2.GettPaiStats(pais)

	//判断各牌型是否满足
	canHu, _, jiang := p.TryHU(handCounts, len(pais))
	is258Jiang = p.Is258(jiang)
	//log.T("长沙平胡 tryHu[%v] ,jiang [%v],258jiang[%v]", canHu, jiang, is258Jiang)

	qingYiSe := p.ChangShaIsQingYiSe(handPai, hupai)

	ppHu := p.ChangShaIsPengPengHu(handPai, hupai)
	jjHu := p.ChangShaIsJiangJiangHu(handPai, hupai)

	qiDui := p.ChangShaIsQiDui(handPai, hupai)
	qiDuiHaoHua := p.ChangShaIsLongQiDui(handPai, hupai)
	qiDuiHaoHuaDouble := p.ChangShaIsDoubleLongQiDui(handPai, hupai)

	//统计分数/牌型
	if canHu && is258Jiang {
		pingHu = true
		//长沙平胡
		//log.T("平胡")
		hu = true
		score = p.addXiaoHuScore(score, zimo)
		cardStr = append(cardStr, "平胡")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_PINGHU)
	}

	if ppHu {
		//碰碰胡
		//log.T("碰碰胡")
		hu = true
		score = p.addDaHuScore(score, zimo)
		cardStr = append(cardStr, "碰碰胡")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_PENGPENGHU)
	}
	if jjHu {
		//将将胡
		//log.T("将将胡")
		hu = true
		score = p.addDaHuScore(score, zimo)
		cardStr = append(cardStr, "将将胡")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_JIANGJIANGHU)
	}

	if qiDui {
		//七对
		log.T("七对")
		hu = true
		score = p.addDaHuScore(score, zimo)
		cardStr = append(cardStr, "七小对")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QIXIAODUI)
	}

	if qiDuiHaoHua {
		//豪华七小对
		log.T("豪华七小对")
		hu = true
		score = p.addDaHuScore(score, zimo)
		score = p.addDaHuScore(score, zimo)

		cardStr = append(cardStr, "豪华七小对")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QIXIAODUI_HAOHUA)
	}

	if qiDuiHaoHuaDouble {
		//双豪华七小对
		log.T("双豪华七小对")
		hu = true
		score = p.addDaHuScore(score, zimo)
		score = p.addDaHuScore(score, zimo)
		score = p.addDaHuScore(score, zimo)

		cardStr = append(cardStr, "双豪华七小对")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QIXIAODUI_HAOHUA_DOUBLE)
	}

	if (pingHu || ppHu || jjHu || qiDui || qiDuiHaoHua || qiDuiHaoHuaDouble) && qingYiSe {
		//清一色 满足任意一个牌型即可
		//log.T("清一色")
		hu = true
		score = p.addDaHuScore(score, zimo)
		cardStr = append(cardStr, "清一色")
		pt = append(pt, mjproto.PaiType_H_CHANGSHA_QINGYISE)
	}

	return hu, fan, score, cardStr, pt, is258Jiang, pingHu
}

/******************** 长沙麻将 小胡 ********************/

//是否是大四喜 4张一样的牌
func (p *HuParserChangSha) ChangShaIsSiXi(handPai *majiang.MJHandPai, hupai *majiang.MJPai) bool {
	joinedPais := p.JoinHandPaiPaisAndHuPai(handPai, hupai)
	handCounts := majiangv2.GettPaiStats(joinedPais)
	for i := 0; i < len(handCounts); i++ {
		if handCounts[i] == 4 {
			return true
		}
	}
	return false
}

//是否是板板胡 没有一张 2、5、8
func (p *HuParserChangSha) ChangShaIsBanBanHu(handPai *majiang.MJHandPai, hupai *majiang.MJPai) bool {
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
func (p *HuParserChangSha) ChangShaIsQueYiMen(handPai *majiang.MJHandPai, hupai *majiang.MJPai) bool {
	joinedPais := p.JoinAllHandPaiAndHuPai(handPai, hupai)
	t, s, w := 0, 0, 0
	for i := 0; i < len(joinedPais); i++ {
		switch *joinedPais[i].Flower {
		case majiang.T:
			t++
		case majiang.S:
			s++
		case majiang.W:
			w++
		}
	}
	return t == 0 || s == 0 || w == 0
}

//是否是六六胡 至少两副刻子
func (p *HuParserChangSha) ChangShaIsLiuLiuHu(handPai *majiang.MJHandPai, hupai *majiang.MJPai) bool {
	joinedPais := p.JoinAllHandPaiAndHuPai(handPai, hupai)
	handCounts := majiangv2.GettPaiStats(joinedPais)
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
func (p *HuParserChangSha) ChangShaIsPengPengHu(handPai *majiang.MJHandPai, hupai *majiang.MJPai) bool {
	//碰碰胡不能有吃牌 吃牌不为空 吃牌数组大于0
	if chiPais := handPai.GetChiPais(); chiPais != nil || len(chiPais) > 0 {
		return false
	}

	joinedPais := p.JoinHandPaiPaisAndHuPai(handPai, hupai)
	handCounts := majiangv2.GettPaiStats(joinedPais)
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
func (p *HuParserChangSha) ChangShaIsJiangJiangHu(handPai *majiang.MJHandPai, hupai *majiang.MJPai) bool {
	//log.T("判断是否是将将胡")
	if handPai.GangPais != nil && len(handPai.GangPais) > 0 {
		//log.T("杠牌不为空 将将胡失败")
		return false
	}
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
func (p *HuParserChangSha) ChangShaIsQingYiSe(handPai *majiang.MJHandPai, hupai *majiang.MJPai) bool {
	pais := p.JoinAllHandPaiAndHuPai(handPai, hupai)
	return p.IsQingYiSe(pais)
}

//是否是七小对 任意花色组成的七对牌
func (p *HuParserChangSha) ChangShaIsQiDui(handPai *majiang.MJHandPai, hupai *majiang.MJPai) bool {
	//log.T("判断是否是七小对")
	joinedPais := p.JoinAllHandPaiAndHuPai(handPai, hupai)
	handCounts := majiangv2.GettPaiStats(joinedPais)
	duiCount := 0
	for i := 0; i < len(handCounts); i++ {
		if handCounts[i] == 0 {
			continue
		}
		if handCounts[i] != 2 {
			return false
		} else {
			duiCount++
		}
	}
	if duiCount != 7 {
		return false
	}
	return true
}

//是否是豪华七小对
func (p *HuParserChangSha) ChangShaIsLongQiDui(handPai *majiang.MJHandPai, hupai *majiang.MJPai) bool {
	longCount := 0 //勾数
	duiCount := 0  //对数
	//log.T("判断是否是豪华七小对")
	//不能有杠
	if handPai.GangPais != nil && len(handPai.GangPais) >= 0 {
		//log.T("杠牌不为空 豪华七小对失败")
		return false
	}

	joinedPais := p.JoinAllHandPaiAndHuPai(handPai, hupai)
	handCounts := majiangv2.GettPaiStats(joinedPais)

	for i := 0; i < len(handCounts); i++ {
		if handCounts[i] == 0 {
			continue
		}
		if handCounts[i] == 4 {
			//勾
			longCount++
		} else if (handCounts[i] == 2) {
			duiCount++
		} else {
			//牌数不符合
			return false
		}
	}
	if longCount == 1 && duiCount == 5 {
		return true
	}

	//log.T("四张数[%v]对数[%v] 豪华七小对失败", longCount, duiCount)

	return false
}

//是否是双豪华七小对
func (p *HuParserChangSha) ChangShaIsDoubleLongQiDui(handPai *majiang.MJHandPai, hupai *majiang.MJPai) bool {
	longCount := 0 //勾数
	duiCount := 0  //对数

	//不能有杠
	if handPai.GangPais != nil && len(handPai.GangPais) >= 0 {
		return false
	}

	joinedPais := p.JoinAllHandPaiAndHuPai(handPai, hupai)
	handCounts := majiangv2.GettPaiStats(joinedPais)

	for i := 0; i < len(handCounts); i++ {
		if handCounts[i] == 4 {
			//勾
			longCount++
		} else if (handCounts[i] == 2) {
			duiCount++
		} else {
			//牌数不符合
			return false
		}
	}
	if longCount == 2 && duiCount == 4 {
		return true
	}
	return false
}

////////////////////////////////////////////////////////////////////////

//将手牌碰牌刚牌与huPai拼接成数组
func (p *HuParserChangSha) JoinAllHandPaiAndHuPai(handPai *majiang.MJHandPai, hupai *majiang.MJPai) []*majiang.MJPai {
	pais := []*majiang.MJPai{}
	if hupai != nil {
		pais = append(pais, hupai)
	}

	if handPai.Pais != nil {
		pais = append(pais, handPai.Pais...)
	}

	if handPai.GangPais != nil {
		pais = append(pais, handPai.GangPais...)
	}

	if handPai.PengPais != nil {
		pais = append(pais, handPai.PengPais...)
	}
	return pais
}

//将手牌与huPai拼接成数组
func (p *HuParserChangSha) JoinHandPaiPaisAndHuPai(handPai *majiang.MJHandPai, hupai *majiang.MJPai) []*majiang.MJPai {
	pais := []*majiang.MJPai{}
	if hupai != nil {
		pais = append(pais, hupai)
	}

	if handPai.Pais != nil {
		pais = append(pais, handPai.Pais...)
	}

	return pais
}
