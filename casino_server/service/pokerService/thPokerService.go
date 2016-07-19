package pokerService

import (
	"casino_server/utils"
	"casino_server/msg/bbprotogo"
	"casino_server/common/log"
	"sort"
)

//德州扑克的纸牌

/**
	德州扑克判断牌的大小
	1,手牌加上公共牌中选5张需要对比的牌
	2,选出来的牌再和其他人的牌来比大小
 */

//比较两幅德州牌的大小Big
var THPOKER_COMPARE_BIG = 1
var THPOKER_COMPARE_SMALL = 2
var THPOKER_COMPARE_EQUALS = 3


//const
var (
	THPOKER_TYPE_GAOPAI 				int32	=	1;
	THPOKER_TYPE_YIDUI 				int32	=	2;
	THPOKER_TYPE_LIANGDUI 				int32	=	3;
	THPOKER_TYPE_SANTIAO 				int32	=	4;
	THPOKER_TYPE_SHUNZI 				int32	=	5;
	THPOKER_TYPE_TONGHUA 				int32	=	6;
	THPOKER_TYPE_HULU 				int32	=	7;
	THPOKER_TYPE_SITIAO 				int32	=	8;
	THPOKER_TYPE_TONGHUASHUN 			int32	=	9;
	THPOKER_TYPE_HUANGJIATONGHUASHUN 		int32	=	10;
)


//config
type  ThCardsList 	[]*ThCards	//需要比较的牌
type  CardsList		[]*bbproto.Pai	//对牌进行排序

type ThCards struct {
	ThType		int32
	KeyValue	[]int32
	Cards	[]*bbproto.Pai
	DuiziCount	int32
	SanTiaoCount	int32
	SiTiaoCount     int32		//四条
	CardsStatistics  []int32
	IsShunzi	bool		//是否是顺子
	IsTongHua	bool		//是否是同花
	IsSiTiao	bool		//是否是四条
	IsSanTiao	bool		//是否是三条
	IsHulu		bool		//是否是葫芦
	IsGaoPai	bool		//是否是高牌
	IsDuiZi		bool		//是否是对子
	IsLiangDui	bool		//是否是两队
	IsWin		bool
}

//创建一个德州的牌面,并且对属性做0值初始化
func NewThCards() *ThCards{

	result := &ThCards{}
	result.ThType = 0
	result.KeyValue = make([]int32,5)
	result.DuiziCount = 0
	result.SanTiaoCount = 0
	result.SiTiaoCount = 0
	result.CardsStatistics = make([]int32,14)
	result.IsShunzi = false
	result.IsTongHua = false
	result.IsSiTiao = false
	result.IsSanTiao = false
	result.IsHulu = false
	result.IsGaoPai = false
	result.IsDuiZi = false
	result.IsLiangDui = false
	result.IsWin = false
	return result
}

func (t *ThCards) LogString(){
	//log.T("牌类型[%v]", *t.ThType)
	//log.T("对子数量[%v]",*t.DuiziCount)
	//log.T("三条数量[%v]",*t.SanTiaoCount)
	//log.T("四条数量[%v]",*t.SiTiaoCount)
	//log.T("牌面[%v]",*t.Cards)
}


//检查是否是同花

//由于同花的时候就不可能出现对子,随意比较值的初始化可以在这一步进行
func ( t *ThCards) OnInitTongHuaStatus(){
	cdlist := t.Cards
	var result bool = true
	for i := 0; i < len(cdlist)-1; i++ {
		if *(cdlist[i].Flower) != *(cdlist[i+1].Flower){
			result = false
			break
		}
	}
	t.IsTongHua = result


}

//检查是否是顺子
func ( t *ThCards) OnInitShunZiStatus(){
	cdlist := t.Cards
	var result bool = true
	for i := 0; i < len(cdlist)-1; i++ {
		if (*(cdlist[i].Value) - *(cdlist[i+1]).Value) != 1{
			result = false
			break
		}
	}
	t.IsShunzi = result
}


