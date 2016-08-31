package test

import (
	"testing"
	"fmt"
	"casino_server/utils/security"
)

func TestTemp(t *testing.T) {
	//testN()
	fmt.Println("md5key:",security.SECRET_KEY)

}

type s struct {
	i *int32
}

func testN() {
	s1 := &s{}
	s1.i = new(int32)
	*s1.i = int32(99)
	fmt.Println("s1.i", *s1.i)
	i := s1.i
	fmt.Println("i,", *i)
	s1 = nil
	fmt.Println("i,", *i)
}

func testMongoUtils2() {
}




