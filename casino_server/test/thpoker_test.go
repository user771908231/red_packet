package mongodb

import (
	"testing"
	"casino_server/service/pokerService"
	"casino_server/msg/bbprotogo"
	"fmt"
	"casino_server/common/log"
)

func TestThPoker(t *testing.T) {
	initSys()
	//testCompare()
	initCardTest()
}



// testCompare


func testCompare() {
	//1,构造两幅德州牌

	fmt.Println()
	//th1index := []int32{23,49,9,48,43}
	th1index := []int32{26, 36, 34, 47, 44}

	th1 := pokerService.NewThCards()
	th1.Cards = make([]*bbproto.Pai, 5)
	for i := 0; i < 5; i++ {
		th1.Cards[i] = bbproto.CreatePorkByIndex(th1index[i])
	}
	th1.OnInit()
	fmt.Println("th1的牌型[%v],key值[%v],牌:[%v]", th1.ThType, th1.KeyValue, th1.Cards)

	//th2index := []int32{23,49,11,50,43}
	th2index := []int32{25, 50, 36, 34, 47}
	th2 := pokerService.NewThCards()
	th2.Cards = make([]*bbproto.Pai, 5)
	for i := 0; i < 5; i++ {
		th2.Cards[i] = bbproto.CreatePorkByIndex(th2index[i])
	}
	th2.OnInit()

	fmt.Println("th2,牌型[%v],key值[%v],牌:[%v]", th2.ThType, th2.KeyValue, th2.Cards)

	ret1 := pokerService.ThCompare(th1, th2)
	fmt.Println("比较的结果:[%v] //1大,2小,3相等:--", ret1)
}

func initCardTest() {
	//轮子13 2 17 3 41 1 20 35 24 29 40 33 47 27 18 39 37 36 51 4 16
	wheel := []int32{13, 2, 17, 3, 41, 1, 20, 35, 24, 29, 40, 33, 47, 27, 18, 39, 37, 36, 51, 4, 16}
	testRandomCards(wheel)

}

func testRandomCards(randomIndex []int32) {

	var total = 21; //人数*手牌+5张公共牌
	totalCards := pokerService.RandomTHPorkCards(total)        //得到牌
	//得到所有的牌,前五张为公共牌,后边的每两张为手牌
	publicPai := totalCards[0:5]
	log.T("初始化得到的公共牌的信息:")
	//给每个人分配手牌
	for i := 0; i < 5; i++ {
		begin := i * 2 + 5
		end := i * 2 + 5 + 2
		cards := totalCards[begin:end]
		log.T("用户[%v]的手牌[%v]", i, cards)
		thCards := pokerService.GetTHPoker(cards, publicPai, 5)
		log.T("用户[%v]的:拍类型,所有牌[%v],th[%v]", i, thCards.ThType, thCards.Cards, thCards)
	}
	log.T("开始一局新的游戏,初始化牌的信息完毕...")

}