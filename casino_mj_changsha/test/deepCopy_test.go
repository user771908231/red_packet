package test

import (
	"testing"
	"github.com/golang/protobuf/proto"
)

type User struct {
	name  *string
	name2 string
}

func TestDeepCp(t *testing.T) {
	list := make([]*User, 4)
	u1 := &User{
		name: proto.String("u1n"),
		name2:"u1n2",
	}

	u2 := &User{
		name: proto.String("u2n"),
		name2:"u2n2",
	}

	u3 := &User{
		name: proto.String("u3n"),
		name2:"u3n2",
	}

	u4 := &User{
		name: proto.String("u4n"),
		name2:"u4n2",
	}

	list[0] = u1
	list[1] = u2
	list[2] = u3
	list[3] = u4

	list2 := make([]*User, 4)
	t.Logf("&list:%p list :%v ", &list, list)
	t.Logf("&list2:%p list2 :%v ", &list2, list2)

	//list2 = append(list2, list...)
	copy(list2, list)
	t.Logf("&list:%p list :%v ", &list, list)
	t.Logf("&list2:%p list2 :%v ", &list2, list2)

	list2[3].name = proto.String("x-u4n")
	list2[3].name2 = "x-u4n2"
	t.Logf("&list:%p list :%v ", &list, list)
	t.Logf("&list2:%p list2 :%v ", &list2, list2)

	//list3 := util.DeepClone(list)
	//t.Logf("&list3:%p list3 :%v ", &list3, list3)
}
