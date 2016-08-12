package mongodb

import (
	"testing"
	"fmt"
)

var CCC struct {
	I1 int
	I2 int
	S1 string
	S2 string

}

func TestTemp(t *testing.T) {

	CCC.I1 = 1
	CCC.I2 = 1
	CCC.S1 = "a"
	CCC.S2 = "b"

	fmt.Println(CCC.I1)
	fmt.Println(CCC.I2)
	fmt.Println(CCC.S1)
	fmt.Println(CCC.S2)


}


