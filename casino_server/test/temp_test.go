package mongodb

import (
	"testing"
	"fmt"
	"casino_server/msg/bbprotogo"
)

var a int32 = 20

func TestTemp(t *testing.T) {

	ha := int32(1)
	if ha == int32(bbproto.EProtoId_REG) {
		fmt.Println("哈哈哈")
	}else{
		fmt.Println("gagaga ")
	}

}

