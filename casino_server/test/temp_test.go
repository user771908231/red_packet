package mongodb

import (
	"testing"
	"casino_server/mode"
	"fmt"
)


func TestTemp(t *testing.T) {

	u := &mode.T_user{}
	//u.Mid = bson.NewObjectId()
	fmt.Println("v:",u.Mid)
	fmt.Println("v:",u.Mid.Hex())
	fmt.Println("v:",u.Mid.String())
	if u.Mid.Hex() == "" {
		fmt.Println("u.mid == nil")
	}else{
		fmt.Println(" u.mid != nil")
	}

}

