package skeleton

import (
	"casino_common/common/log"
	"casino_majianagv2/core/data"
	"casino_majianagv2/core/api"
	"github.com/golang/protobuf/proto"
)

//常见的get set 方法 需要放置在这里
func (f *SkeletonMJDesk) GetMJConfig() *data.SkeletonMJConfig {
	log.Debug("玩家[%v]进入fdesk")
	return f.config
}

//得到麻将的Status
func (r *SkeletonMJDesk) GetStatus() *data.MjDeskStatus {
	return r.status
}

//todo 日志信息
func (r *SkeletonMJDesk) DlogDes() string {
	return "todo"
}

//todo 通过userId 找到对应的User
func (r *SkeletonMJDesk) GetUserByUserId(userId uint32) api.MjUser {
	return nil
}

//todo 广播
func (d *SkeletonMJDesk) BroadCastProto(p proto.Message) {

}

//todo 是否有牌可以摸
func (d *SkeletonMJDesk) HandPaiCanMo() bool {
	return false
}
func (d *SkeletonMJDesk) GetCheckCase() *data.CheckCase {
	return d.CheckCase
}

// todo 游戏中玩家的人数
func (d *SkeletonMJDesk) GetGamingCount() int32 {
	return 0
}
