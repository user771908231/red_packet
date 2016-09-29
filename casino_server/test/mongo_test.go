package test

import (
	"testing"
	"casino_server/msg/bbprotogo"
	"casino_server/service/room"
	"fmt"
	"sort"
)

func TestMongoUtils(t *testing.T) {
	r1 := newinfo()
	r2 := newinfo()
	r3 := newinfo()

	*r1.UserId = 1
	*r1.Balance = 4000
	*r1.EndTime = 1475053995025919639

	*r2.UserId = 2
	*r2.Balance = 9000
	*r2.EndTime = 1475053999131486792

	*r3.UserId = 3
	*r3.Balance = 3000
	*r3.EndTime = 1475053996982643335

	var list room.RankList
	list = append(list, r1)
	list = append(list, r2)
	list = append(list, r3)
	fmt.Println("1:", list)
	sort.Sort(list)
	fmt.Println("2:", list)
	sort.Sort(list)
	fmt.Println("3:", list)
	sort.Sort(list)
	fmt.Println("4:", list)
}

func newinfo() *bbproto.CsThRankInfo {
	ret := &bbproto.CsThRankInfo{}
	ret.UserId = new(uint32)
	ret.Balance = new(int64)
	ret.EndTime = new(int64)
	return ret
}