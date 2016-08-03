package mongodb

import (
	"testing"
	"casino_server/service/pokerService"
	"casino_server/msg/bbprotogo"
	"fmt"
)

func TestThPoker(t *testing.T) {
	testCompare()

}



// testCompare


func testCompare() {
	//1,构造两幅德州牌

	fmt.Println()
	//th1index := []int32{23,49,9,48,43}
	th1index := []int32{26,36,34,47,44}

	th1 := pokerService.NewThCards()
	th1.Cards = make([]*bbproto.Pai, 5)
	for i := 0; i < 5; i++ {
		th1.Cards[i] = bbproto.CreatePorkByIndex(th1index[i])
	}
	th1.OnInit()
	fmt.Println("th1的牌型[%v],key值[%v],牌:[%v]", th1.ThType,th1.KeyValue, th1.Cards)

	//th2index := []int32{23,49,11,50,43}
	th2index := []int32{25,50,36,34,47}
	th2 := pokerService.NewThCards()
	th2.Cards = make([]*bbproto.Pai, 5)
	for i := 0; i < 5; i++ {
		th2.Cards[i] = bbproto.CreatePorkByIndex(th2index[i])
	}
	th2.OnInit()

	fmt.Println("th2,牌型[%v],key值[%v],牌:[%v]", th2.ThType,th2.KeyValue, th2.Cards)

	ret1 := pokerService.ThCompare(th1, th2)
	fmt.Println("比较的结果:[%v] //1大,2小,3相等:--", ret1)
}