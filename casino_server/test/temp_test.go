package mongodb

import (
	"testing"
	"fmt"
	"casino_server/utils/strUtils"
)


func TestTemp(t *testing.T) {
	s := "abc"
	b := strUtils.Str2Bytes(s)
	fmt.Println(b)

	c := strUtils.Bytes2Str(b)
	fmt.Println(c)

}


