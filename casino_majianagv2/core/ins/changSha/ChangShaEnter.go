package changSha

import (
	"casino_majiang/msg/protogo"
	"casino_common/common/log"
	"time"
)

func (d *ChangShaFMJDesk) SendGameInfo(userId uint32, reconnect mjproto.RECONNECT_TYPE) {
	log.T("[%v]用户[%v]进入房间,reconnect[%v]之后", d.DlogDes(), userId, reconnect)

	//获取房间的信息
	gameinfo := d.GetGame_SendGameInfo(userId, reconnect)
	gameinfo.GetDeskGameInfo().GetRoomTypeInfo().ChangShaPlayOptions = d.ChangShaPlayOptions //长沙麻将的配置
	d.BroadCastProto(gameinfo)

	//如果是重新进入房间，需要发送重近之后的处理
	if reconnect == mjproto.RECONNECT_TYPE_RECONNECT {
		time.Sleep(time.Second * 3)
		d.SendReconnectOverTurn(userId)
	}
}
