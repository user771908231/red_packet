package mongodb

import (
	"testing"
	"fmt"
)

type CCC struct {
	I1 int
	I2 int
	S1 string
	S2 string

}

func TestTemp(t *testing.T) {

	var c = &CCC{}
	c.I1 = 1
	c.I2 = 1
	c.S1 = "a"
	c.S2 = "b"

	fmt.Println(c.I1)
	fmt.Println(c.I2)
	fmt.Println(c.S1)
	fmt.Println(c.S2)


}


