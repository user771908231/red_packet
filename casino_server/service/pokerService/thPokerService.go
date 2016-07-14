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

//config


type  ThCardsList 	[]*ThCards	//需要比较的牌
type  CardsList		[]*bbproto.Pai	//对牌进行排序

type ThCards struct {
	ThType		*int32
	KeyValue	[]int32
	Cards	[]*bbproto.Pai
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

func (list ThCardsList) Less(i,j int) bool{
	return true

}

func (list ThCardsList) Swap(i,j int){
}

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
func GetTHMax(hand,public []*bbproto.Pai,count int) *ThCards{

	//把公共牌增加到手牌中
	for i := 0; i < len(public); i++ {
		if public[i] !=nil {
			hand = append(hand,public[i])
		}
	}
	//log.T("总共有[%v]张牌",len(hand))
	tcsList := Com(7,count,hand)
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
		tcs := &ThCards{}
		tcs.Cards = make([]*bbproto.Pai,5)
		for i := 1;i<=k ;i++  {
			tcs.Cards[1] = allCards[a[i]-1]
		}

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