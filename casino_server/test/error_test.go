package test

import (
	"testing"
	"fmt"
)

func TestError(t *testing.T) {
	s := testErrorGetString()
	fmt.Println("s:", s)

}

func testErrorGetString() (ret string) {
	//测试捕获异常
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[W]", r)
			ret = "err1"
		}
	}()

	a := 0
	b := 1 / a
	fmt.Println("b", b)
	ret = "ok1"
	return ret
}
