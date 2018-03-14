package main

import (
	"testing"
	"casino_laowangye/service/niuniu"
	"casino_common/proto/ddproto"
	"casino_common/utils/chessUtils"
	"log"
	"math/rand"
	"time"
	"github.com/golang/protobuf/proto"
)

func ShowDesk(desk *laowangye.Desk) {
	for _,u := range desk.Users {
		log.Println(*u)
	}
}

func TestPokerList(t *testing.T) {
	poker_list := laowangye.PokerList
	//poker_list_no_jqk := niuniu.PokerListNoJQK

	t.Log(poker_list)
	//t.Log(poker_list_no_jqk)
}

// 测试： 老王爷牌型解析
func TestParsePoker(t *testing.T)  {
	niu_map := map[ddproto.LwyEnum_PokerType]int{}
	option := &ddproto.LwyDeskOption{
		IsFlowerPlay: proto.Bool(true),
	}
	for i := 0; i< 100000; i++ {
		<-time.After(1*time.Nanosecond)
		index := chessUtils.Xipai(0, 52)
		pais := []*ddproto.CommonSrvPokerPai{}
		for j := 0;j <5; j++ {
			pais = append(pais, laowangye.PokerList[int(index[j])])
		}

		user_poker := laowangye.ParseNiuPoker(pais, option)
		//如果出现一条龙、老王爷等牌型则重新洗牌
		rand_num := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)
		if poker_rand,ok := laowangye.Exchange_poker_score_map[user_poker.GetType()];ok && rand_num<poker_rand && len(index)-((1-1+1)*5)>=2 {
			//需要换牌
			user_poker.Pais[1] = laowangye.PokerList[int(index[len(index)-1])]
			index = index[:len(index)-1]
			user_poker.Pais[3] = laowangye.PokerList[int(index[len(index)-1])]
			index = index[:len(index)-1]
			user_poker = laowangye.ParseNiuPoker(user_poker.Pais, option)
		}

		if _,ok := niu_map[user_poker.GetType()];ok {
			niu_map[user_poker.GetType()]++
		}else {
			niu_map[user_poker.GetType()] = 1
		}
		if user_poker.GetType() == ddproto.LwyEnum_PokerType_NIU_ZHA_DAN {
			//t.Log(user_poker.GetType(), user_poker)
		}
	}
	//t.Log(niu_map)
}

//测试：无花牌
func TestFiltFlowerPoker(t *testing.T) {
	// 发牌
	rand_index := chessUtils.Xipai(0, 51)

	rand_index = []int32{48, 9, 27, 43, 13, 35, 38, 34, 31, 20, 14, 32, 19, 8, 24, 49, 28, 45, 41, 1, 44, 26, 12, 6, 17, 23, 16, 10, 5, 39, 7, 2, 47, 22, 37, 36, 18, 40, 4, 0, 21, 42, 11, 33, 30, 50, 15, 25, 46, 3, 29}

	t.Log(rand_index)
	rand_filter := []int32{}
	//如果没有花牌，则把j、q、k过滤掉（移动到尾部）

	for _,index := range rand_index {
		switch index {
		case 8, 9, 10, 21, 22, 23, 34, 35, 36, 47, 48, 49:
			rand_filter = append(rand_filter, index)
		}
	}

	t.Log(rand_index)
	for i,index := range rand_index {
		switch index {
		case 8, 9, 10, 21, 22, 23, 34, 35, 36, 47, 48, 49:
			t.Log("2", i, index)
		}
	}

}

//测试： 老王爷发牌
func TestSendPoker(t *testing.T) {
	//desk := GetDesk()
	//desk.AddUser(1, nil)
	//desk.AddUser(2, nil)
	//desk.DoSendPoker()
	//ShowDesk(desk)
}

//根据id初始化一副牌
func InitPokerByIds(index ...int32) *ddproto.LwySrvPoker {
	option := &ddproto.LwyDeskOption{
		IsFlowerPlay: proto.Bool(true),
	}
	return laowangye.ParseNiuPoker([]*ddproto.CommonSrvPokerPai{
		laowangye.InitPaiByIndex(index[0]),
		laowangye.InitPaiByIndex(index[1]),
		laowangye.InitPaiByIndex(index[2]),
		laowangye.InitPaiByIndex(index[3]),
		laowangye.InitPaiByIndex(index[4]),
	}, option)
}

//测试：老王爷比牌
func TestIsBigThanBanker(t *testing.T) {
	banker_poker := InitPokerByIds(38, 2, 27, 42, 41)
	user_poker := InitPokerByIds(15, 40, 4, 46, 20)
	t.Log(banker_poker.GetType(), user_poker.GetType(), laowangye.IsBigThanBanker(banker_poker, user_poker))

	//老王爷7大
	banker_poker = InitPokerByIds(0, 4, 7, 1, 3)
	//老王爷Q大
	user_poker = InitPokerByIds(0, 30, 27, 42, 48)
	t.Log(banker_poker.GetType(), user_poker.GetType(), laowangye.IsBigThanBanker(banker_poker, user_poker))
}