//检测是否是四条
func ( t *ThCards) OnInitSiTiaoStatus() {
	if t.SiTiaoCount == 1 {
		(t.IsSiTiao) = true
		//四条的规则,先比较四条再比较单张,
		s := t.CardsStatistics
		for i := 0; i < len(s); i++ {
			if s[i] == 4 {
				t.KeyValue[0] = int32(i)
			}else if s[i] == 1{
				t.KeyValue[1] = int32(i)
			}
		}
	}
}

//检测是否是三条
func ( t *ThCards) OnInitSantiaoStatus() {
	if t.SanTiaoCount == 1 && t.DuiziCount == 0 {
		t.IsSanTiao = true

		//四条的规则,先比较四条再比较单张,
		s := t.CardsStatistics
		for i := 0; i < len(s); i++ {
			if s[i] == 3 {
				t.KeyValue[0] = int32(i)
			}
		}

		for i := 0; i < len(s); i++ {
			if s[i] == 1 {
				t.KeyValue[1] = int32(i)
			}
		}

		for i := 0; i < len(s); i++ {
			if s[i] == 1 &&  t.KeyValue[1] != int32(i) {
				t.KeyValue[2] = int32(i)
			}
		}
	}

}


//检测是否是两队
func ( t *ThCards) OnInitLiangDuiStatus(){
	if t.DuiziCount == 2 {
		t.IsLiangDui = true

		s := t.CardsStatistics
		//初始化比较值
		for i := 0; i < len(s); i++ {
			if s[i] == 2 {
				t.KeyValue[0] = int32(i)
			}
		}

		for i := 0; i < len(s); i++ {
			if s[i] == 1 &&  t.KeyValue[0] != int32(i) {
				t.KeyValue[1] = int32(i)
			}
		}

		for i := 0; i < len(s); i++ {
			if s[i] == 1 {
				t.KeyValue[2] = int32(i)
			}
		}

		
	}
}


//初始化一对
func ( t *ThCards) OnInitYiDuiStatus() {
	if t.DuiziCount == 1 && t.SanTiaoCount == 0 {
		//只有一个对子,且三条的个数为0
		t.IsSanTiao = true
		//初始化比较值
		s := t.CardsStatistics
		for i := 0; i < len(s); i++ {
			if s[i] == 3 {
				t.KeyValue[0] = int32(i)
			}
		}

		s2 := s
		for i := 1; i < 5; i++ {
			for j := 0; j < len(s2); j++ {
				if s2[j] == 1 {
					t.KeyValue[i] = int32(j)
					s2[j] = 0
				}
			}
		}

	}

}


//初始化葫芦
func (t *ThCards) OnInitHuLuStatus(){
	if t.SanTiaoCount == 1 && t.DuiziCount == 1 {
		t.ThType = THPOKER_TYPE_HULU
		//初始化比较值
		s := t.CardsStatistics
		for i := 0; i < len(s); i++ {
			if s[i] == 3 {
				t.KeyValue[0] = int32(i)
			}else if s[i] == 2{
				t.KeyValue[1] = int32(i)
			}
		}
	}
}


//分析牌面
func (c *ThCards) OnInitStatisticsCard() error{
	list := c.Cards

	//可以通过这个统计来计算,对子有多少,三条有多少,四条有多少
	c.CardsStatistics  = make([]int32,15)
	for i := 0; i < len(list); i++ {
		c.CardsStatistics[*list[i].Value] ++
	}

	//fmt.Println("统计出来的牌数量:s",c.CardsStatistics)
	for i := 0; i < len(c.CardsStatistics); i++ {
		//fmt.Println("开始检测:",c.CardsStatistics[i])

		if c.CardsStatistics[i] == 2  {
			//fmt.Println("检测到的对子是:",c.CardsStatistics[i])
			c.DuiziCount ++
		}else if c.CardsStatistics[i] == 3 {
			//fmt.Println("检测到的三条是:",c.CardsStatistics[i])
			c.SanTiaoCount ++
		}else if c.CardsStatistics[i] == 4 {
			//fmt.Println("检测到的四条是:",c.CardsStatistics[i])
			c.SiTiaoCount ++
		}
	}

	//如果是同花,排序之后的牌就是比较值的值
	for i := 0; i < len(c.Cards); i++ {
		c.KeyValue[i]= *(c.Cards[i].Value)
	}


	//log.T("统计之后的牌面值")

	c.OnInitTongHuaStatus()
	c.OnInitShunZiStatus()
	c.OnInitSiTiaoStatus()
	c.OnInitHuLuStatus()
	c.OnInitSantiaoStatus()
	c.OnInitLiangDuiStatus()
	c.OnInitYiDuiStatus()

	return nil
}

