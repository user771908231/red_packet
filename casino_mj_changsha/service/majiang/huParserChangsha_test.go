package majiang

import (
	"testing"
	"casino_mj_changsha/msg/protogo"
	"casino_common/common/log"
	"strings"
	"time"
)

func TestHuParserChangSha_GetCanHu(t *testing.T) {
	var handPai = new(MJHandPai)
	var huPai = new(MJPai)
	p := new(HuParserChangSha)

	//log.T(" ")
	//
	//log.T("开始测试起手胡: 大四喜")
	//handPai, huPai = getChangShaDaSiXi()
	//printCanHuRet(p.GetCanHu(handPai, huPai, true, mjproto.HuType_H_changsha_qishouhu))
	//
	//log.T(" ")
	//
	//log.T("开始测试起手胡: 板板胡")
	//handPai, huPai = getChangShaBanBanHu()
	//printCanHuRet(p.GetCanHu(handPai, huPai, true, mjproto.HuType_H_changsha_qishouhu))
	//
	//log.T(" ")
	//
	//log.T("开始测试起手胡: 缺一色")
	//handPai, huPai = getChangShaQueYiSe()
	//printCanHuRet(p.GetCanHu(handPai, huPai, true, mjproto.HuType_H_changsha_qishouhu))

	//log.T(" ")
	//
	//log.T("开始测试起手胡: 六六顺")
	//handPai, huPai = getChangSha66Shun()
	//printCanHuRet(p.GetCanHu(handPai, huPai, true, mjproto.HuType_H_changsha_qishouhu))

	//log.T(" ")
	//
	//log.T("开始测试起手胡: 碰碰胡")
	//handPai, huPai = getChangShaPengPengHu()
	//printCanHuRet(p.GetCanHu(handPai, huPai, true, mjproto.HuType_H_GangShangPao))

	//log.T(" ")
	//
	//log.T("开始测试起手胡: 将将胡")
	//handPai, huPai = getChangShaJJHu()
	//printCanHuRet(p.GetCanHu(handPai, huPai, true, mjproto.HuType_H_GangShangPao))

	//log.T(" ")
	//
	//log.T("开始测试起手胡: 清一色")
	//handPai, huPai = getChangShaQingyise()
	//printCanHuRet(p.GetCanHu(handPai, huPai, true, mjproto.HuType_H_GangShangPao))

	//log.T(" ")
	//
	//log.T("开始测试起手胡: 七对")
	//handPai, huPai = getChangShaQidui()
	//printCanHuRet(p.GetCanHu(handPai, huPai, true, mjproto.HuType_H_GangShangPao))
	log.T(" ")
	//
	//log.T("开始测试起手胡: 豪华七对")
	//handPai, huPai = getJiangqidui()
	//printCanHuRet(p.GetCanHu(handPai, huPai, true, mjproto.HuType_H_GangShangPao))

	log.T("开始测试起手胡: bug")
	handPai, huPai = getTest01()
	//printCanHuRet(p.GetCanHu(handPai, huPai, true, mjproto.HuType_H_GangShangPao))

	//p.GetJiaoPais(handPai.GetPais())
	time.Sleep(time.Duration(3) * time.Second)
}

func printCanHuRet(hu bool, fan int32, score int64, cardStr []string, pt []mjproto.PaiType, is258Jiang bool) {
	log.T("开始打印胡牌描述信息")
	log.T("是否能胡[%v]", hu)
	log.T("得分[%v]", score)
	log.T("胡牌描述[%v]", strings.Join(cardStr, " "))
	log.T("258将对[%v]", is258Jiang)
	if pt == nil || len(pt) <= 0 {
		log.T("打印胡牌描述信息结束")
		return
	}
	for i, _ := range pt {
		log.T("胡牌类型[%v]", pt[i])
	}
	log.T("打印胡牌描述信息结束")
}

