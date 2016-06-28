package mongodb

import (
	"testing"
)

func TestTemp(t *testing.T) {
	var i1 int32 = 0
	var i2 int = 10
	i1 = int32(i2)
	println(i2,i1)
}