/**
	有了5张牌之后初始化牌
 */
func (c *ThCards) OnInit() error{
	//首先对牌进行排序
	var cdList CardsList = c.Cards
	sort.Sort(cdList)
	c.Cards = cdList

	//做统计初始化
	c.OnInitStatisticsCard()
	//fmt.Println("c.DuiziCount",*c.DuiziCount)
	//fmt.Println("c.SanTiaoCount",*c.SanTiaoCount)
	//fmt.Println("c.SiTiaoCount",*c.SiTiaoCount)
	//
	//fmt.Println("c.SiTiaoCount",*c.IsDuiZi)
	//fmt.Println("c.IsLiangDui",*c.IsLiangDui)
	//fmt.Println("c.IsTongHua",*c.IsTongHua)
	//fmt.Println("c.IsSanTiao",*c.IsSanTiao)
	//fmt.Println("c.IsSiTiao",*c.IsSiTiao)
	//fmt.Println("c.IsGaoPai",*c.IsGaoPai)
	//fmt.Println("c.IsHulu",*c.IsHulu)
	//fmt.Println("c.IsShunzi",*c.IsShunzi)


	//解析牌的keyValue值,属性
	if c.IsTongHua {
		if c.IsShunzi {
			if *cdList[0].Value == 14 {
				//如果值是A(14),表示这个牌是皇家同花顺
				c.ThType = THPOKER_TYPE_HUANGJIATONGHUASHUN
			}else{
				//同花顺
				c.ThType = THPOKER_TYPE_TONGHUASHUN
			}
		}else{
			//同花
			c.ThType = THPOKER_TYPE_TONGHUA
		}
	}else{
		if c.IsShunzi {
			c.ThType = THPOKER_TYPE_SHUNZI
		}else if c.IsLiangDui {
			c.ThType = THPOKER_TYPE_LIANGDUI
		}else if c.IsHulu{
			c.ThType = THPOKER_TYPE_HULU
		}else if c.IsDuiZi{
			c.ThType = THPOKER_TYPE_YIDUI
		}else if c.IsSanTiao{
			c.ThType = THPOKER_TYPE_SANTIAO
		}else if c.IsSiTiao{
			c.ThType = THPOKER_TYPE_SITIAO
		}else {
			c.ThType = THPOKER_TYPE_GAOPAI
		}
	}


	return nil
}


//对牌进行排序
/**
使用sort包需要实现的方法
 */
func ( list CardsList) Len() int{
	return len(list)
}

//------------------------------------------------------------实现扑克牌的排序-----------------------------------------
//由大到小的排序
func ( list CardsList) Less(i,j int) bool{
	if list[i].GetValue() > list[j].GetValue(){	//比较类型
		return true
	}else {
		return false
	}
}


//交换函数
func ( list CardsList) Swap(i,j int){
	var temp *bbproto.Pai = list[i]
	list[i] = list[j]
	list[j] = temp
}

//返回牌的列表
func RandomTHPorkCards(total int ) []*bbproto.Pai{
	result := make([]*bbproto.Pai,total)	//返回值
	indexs := RandomTHPorkIndex(0,52,total)
	for i := 0;i< total;i++ {
		result[i] =  bbproto.CreatePorkByIndex(indexs[i])
	}
	return result
}
//----------------------------------------------------------实现扑克牌的排序结束----------------------------------------


//--------------------------------------------------------------实现德州牌的排序---------------------------------------

func (list ThCardsList) Len() int{
	return len(list)
}

//判断是否大于
func (list ThCardsList) Less(i,j int) bool{
	if list[i].ThType > list[j].ThType {
		return true
	}else if list[i].ThType == list[j].ThType{
		flag := true
		for m := 0; m < len(list[i].KeyValue) ; m++ {
			if list[i].KeyValue[m] < list[j].KeyValue[m] {
				flag = false
				break
			}
		}
		return flag
	}else {
		return false
	}
}

