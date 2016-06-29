package mongodb

import (
	"testing"
	"fmt"
)

var a int32 = 20

func TestTemp(t *testing.T) {
	fmt.Println(a)
	fmt.Println(&a)
}

