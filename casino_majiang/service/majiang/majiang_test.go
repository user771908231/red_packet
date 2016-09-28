package majiang

import (
	"testing"
	"fmt"
)

func Test(t *testing.T) {
	fmt.Println("开始测试...")
	pairPai := XiPai()
	fmt.Println("开始洗牌得到的牌:")
	for i, p := range pairPai {
		fmt.Println(i, "牌-index-----   ", p.GetIndex())
	}

}
