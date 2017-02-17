package majiangv2

import (
	"sort"
	"fmt"
	"casino_common/common/log"
	"casino_majiang/service/majiang"
)

//牌型解析
type MjParser struct {
	count_t    int              //筒子的数量
	count_s    int              //条子的数量
	count_w    int              //万子的数量
	count_less int32            //最少的花色
	LianZis    []*MjPBLianzi    //链子数量
	DuiZis     []*MjPBDuizi     //对子数量
	SanTiao    []*MjPbSanTiao   //三条数量
	SiTiao     []*MjPbSiTiao    //四条数量
	ALL        []*majiang.MJPai //初始化的所有牌
	Dan        []*majiang.MJPai //初始化的单牌
	QUE        []*majiang.MJPai //缺的牌
}

//牌型里面的链子
/**
	如链子 3,4,5,6,7 万
	flower : W
	min  :3
	max  :7
 */
type MjPBLianzi struct {
	pais   []*majiang.MJPai //链子的牌
	flower int32            //花色
	min    int32            //最小的值
	max    int32            //最大的值
}

func (l *MjPBLianzi) GetCount() int {
	return len(l.pais)
}

func (l *MjPBLianzi) AddPai(pai *majiang.MJPai) {
	l.pais = append(l.pais, pai)
	//做初始化
	l.flower = pai.GetFlower()
	var paiList majiang.MjPaiList = l.pais
	sort.Sort(paiList)
	l.min = paiList[0].GetValue()
	l.max = paiList[l.GetCount()-1].GetValue()
}

//牌型里面的对子
type MjPBDuizi struct {
	pais []*majiang.MJPai
}

//牌型里面的三条
type MjPbSanTiao struct {
	pais []*majiang.MJPai
}

//牌型里面的四条
type MjPbSiTiao struct {
	pais []*majiang.MJPai
}

//解析可能会用到 GettPaiStats
/*
	解析的过程
	1,为了保证机器人能够更容易胡牌，所以需要保证链子先被解析
	2,链子解析之后 再解析 4，3，2
	3，最后剩下的就是单张

*/

func (p *MjParser) Parse(flower int32) {
	p.parseSort()
	p.ParseQue(flower)
	p.ParseCount()
	p.ParseLian(flower)
	p.Parse432()
	p.ParseDan()
}

//首先解析数量
func (p *MjParser) ParseCount() {
	//init
	p.count_t = 0
	p.count_s = 0
	p.count_w = 0

	//count
	for _, pai := range p.ALL {
		if pai != nil {
			if pai.GetFlower() == T {
				p.count_t ++
			} else if pai.GetFlower() == W {
				p.count_w ++
			} else if pai.GetFlower() == S {
				p.count_s ++
			}
		}
	}

	//去张数最少的牌
	if p.count_t < p.count_s && p.count_t < p.count_w {
		p.count_less = T
	} else {
		if p.count_s < p.count_w {
			p.count_less = S
		} else {
			p.count_less = W
		}
	}
}

//删除一个list
func (p *MjParser) rmPaiList(list []*majiang.MJPai) {
	for _, pai := range list {
		p.rmTempPai(pai)
	}
}

//删除temp中指定的牌
func (p *MjParser) rmTempPai(rmp *majiang.MJPai) {
	index := -1
	for i, pai := range p.Dan {
		if pai != nil && pai.GetIndex() == rmp.GetIndex() {
			index = i
			break
		}
	}
	if index > -1 {
		p.Dan = append(p.Dan[:index], p.Dan[index+1:]...)
	} else {
		log.E("服务器错误：删除手牌的时候出错，没有找到对应的手牌[%v]", rmp)
	}
}

//排序
func (p *MjParser) parseSort() {
	var list majiang.MjPaiList = p.Dan
	sort.Sort(list)
}

//初始化缺牌 todo 这里暂时可以不做,外部出牌的地方已经做了，之后可以移动过来
func (p *MjParser) ParseQue(flower int32) {
	for _, quepai := range p.Dan {
		if quepai.GetFlower() == flower {
			p.QUE = append(p.QUE, quepai)
		}
	}
}

func (p *MjParser) GetDanPaiByCountInex(index int32) *majiang.MJPai {
	for _, pai := range p.Dan {
		if pai.GetClientId() == (index + 1) {
			return pai
		}
	}
	return nil
}

//取链子的时候不用判断缺牌
func (p *MjParser) ParseLian(flower int32) {
	//log.T("开始解析链子")
	l := new(MjPBLianzi)

	// 1,2,3 wst
	//log.T("需要解析的链子1:%v", ServerPais2string(p.Dan))
	count := GettPaiStats(p.Dan)
	//log.T("需要解析的链子,解析之后的count:%v", count)

	//循环得到链子
	for f := int32(1); f < 4; f++ {

		//如果是缺的花色，不用循环
		if f == flower {
			continue
		}

		isBreak := false
		//循环不是缺的花色
		for i := (f - 1) * 9; i < 7*f; i++ {
			if count[i] != 0 && count[i+1] != 0 && count[i+2] != 0 {
				l.pais = make([]*majiang.MJPai, 3)
				l.pais[0] = p.GetDanPaiByCountInex(i)
				l.pais[1] = p.GetDanPaiByCountInex(i + 1)
				l.pais[2] = p.GetDanPaiByCountInex(i + 2)
				p.LianZis = append(p.LianZis, l)
				isBreak = true
				break
			}
		}
		//跳出2层循环
		if isBreak {
			break
		}
	}

	if l != nil && l.GetCount() == 3 {
		//但单牌list 中 删除得到的链子
		p.rmPaiList(l.pais)
		//删除之后继续检测链子
		p.ParseLian(flower)
	}

}

