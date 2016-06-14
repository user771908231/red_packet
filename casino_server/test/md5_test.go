package mongodb

import (
	"testing"
	"casino_server/utils/security"
)

func TestMd5(t *testing.T){
	s := "abcd"
	data := []byte(s)
	security.Md5(data)
}
