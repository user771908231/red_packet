package a

import (
	"fmt"
)

type A struct {
	name string
}

func (a *A) Aa(){
	fmt.Println("aa")
}

