package wxRobotModel

import (
	"gopkg.in/mgo.v2/bson"
	"casino_common/utils/db"
	"casino_common/common/consts/tableName"
	"casino_common/proto/ddproto"
	"strings"
	"casino_common/common/service/rpcService"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/proto"
)

type T_WECHAT_GROUP_INFO struct {
	ObjId bson.ObjectId `bson:"_id"`
	GroupName string  //群名称
	OwnerId uint32    //群主id
	GameType  string     //游戏类型
}

//通过群名称获取群信息
func GetRroupInfoByName(group_name string) *T_WECHAT_GROUP_INFO {
	row := &T_WECHAT_GROUP_INFO{}
	err := db.C(tableName.DBT_ROBOT_WECHAT_GROUP_INFO).Find(bson.M{"groupname": group_name}, row)
	if err != nil {
		return nil
	}
	return row
}

//是否存在空闲房间
func (t *T_WECHAT_GROUP_INFO) GetFreeRoom(gamer_num int, keywords []string) *ddproto.CommonDeskByAgent {
	list, err := rpcService.GetHall().GetAgentRoomList(context.Background(), &ddproto.HallReqAgentRoomGamingList{UserId: proto.Uint32(t.OwnerId)})
	if err != nil || list == nil {
		return nil
	}
	rooms := list.List
	for _, room := range rooms {
		if room.GetStatus() != 0 || len(room.GetUsers()) >= gamer_num {
			continue
		}

		has_all_keywords := true
		for _, keyword := range keywords {
			if !strings.Contains(room.GetTips(), keyword) {
				has_all_keywords = false
				break
			}
		}

		if  has_all_keywords {
			return room
		}
	}
	return nil
}
