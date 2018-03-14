package laowangye

import (
	"casino_common/proto/ddproto"
	"casino_common/utils/redisUtils"
	"errors"
	"casino_common/utils/numUtils"
	"casino_common/common/consts"
	"casino_common/common/log"
	"github.com/golang/protobuf/proto"
	"casino_common/common/Error"
)

func (room *Room) GetSnapShotIdList() *ddproto.LwySrvDeskSnapshotIdIndex {
	desk_id_list_key := consts.RKEY_LWY_SNAPSHOT_ID_LIST
	list := redisUtils.GetObj(desk_id_list_key, &ddproto.LwySrvDeskSnapshotIdIndex{})
	if list == nil {
		return &ddproto.LwySrvDeskSnapshotIdIndex{
			DeskId: []int32{},
		}
	}
	return list.(*ddproto.LwySrvDeskSnapshotIdIndex)
}

//牌桌创建快照
func (desk *Desk) NewSnapShot() error {
	return nil
	if desk == nil {
		return errors.New("desk nil.")
	}
	log.T("房间%d创建快照id。", desk.GetDeskId())
	desk_id_list_key := consts.RKEY_LWY_SNAPSHOT_ID_LIST

	list := desk.Room.GetSnapShotIdList()

	for _, ex_desk_id := range list.DeskId {
		if ex_desk_id == desk.GetDeskId() {
			//已存在快照
			return nil
		}
	}

	list.DeskId = append(list.DeskId, desk.GetDeskId())
	return redisUtils.SetObj(desk_id_list_key, list)
}

//牌桌删除快照
func (desk *Desk) RemoveSnapShot() error {
	return nil
	log.T("开始删除牌桌pwd:%v id:%v快照.", desk.GetPassword(), desk.GetDeskId())
	if desk == nil {
		return errors.New("desk nil.")
	}

	desk_id_list_key := consts.RKEY_LWY_SNAPSHOT_ID_LIST
	list := desk.Room.GetSnapShotIdList()
	for i, ex_desk_id := range list.DeskId {
		if ex_desk_id == desk.GetDeskId() {
			//删除快照id
			list.DeskId = append(list.DeskId[:i], list.DeskId[i+1:]...)
			redisUtils.SetObj(desk_id_list_key, list)
			log.T("牌桌%d删除快照id", desk.GetDeskId())
			//删除快照
			desk_snap_key := redisUtils.K_STRING(consts.RKEY_LWY_SNAPSHOT_DATA, numUtils.Int2String2(desk.GetDeskId()))
			redisUtils.Del(desk_snap_key)
			log.T("已删除牌桌%d的快照", desk.GetDeskId())
			return nil
		}
	}

	log.E("删除牌桌pwd:%v id:%v快照失败,未找到该快照。", desk.GetPassword(), desk.GetDeskId())
	return nil
}

//备份牌桌快照
func (desk *Desk) WipeSnapShot() error {
	return nil
	go func() error {
		defer Error.ErrorRecovery("WipeSnapShot()")
		if desk == nil {
			return errors.New("desk nil.")
		}
		//如果该房间已被销毁，则不执行备份快照操作(避免多个defer导致销毁房间后重复备份的问题)
		if _, err := desk.Room.GetDeskById(desk.GetDeskId()); err != nil {
			return err
		}

		log.T("牌桌%d正在备份快照。", desk.GetDeskId())
		new_snapshot := &ddproto.LwySrvDeskSnapshot{
			DeskState: desk.LwySrvDesk,
			Users: []*ddproto.LwySrvUser{},
		}

		for _,u := range desk.Users {
			new_snapshot.Users = append(new_snapshot.Users, u.LwySrvUser)
		}

		desk_snap_key := redisUtils.K_STRING(consts.RKEY_LWY_SNAPSHOT_DATA, numUtils.Int2String2(desk.GetDeskId()))

		return redisUtils.SetObj(desk_snap_key, new_snapshot)
	}()
	return nil
}

//恢复牌桌快照
func (room *Room) RecoveryDeskSnapShot(desk_id int32) {
	return
	log.T("正在恢复房间%d的快照。", desk_id)
	desk_snap_key := redisUtils.K_STRING(consts.RKEY_LWY_SNAPSHOT_DATA, numUtils.Int2String2(desk_id))

	redis_snap_shot := redisUtils.GetObj(desk_snap_key, &ddproto.LwySrvDeskSnapshot{})
	if redis_snap_shot == nil {
		//如果该快照不存在，则删除该快照id和对应的快照
		desk := Desk{
			Room: room,
			LwySrvDesk: &ddproto.LwySrvDesk{
				DeskId: proto.Int32(desk_id),
			},
		}
		desk.RemoveSnapShot()
		return
	}

	snap_shot := redis_snap_shot.(*ddproto.LwySrvDeskSnapshot)

	if snap_shot.DeskState.GetRoomId() != room.GetRoomId() {
		//过滤非room的牌桌
		return
	}

	new_desk := &Desk{
		Room: room,
		LwySrvDesk: snap_shot.DeskState,
		Users: []*User{},
	}

	for _,u := range snap_shot.Users {
		new_desk.Users = append(new_desk.Users, &User{
			Agent: nil,
			Desk: new_desk,
			LwySrvUser: u,
		})
	}

	room.Desks = append(room.Desks, new_desk)

	//牌桌恢复切片后的处理
	new_desk.OnAfterRecoveryDo()

}

//牌桌恢复快照后的处理
func (desk *Desk) OnAfterRecoveryDo() {
	//如果当前状态为抢庄抢庄
	if desk.GetStatus() == ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_QIANGZHUANG {
		desk.StartQiangzhuangTimer()
	}

	//如果当前状态为解散房间,则清除解散房间相关的状态
	if desk.GetIsOnDissolve() == true {
		*desk.IsOnDissolve = false
		*desk.DissolveTime = 0

		for _, u := range desk.Users {
			if u != nil {
				*u.DissolveState = 0
			}
		}
	}

	//机器人-准备处理
	for _,u := range desk.Users {
		if !u.GetIsRobot() {
			continue
		}

		switch desk.GetStatus() {
		case ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_READY:
			u.DoRobotReady(false)
		case ddproto.LwyEnumDeskStatus_LWY_DESK_STATUS_WAIT_QIANGZHUANG:
			u.DoRobotQiangzhuang()
		}
		u.Desk.SendQiangzhuangOt()
	}
}
