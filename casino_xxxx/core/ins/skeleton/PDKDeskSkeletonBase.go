package skeleton

import "github.com/golang/protobuf/proto"

func (d *PDKDeskSkeleton) BC(p proto.Message) {
	for _, u := range d.Users {
		u.WriteMsg(p)
	}
}
