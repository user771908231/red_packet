package majiang

import (
	"testing"
	"casino_common/common/log"
)

var handpai *MJHandPai

func init() {
	log.InitLogger("", "") //初始化日志处理
	handpai = new(MJHandPai)
	//初始化
	index := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	for _, i := range index {
		handpai.Pais = append(handpai.Pais, InitMjPaiByIndex(i))
	}
	handpai.InPai = InitMjPaiByIndex(14)
}

func TestNewDefaultMjParser(t *testing.T) {
	t.Logf("需要解析的牌:%v", ServerPais2string(handpai.Pais))
	t.Logf("需要解析的inpai 牌:%v", handpai.InPai.LogDes())

	for _, pai := range handpai.Pais {
		t.Logf("打印牌[%v],des[%v]clientId[%v]", pai, pai.LogDes(), pai.GetClientId())
	}

	p := NewDefaultMjParser(handpai)
	t.Logf("解析出来的链子:%v", p.LianZis)
	for i, l := range p.LianZis {
		t.Logf("解析出来的链子%v:%v", i, ServerPais2string(l.pais))
	}

	t.Logf("解析出来的对子:%v", p.DuiZis)
	for i, d := range p.DuiZis {
		t.Logf("解析出来的对子%v:%v", i, ServerPais2string(d.pais))
	}

	t.Logf("解析出来的三条:%v", p.SanTiao)
	for i, s := range p.SanTiao {
		t.Logf("解析出来的对子%v:%v", i, ServerPais2string(s.pais))
	}

	t.Logf("解析出来的四条:%v", p.SiTiao)
	for i, si := range p.SiTiao {
		t.Logf("解析出来的对子%v:%v", i, ServerPais2string(si.pais))
	}
	t.Logf("解析出来的单张:%v", ServerPais2string(p.Dan))
}

func TestGetOutPai(t *testing.T) {
	t.Logf("解析之前的牌: %v", ServerPais2string(handpai.Pais))
	ret := GetOutPai(handpai)
	t.Logf("测试完成.ret :%v", ret.LogDes())
}

func TestGetDingQUe(t *testing.T) {

}
