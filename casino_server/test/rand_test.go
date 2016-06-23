package mongodb

import (
	"testing"
	"casino_server/utils"
	"fmt"
)

func TestRand(t *testing.T){
	//for i := 0; i < 100000; i++ {
	//	t.Log(i,config.RandNikeName())
	//}

	for i := 0; i < 100; i++ {
		fmt.Println(utils.Randn(100))
	}

}
