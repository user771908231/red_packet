package majiang

import (
	"testing"
	"fmt"
)

func Test(t *testing.T) {
	fmt.Println("开始测试...")
	pairPai := XiPai()
	fmt.Println("开始洗牌得到的牌:")

	for _, p := range pairPai {
		ei := 0
		for _, p2 := range pairPai {
			if p.GetIndex() == p2.GetIndex() {
				ei ++;
				if ei > 1 {
					fmt.Println("出错了...index", p.GetIndex())
				}
			}

		}
	}
	fmt.Println("出错判定完毕...")

	for i, p := range pairPai {
		fmt.Println(i, "牌-index-----   ", p.GetIndex())
	}

}
