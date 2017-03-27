package majiang

import (
	"testing"
	"casino_common/common/game"
	"github.com/golang/protobuf/proto"
)

func TestGameUser(t *testing.T) {
	nu := NewMjUser()
	nu.UserId = proto.Uint32(100)
	var u game.GameUserApi = nu
	t.Logf("ID：%v", u.GetUserId())
}