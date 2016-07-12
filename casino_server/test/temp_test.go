package mongodb

import (
	"testing"
	"fmt"
)

func TestTemp(t *testing.T) {
	var a *int32
	a = new(int32)
	*a = 1232

	fmt.Println(a)
	fmt.Println(a)

}