func getChangShaNoHu() (*MJHandPai, *MJPai) {
	huPaiDes := "T_4"
	inPaiDes := "T_4"
	paisDes := []string{"S_6", "S_7", "S_8", "T_7", "T_7", "T_9", "T_4", "T_4", "T_4"}
	pengPaisDes := []string{"T_1", "T_1", "T_1"}
	gangPaisDes := []string{"T_3", "T_3", "T_3", "T_3"}
	handPai := getMjHandPai(inPaiDes, pengPaisDes, gangPaisDes, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

//平胡 1番
func getChangShaPinghu() (*MJHandPai, *MJPai) {
	huPaiDes := "W_9"
	inPaiDes := "W_9"
	//paisDes		:= []string{"S_6", "S_7", "S_8", "T_7", "T_8", "T_9", "T_4", "T_4", "T_4", "T_2"}
	//paisDes := []string{"S_6", "S_7", "S_8", "T_7", "T_7", "T_9", "W_4", "W_4", "W_4", "T_2"}
	paisDes := []string{"W_7", "W_7", "W_8", "W_8", "W_9", "S_1", "S_2", "S_3", "S_1", "S_2", "S_3", "T_2", "T_2"}

	handPai := getMjHandPai(inPaiDes, nil, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

func getChangShaDaSiXi() (*MJHandPai, *MJPai) {
	huPaiDes := "W_9"
	inPaiDes := "W_9"
	//paisDes		:= []string{"S_6", "S_7", "S_8", "T_7", "T_8", "T_9", "T_4", "T_4", "T_4", "T_2"}
	//paisDes := []string{"S_6", "S_7", "S_8", "T_7", "T_7", "T_9", "W_4", "W_4", "W_4", "T_2"}
	paisDes := []string{"W_7", "W_7", "W_8", "W_8", "W_9", "S_1", "S_1", "S_1", "S_1", "T_2", "T_2"}

	handPai := getMjHandPai(inPaiDes, nil, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

func getChangShaBanBanHu() (*MJHandPai, *MJPai) {
	huPaiDes := "W_1"
	inPaiDes := "W_1"
	//paisDes		:= []string{"S_6", "S_7", "S_8", "T_7", "T_8", "T_9", "T_4", "T_4", "T_4", "T_2"}
	//paisDes := []string{"S_6", "S_7", "S_8", "T_7", "T_7", "T_9", "W_4", "W_4", "W_4", "T_2"}
	paisDes := []string{"W_1", "W_1", "W_3", "W_3", "W_4", "W_4", "S_1", "S_1", "S_1", "S_1", "T_3", "T_3"}

	handPai := getMjHandPai(inPaiDes, nil, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

func getChangShaQueYiSe() (*MJHandPai, *MJPai) {
	huPaiDes := "W_6"
	inPaiDes := "W_6"
	//paisDes		:= []string{"S_6", "S_7", "S_8", "T_7", "T_8", "T_9", "T_4", "T_4", "T_4", "T_2"}
	//paisDes := []string{"S_6", "S_7", "S_8", "T_7", "T_7", "T_9", "W_4", "W_4", "W_4", "T_2"}
	paisDes := []string{"W_7", "W_7", "W_8", "W_8", "W_6", "T_1", "T_1", "T_1", "T_1", "T_3", "T_3"}

	handPai := getMjHandPai(inPaiDes, nil, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

func getChangSha66Shun() (*MJHandPai, *MJPai) {
	huPaiDes := "W_6"
	inPaiDes := "W_6"
	//paisDes		:= []string{"S_6", "S_7", "S_8", "T_7", "T_8", "T_9", "T_4", "T_4", "T_4", "T_2"}
	//paisDes := []string{"S_6", "S_7", "S_8", "T_7", "T_7", "T_9", "W_4", "W_4", "W_4", "T_2"}
	paisDes := []string{"W_7", "W_7", "W_7", "W_6", "W_6", "T_1", "T_1", "T_1", "T_1", "T_3", "T_3"}

	handPai := getMjHandPai(inPaiDes, nil, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

func getChangShaPengPengHu() (*MJHandPai, *MJPai) {
	huPaiDes := "W_8"
	inPaiDes := "W_8"
	//paisDes		:= []string{"S_6", "S_7", "S_8", "T_7", "T_8", "T_9", "T_4", "T_4", "T_4", "T_2"}
	//paisDes := []string{"S_6", "S_7", "S_8", "T_7", "T_7", "T_9", "W_4", "W_4", "W_4", "T_2"}
	paisDes := []string{"W_7", "W_7", "W_7", "W_8", "W_8", "S_1", "S_1", "S_1", "T_2", "T_2", "T_2", "T_7", "T_7"}

	pengPaisDes := []string{"T_1", "T_1", "T_1"}
	gangPaisDes := []string{"T_3", "T_3", "T_3", "T_3"}

	handPai := getMjHandPai(inPaiDes, pengPaisDes, gangPaisDes, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)

}

func getChangShaJJHu() (*MJHandPai, *MJPai) {
	huPaiDes := "W_8"
	inPaiDes := "W_8"
	//paisDes		:= []string{"S_6", "S_7", "S_8", "T_7", "T_8", "T_9", "T_4", "T_4", "T_4", "T_2"}
	//paisDes := []string{"S_6", "S_7", "S_8", "T_7", "T_7", "T_9", "W_4", "W_4", "W_4", "T_2"}
	paisDes := []string{"W_2", "W_2", "W_2", "W_8", "W_8", "S_2", "S_2", "S_2", "T_2", "T_2", "T_2", "T_5", "T_5"}

	pengPaisDes := []string{"T_8", "T_8", "T_8"}
	//gangPaisDes := []string{"T_3", "T_3", "T_3", "T_3"}

	handPai := getMjHandPai(inPaiDes, pengPaisDes, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)

}

//对对胡 2番
func getChangShaDuiduihu() (*MJHandPai, *MJPai) {
	inPaiDes := "T_4" //4T
	huPaiDes := "T_4"
	paisDes := []string{"S_6", "S_6", "S_6", "T_7", "T_7", "T_7", "T_4"} //666S 777T 4T
	pengPaisDes := []string{"T_2", "T_2", "T_2"}                         //111T 222T
	gangPaisDes := []string{"T_1", "T_1", "T_1", "T_1"}
	handPai := getMjHandPai(inPaiDes, pengPaisDes, gangPaisDes, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

//清一色 3番
func getChangShaQingyise() (*MJHandPai, *MJPai) {
	inPaiDes := "T_9" //4T
	huPaiDes := "T_9"
	paisDes := []string{"T_2", "T_3", "T_4", "T_5", "T_6", "T_7", "T_8", "T_8", "T_9", "T_9"} //123T 456T 4T
	gangPaisDes := []string{"T_1", "T_1", "T_1", "T_1"}
	handPai := getMjHandPai(inPaiDes, nil, gangPaisDes, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

//七对 3番
func getChangShaQidui() (*MJHandPai, *MJPai) {
	inPaiDes := "S_6"
	huPaiDes := "S_6"                                                                                              //6S
	paisDes := []string{"S_1", "S_1", "S_2", "S_2", "S_4", "S_4", "T_9", "T_9", "T_7", "T_7", "S_7", "S_7", "S_6"} //11S 22S 44S 99T 77T 6S
	handPai := getMjHandPai(inPaiDes, nil, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

//将七对 3番
func getJiangqidui() (*MJHandPai, *MJPai) {
	inPaiDes := "S_6"                                                                                              //6S
	huPaiDes := "S_6"                                                                                              //6S
	paisDes := []string{"T_2", "T_2", "S_2", "S_2", "S_5", "S_5", "T_8", "T_8", "T_5", "T_5", "S_6", "S_6", "S_6"} //11S 22S 44S 99T 77T 6S
	handPai := getMjHandPai(inPaiDes, nil, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

//清对 4番
func getChangShaQingdui() (*MJHandPai, *MJPai) {
	inPaiDes := "T_9"                                                    //9T
	huPaiDes := "T_9"                                                    //
	paisDes := []string{"T_5", "T_5", "T_5", "T_7", "T_7", "T_7", "T_9"} //555T 777T 9T
	pengPaisDes := []string{"T_1", "T_1", "T_1", "T_3", "T_3", "T_3"}    //111T 333T
	handPai := getMjHandPai(inPaiDes, pengPaisDes, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

//将对 4番
func getChangShaJiangdui() (*MJHandPai, *MJPai) {
	inPaiDes := "S_2" //2S
	huPaiDes := "S_2"
	paisDes := []string{"S_5", "S_5", "S_5", "S_8", "S_8", "S_8", "S_2"} //555S 888S 2S
	pengPaisDes := []string{"T_2", "T_2", "T_2", "T_5", "T_5", "T_5"}    //222T 555T
	gangPaisDes := []string{"T_8", "T_8", "T_8", "T_8"}
	handPai := getMjHandPai(inPaiDes, pengPaisDes, gangPaisDes, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

//龙七对 5番
func getChangShaLongqidui() (*MJHandPai, *MJPai) {
	inPaiDes := "T_7" //7T
	huPaiDes := "T_7"
	paisDes := []string{"S_1", "S_1", "S_2", "S_2", "S_4", "S_4", "T_9", "T_9", "T_6", "T_6", "T_7", "T_7", "T_7"} //11S 22S 44S 99T 777T
	handPai := getMjHandPai(inPaiDes, nil, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

//清七对 5番
func getChangShaQingqidui() (*MJHandPai, *MJPai) {
	inPaiDes := "S_7"
	huPaiDes := "S_7"                                                                                              //7S
	paisDes := []string{"S_1", "S_1", "S_2", "S_2", "S_3", "S_3", "S_5", "S_5", "S_6", "S_6", "S_9", "S_9", "S_7"} //11S 22S 55S 66S 7S
	handPai := getMjHandPai(inPaiDes, nil, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

//清幺九 5番
func getChangShaQingyaojiu() (*MJHandPai, *MJPai) {
	inPaiDes := "S_4" //4S
	huPaiDes := "S_4"
	paisDes := []string{"S_1", "S_2", "S_3", "S_1", "S_2", "S_3", "S_1", "S_2", "S_3", "S_4"} //123S 123S 123S 4S
	pengPaisDes := []string{"S_9", "S_9", "S_9"}                                              //999S
	handPai := getMjHandPai(inPaiDes, pengPaisDes, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

//清龙七对 6番
func getChangShaQinglongqidui() (*MJHandPai, *MJPai) {
	inPaiDes := "T_7" //7T
	huPaiDes := "T_7"
	paisDes := []string{"T_1", "T_1", "T_2", "T_2", "T_3", "T_3", "T_4", "T_4", "T_5", "T_5", "T_7", "T_7", "T_7"} //11T 22T 44T 777T
	handPai := getMjHandPai(inPaiDes, nil, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}

func getTest01() (*MJHandPai, *MJPai) {
	inPaiDes := "S_5" //7T
	huPaiDes := "S_5"
	paisDes := []string{"T_2", "T_2", "T_2", "T_2", "T_3", "T_3", "T_3", "T_3", "T_4", "S_5"} //11T 22T 44T 777T
	handPai := getMjHandPai(inPaiDes, nil, nil, paisDes)
	return handPai, InitMjPaiByDes(huPaiDes, handPai)
}
