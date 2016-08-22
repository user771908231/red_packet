package test

import (
	"testing"
	"fmt"
	"sync/atomic"
	"time"
)

func TestAtomic(t *testing.T){

	fmt.Println("")
	var a int32 = 0
	var b int32 = 10
	fmt.Println("a,b,",a,b)
	atomic.AddInt32(&a,1)
	fmt.Println("a,b,",a,b)
	atomic.AddInt32(&a,1)
	atomic.AddInt32(&a,1)
	fmt.Println("a,b,",a,b)
	atomic.AddInt32(&a,100)
	fmt.Println("a,b,",a,b)
	atomic.AddInt32(&a,-50)
	fmt.Println("a,b,",a,b)
	atomic.AddInt32(&b,-5)
	fmt.Println("a,b,",a,b)

	for i:=0;i<10000;i++ {
		go func() {
			atomic.AddInt32(&a,1)
			//a = a+1
		}()
	}

	time.Sleep(time.Second*5)
	fmt.Println("a,b,",a,b)


}