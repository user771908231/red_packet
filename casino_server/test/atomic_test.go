package test

import (
	"testing"
	"fmt"
	"sync/atomic"
)

func TestAtomic(t *testing.T){

	fmt.Println("")
	var a int32 = 100
	var b int32 = 100
	var coin int32 = 10
	atomic.AddInt32(&a,coin)
	atomic.AddInt32(&b,-coin)

	fmt.Println("a,b",a,b)


}