package mongodb

import (
	"testing"
	"fmt"
	"sync"
)

func TestBingfa(t *testing.T) {
	chs := make([]chan int, 10)

	for i := 0; i < 10; i++ {
		u := bfuser{
			userid : 1,
		}
		chs[i] = make(chan int)
		go u.add1(chs[i])
	}

	for _, ch := range chs {
		<-ch
	}
	fmt.Println("u1最终的结果-----", amount)

}

var amount int32 = 0;

type bfuser struct {
	sync.Mutex
	userid     int32
	userAmount int32
}

func (u *bfuser) add1(ch chan int) {
	u.Lock()
	defer u.Unlock()
	u.userAmount = amount
	fmt.Println("u1加之前,", u.userAmount)
	u.userAmount += 1
	amount = u.userAmount
	fmt.Println("u1加之后,", u.userAmount, amount)
	ch <- 1
}
