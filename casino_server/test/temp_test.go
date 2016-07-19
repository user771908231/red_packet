package mongodb

import (
	"testing"
	"fmt"
	"sync"
)

func TestTemp(t *testing.T) {
	d :=&dog{}
	d.name = "小黄"
	d.say()
	d.sing()

}



type dog struct {
	sync.Mutex
	name string
}

func ( d *dog) say(){
	d.Lock()
	defer d.Unlock()
	//go d.sing()
	fmt.Println("哈哈哈:",d.name)
	d = nil
	if d == nil {
		fmt.Println("d==nil")
	}
}

func ( d *dog) sing(){
	d.Lock()
	defer d.Unlock()
	fmt.Println("sing 哈哈哈:",d.name)
}



