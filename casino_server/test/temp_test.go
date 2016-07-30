package mongodb

import (
	"testing"
	"casino_server/service/pokerService"
	"casino_server/common/log"
)

func TestTemp(t *testing.T) {
	initSys()

	var total = 21;
	totalCards := pokerService.RandomTHPorkCards(total)        //得到牌
	log.T("得到的所有牌:[%v]", totalCards)
	//得到所有的牌,前五张为公共牌,后边的每两张为手牌
	PublicPai := totalCards[0:5]
	log.T("得到的公共牌:[%v]", PublicPai)

	//给每个人分配手牌
	for i := 0; i < 1; i++ {
		begin := i * 2 + 5
		end := i * 2 + 5 + 2
		hand:= totalCards[begin:end]
		log.T("用户[%v]的手牌[%v]",i,hand)
		th := pokerService.GetTHPoker(hand, PublicPai, 5)
		log.T("i[%v]th[%v]", i,th)
	}

	for ; ;  {
		
	}
}
