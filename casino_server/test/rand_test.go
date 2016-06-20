package mongodb

import (
	"testing"
	"casino_server/common/config"
)

func TestRand(t *testing.T){
	for i := 0; i < 100000; i++ {
		t.Log(i,config.RandNikeName())
	}
}
