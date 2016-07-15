package mongodb

import (
	"testing"
	"casino_server/service/userService"
	"fmt"
)

func TestTemp(t *testing.T) {
	u := userService.GetUserById(10005)
	fmt.Println("user:",u)
}

