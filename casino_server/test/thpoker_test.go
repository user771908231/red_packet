package mongodb

import (
	"testing"
	"casino_server/service/pokerService"
	"fmt"
)

func TestThPoker(t *testing.T){

	fmt.Println("")
	var total = 21
	totalCards := pokerService.RandomTHPorkCards(total)        //得到牌
	fmt.Println("所有牌%v",totalCards)

	//所有的牌
	public := totalCards[0:5]
	fmt.Println("公共牌%v",public)

	hand   := totalCards[5:7]
	fmt.Println("手牌%v",hand)

	thpoker :=  pokerService.GetTHPoker(hand,public,5)
	fmt.Println("德州牌%v,牌的类型%v",thpoker.Cards,*thpoker.ThType)
}