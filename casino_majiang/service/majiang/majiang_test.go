package majiang

import (
	"testing"
	"fmt"
)

func Test(t *testing.T) {
	fmt.Println("开始测试...")
	pairPai := XiPai()
	p1 := make([]*MJPai, 13)
	p2 := make([]*MJPai, 13)

	fmt.Println("一副", pairPai)
	fmt.Println("p1", p1)
	fmt.Println("p2", p2)

	p1 = pairPai[0:3]
	p2 = pairPai[3:6]
	fmt.Println("一副", pairPai)
	fmt.Println("p1", p1)
	fmt.Println("p2", p2)

	p1 = append(p1, pairPai[26])

	fmt.Println("一副", pairPai)
	fmt.Println("p1", p1)
	fmt.Println("p2", p2)

	p2 = append(p2[:1], p2[2:]...)
	fmt.Println("一副", pairPai)
	fmt.Println("p1", p1)
	fmt.Println("p2", p2)
}
