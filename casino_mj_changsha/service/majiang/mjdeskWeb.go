package majiang

import (
	"casino_common/utils/numUtils"
	mjproto	"casino_mj_changsha/msg/protogo"
)

func (d *MjDesk) GetTransferredStatus() string {
	ret := ""
	switch d.GetStatus() {
	case MJDESK_STATUS_DINGQUE:
		ret = "开始定缺"
	case MJDESK_STATUS_EXCHANGE:
		ret = "开始换牌"
	case MJDESK_STATUS_READY:
		ret = "开始准备"
	case MJDESK_STATUS_RUNNING:
		ret = "定缺后开始打牌"
	case MJDESK_STATUS_QISHOUHU:
		ret = "起手胡阶段"
	default:

	}
	return ret
}

func (d *MjDesk) GetDeskMJInfo() string {
	if d == nil || d.AllMJPai == nil {
		return "暂时没有初始化麻将"
	}
	s := ""
	for i, p := range d.AllMJPai {
		is, _ := numUtils.Int2String(int32(i))
		s = s + " (" + is + "-" + p.LogDes() + ")"
		if (i+1)%27 == 0 {

		}
	}
	return s
}

func (d *MjDesk) GetTransferredRoomType() string {
	ret := ""
	switch mjproto.MJRoomType(d.GetMjRoomType()) {
	case mjproto.MJRoomType_roomType_xueZhanDaoDi:
		ret = "血战到底"
	case mjproto.MJRoomType_roomType_sanRenLiangFang:
		ret = "三人两房"
	case mjproto.MJRoomType_roomType_siRenLiangFang:
		ret = "四人两房"
	case mjproto.MJRoomType_roomType_deYangMaJiang:
		ret = "德阳麻将"
	case mjproto.MJRoomType_roomType_daoDaoHu:
		ret = "倒倒胡"
	case mjproto.MJRoomType_roomType_xueLiuChengHe:
		ret = "血流成河"
	case mjproto.MJRoomType_roomType_liangRenLiangFang:
		ret = "两人两房"
	case mjproto.MJRoomType_roomType_liangRenSanFang:
		ret = "两人三房"
	case mjproto.MJRoomType_roomType_sanRenSanFang:
		ret = "三人三房"
	default:

	}
	return ret
}
