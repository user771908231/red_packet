package mongodb

import (
	"testing"
	"fmt"
	"casino_server/service/pokerService"
)

func TestTemp(t *testing.T) {
	var total = 21;
	totalCards := pokerService.RandomTHPorkCards(total)        //得到牌
	fmt.Println("得到的所有牌:[%v]", totalCards)
	//得到所有的牌,前五张为公共牌,后边的每两张为手牌
	publicPai := totalCards[0:5]
	fmt.Println("得到的公共牌:[%v]", publicPai)

	//给每个人分配手牌
	for i := 0; i < 5; i++ {
		begin := i * 2 + 5
		end := i * 2 + 5 + 2
		c := totalCards[begin:end]
		fmt.Println("用户[%v]的手牌[%v]", i, c)
	}
}

