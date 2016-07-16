package mongodb

import (
	"testing"
	"fmt"
)

func TestTemp(t *testing.T) {
	var a *int32
	var b int32 = 999;
	a = &b
	fmt.Println("b:",b)
	fmt.Println("a:",a)
	fmt.Println("*a:",*a)
	*a ++
	fmt.Println("*a ++:",*a)
	fmt.Println("a ++:",a)
}

