package mongodb

import (
	"testing"
	"casino_server/msg/bbprotogo"
	"fmt"
)

func TestTemp(t *testing.T) {
	d := &bbproto.Game_TounamentBlindBean{}
	d.Reset()

	fmt.Println("开始:")
	fmt.Println(*d)
}