func (list ThCardsList) Swap(i,j int){
	var temp *ThCards = list[i]
	list[i] = list[j]
	list[j] = temp
}

//比较两幅德州扑克牌的大小t1是否大于t2
func ThCompare(t1,t2 *ThCards) int{
	if t1.ThType > t2.ThType {
		return THPOKER_COMPARE_BIG	//返回大于
	}else if t1.ThType == t2.ThType{
		flag := THPOKER_COMPARE_EQUALS
		for m := 0; m < len(t1.KeyValue) ; m++ {
			if t1.KeyValue[m] < t2.KeyValue[m] {
				flag = THPOKER_COMPARE_SMALL
				break
			}else if t1.KeyValue[m] > t2.KeyValue[m]{
				flag = THPOKER_COMPARE_BIG
				break
			}
		}
		return flag
	}else {
		return THPOKER_COMPARE_SMALL
	}
}

//判断两幅德州的牌的带大小




//------------------------------------------------------------实现德州牌的排序结束--------------------------------------


//随机的德州牌的坐标
func RandomTHPorkIndex(min, max,total int) []int32 {
	result := make([]int32,total);
	count := 0;
	for count < total {
		num := utils.Rand(int32(min),int32(max))
		flag := true;
		for j := 0; j < total; j++ {
			if num == result[j] {
				flag = false;
				break;
			}
		}
		if flag {
			result[count] = num;
			count++;
		}
	}
	return result;
}

//通过手牌,和给定的牌得到最大的德州牌
func GetTHPoker(hand,public []*bbproto.Pai,count int) *ThCards{

	var allCards = make([]*bbproto.Pai,7)
	copy(allCards[0:2],hand)
	copy(allCards[2:],public)
	tcsList := Com(7,count, allCards)
	sort.Sort(tcsList)
	return tcsList[0]

}


// 组合函数// 1-n里面取k个组合
func   Com(n,k int,allCards []*bbproto.Pai) ThCardsList{
	var result ThCardsList
	//判断参数是否正确
	if (n < k || n <= 0 || k <= 0) {
		log.E("n,k数据输入不合理")
		return nil;
	}

	a := make([]int,k+1)
	fg := make([]int,k+1)

	for i:=1;i<=k;i++{
		a[i] = i
		fg[i] = i-k+n
	}

	for ; ;  {
		tcs := NewThCards()

		tcs.Cards = make([]*bbproto.Pai,5)
		for i := 1;i<=k ;i++  {
			//fmt.Print(" --%v-- ",a[i]-1)
			tcs.Cards[i-1] = allCards[a[i]-1]
		}

		//fmt.Println("排序之前的牌",tcs.Cards)
		var cs CardsList = tcs.Cards
		sort.Sort(cs)
		//fmt.Println("排序之后的牌",tcs.Cards)
		tcs.OnInit()	//初始化德州
		//fmt.Println("排序之后的类型",*tcs.ThType)

		//fmt.Println("")

		if result ==nil {
			result = make([]*ThCards,1)
			result[0] = tcs
		}else{
			result = append(result,tcs)
		}

		if a[1] == (n - k + 1) {
			break
		}

		for i := k; i >= 1; i-- {
			if a[i] < fg[i] {
				a[i]++
				for j:=i+1;j<=k ;j++  {
					a[j] = a[j-1] +1;
				}
				break
			}
		}
	}
	return result
}


type AllInJackpotList []*AllInJackpot

//每一次all in 都会生成
type AllInJackpot struct {
	UserId       uint32 //
	Jackpopt     int32  //奖金池子
	ThroundCount int32  //押注轮次
	AllInAmount  int32  //allIn的值
}

func (list AllInJackpotList) Len() int{
	return len(list)
}

//由大到小的排序
func ( list AllInJackpotList) Less(i,j int) bool{
	if list[i].ThroundCount < list[j].ThroundCount{	//比较类型
		return true
	}else if list[i].AllInAmount < list[j].AllInAmount{
		return true
	}else{
		return false
	}
	return true
}


//交换函数
func ( list AllInJackpotList) Swap(i,j int){
	var temp *AllInJackpot = list[i]
	list[i] = list[j]
	list[j] = temp
}