package mongodb

import (
	"testing"
	"time"
	"fmt"
)


func TestTemp(t *testing.T) {
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for timeNow := range ticker.C {
			fmt.Println("哈哈哈",timeNow,time.Now())
			a()
		}
	}()
	for ; ;  {
		
	}
}


func a(){
	fmt.Println("开始",time.Now())
	time.Sleep(time.Second*10)
	fmt.Println("结束",time.Now())

}