//解析对子
func (p *MjParser) ParseDuzi(index int) {
	duizi := new(MjPBDuizi)
	duizi.pais = make([]*majiang.MJPai, 0)
	for _, pai := range p.Dan {
		if pai.GetClientId() == int32(index+1) {
			if pai == nil {
				log.E("没有找到 对应index:%v的牌:", index)
			}
			duizi.pais = append(duizi.pais, pai)
		}
	}
	p.DuiZis = append(p.DuiZis, duizi)
	p.rmPaiList(duizi.pais)
}

//解析三条
func (p *MjParser) ParseSanTiao(index int) {
	santiao := new(MjPbSanTiao)
	santiao.pais = make([]*majiang.MJPai, 0)
	for _, pai := range p.Dan {
		if pai.GetClientId() == int32(index+1) {
			santiao.pais = append(santiao.pais, pai)
		}
	}
	p.SanTiao = append(p.SanTiao, santiao)
	p.rmPaiList(santiao.pais)

}

//解析四条
func (p *MjParser) ParseSiTiao(index int) {
	siTiao := new(MjPbSiTiao)
	siTiao.pais = make([]*majiang.MJPai, 0)
	for _, pai := range p.Dan {
		if pai.GetClientId() == int32(index+1) {
			siTiao.pais = append(siTiao.pais, pai)
		}
	}
	p.SiTiao = append(p.SiTiao, siTiao)
	p.rmPaiList(siTiao.pais)

}

//解析对子链子等
func (p *MjParser) Parse432() {
	//log.T("需要计算的牌的信息 : %v", ServerPais2string(p.Dan))
	count := GettPaiStats(p.Dan)
	//log.T("计算出来的牌 couont : %v", count)
	for i, c := range count {
		if c == 2 {
			//组装对子
			p.ParseDuzi(i)
		} else if c == 3 {
			//组装砍
			p.ParseSanTiao(i)
		} else if c == 4 {
			//组装四张
			p.ParseSiTiao(i)
		}
	}
}

//解析单张
func (p *MjParser) ParseDan() {
	rmList := make([]*majiang.MJPai, 0)
	//查看单张是否可以解析到链子中去
	for _, pai := range p.Dan {
		//查看单排是否可以组装到链子中去
		if p.Dan2Lian(pai) {
			rmList = append(rmList, pai)
		}
	}
	//剩下的牌都是单张
	p.rmPaiList(rmList)
}

//判断单排是否可以加入到链子中
func (p *MjParser) Dan2Lian(pai *majiang.MJPai) bool {
	for _, lian := range p.LianZis {
		if lian.flower == pai.GetFlower() {
			if (lian.max == pai.GetValue()-1 ) || lian.min == pai.GetValue()+1 {
				lian.AddPai(pai)
				return true
			}
		}
	}
	return false
}

//得到推荐打的牌
func (p *MjParser) GetOutPai() *majiang.MJPai {
	//1,找缺牌
	if len(p.QUE) > 0 {
		return p.QUE[0]
	}

	//2找单排
	if len(p.Dan) > 0 {
		return p.Dan[0]
	}

	//3找链子的多张4,7
	for _, lian := range p.LianZis {
		if lian.GetCount()%3 == 1 {
			return lian.pais[0]
		}
	}

	//4，找链子的多张 ！= 5
	for _, lian := range p.LianZis {
		if lian.GetCount()%3 != 0 {
			return lian.pais[0]
		}
	}

	//4，以上情况都找不到的话，返回所有牌中的第一张
	return p.ALL[0]

}

//初始化
func NewDefaultMjParser(handpai *majiang.MJHandPai) *MjParser {
	//log.T("开始初始化解析器.")
	//初始化所有牌
	p := new(MjParser)
	p.LianZis = make([]*MjPBLianzi, 0)
	p.DuiZis = make([]*MjPBDuizi, 0)
	p.SanTiao = make([]*MjPbSanTiao, 0)
	p.SiTiao = make([]*MjPbSiTiao, 0)
	p.QUE = make([]*majiang.MJPai, 0)

	p.ALL = handpai.Pais
	if handpai.InPai != nil {
		p.ALL = append(p.ALL, handpai.InPai)
	}
	p.Dan = make([]*majiang.MJPai, len(p.ALL))
	copy(p.Dan, p.ALL)

	//开始解析
	p.Parse(handpai.GetQueFlower())

	//log.T("解析出来的链子:%v", p.LianZis)
	//log.T("解析出来的对子:%v", p.DuiZis)
	//log.T("解析出来的三条:%v", p.SanTiao)
	//log.T("解析出来的四条:%v", p.SiTiao)
	//log.T("解析出来的单张:%v", p.Dan)
	return p
}

func GetOutPai(handpai *majiang.MJHandPai) (ret *majiang.MJPai) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			log.E("解析麻将的时候出现了错误")
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			ret = handpai.Pais[0]
		}
	}()

	ret = NewDefaultMjParser(handpai).GetOutPai()
	if ret == nil {
		ret = handpai.Pais[0]
	}
	return ret
}

func GetDingQUe(handpai *majiang.MJHandPai) (ret int32) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			log.E("解析麻将的时候出现了错误")
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			ret = W
		}
	}()
	//得到返回值
	ret = int32(NewDefaultMjParser(handpai).count_less)
	return ret
}
