package redModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_common/utils/db"
	"casino_redpack/model/userModel"
	"sync"
	"casino_common/common/log"
)

type RoomType int

const (
	//五人房红包
	RoomTypeWurenDZ RoomType = 1
	//扫雷（炸弹接龙）
	RoomTypeSaoLei RoomType = 5

	//牛牛
	RoomTypeNiuniu RoomType = 2

	//二八杠
	RoomTypeErBaGang RoomType = 4
)

func init() {

	//炸弹接龙
	Rooms = append(Rooms, &Room{
		ObjId:       bson.NewObjectId(),
		Type:        RoomTypeSaoLei,
		Id:          10000,
		CreatorId:   10000,
		CreatorName: "管理员",
		CreatorHead: "",
		Money:       15,
		MaxUser:     99999,
		Users:       []uint32{},
		CreateTime:  time.Now(),
		RedpackList: []*Redpack{},
	})

	//五人对战
	Rooms = append(Rooms, &Room{
		ObjId: bson.NewObjectId(),
		Type: RoomTypeWurenDZ,
		Id: 10001,
		CreatorId: 10000,
		CreatorName: "管理员",
		CreatorHead: "",
		Money: 15,
		MaxUser: 99999,
		Users: []uint32{},
		CreateTime: time.Now(),
		RedpackList: []*Redpack{
			&Redpack{
				ObjId: bson.NewObjectId(),
				CreatorUser: 10001,
				CreatorName: "郑秀娣",
				CreatorHead: "http://wx.qlogo.cn/mmopen/ajNVdqHZLLDR9YkFYEz0XhumSbNtrpn98PlbDp7K87CxAGYMhkRwV6LEiaYPNRftBoktV2yXTQlodYEUA7SpZkg/0",
				Id: 100004,
				Money: 5,
				Lost: 5,
				Piece: 5,
				Type: RoomTypeWurenDZ,
				TailNumber: 0,
				RoomId: 10001,
				OpenRecord: []*OpenRecordItem{},
				Time: time.Now(),
			},
		},
	})

	//牛牛
	Rooms = append(Rooms, &Room{
		ObjId: bson.NewObjectId(),
		Type: RoomTypeNiuniu,
		Id: 10002,
		CreatorId: 10000,
		CreatorName: "管理员",
		CreatorHead: "",
		Money: 15,
		MaxUser: 99999,
		Users: []uint32{},
		CreateTime: time.Now(),
		RedpackList: []*Redpack{
			&Redpack{
				ObjId: bson.NewObjectId(),
				CreatorUser: 10001,
				CreatorName: "李易峰",
				CreatorHead: "http://wx.qlogo.cn/mmopen/ajNVdqHZLLDR9YkFYEz0XhumSbNtrpn98PlbDp7K87CxAGYMhkRwV6LEiaYPNRftBoktV2yXTQlodYEUA7SpZkg/0",
				Id: 100005,
				Money: 5,
				Lost: 5,
				Piece: 5,
				Type: RoomTypeNiuniu,
				TailNumber: 0,
				RoomId: 10002,
				OpenRecord: []*OpenRecordItem{},
				Time: time.Now(),
			},
		},
	})

	//二八杠
	Rooms = append(Rooms, &Room{
		ObjId: bson.NewObjectId(),
		Type: RoomTypeErBaGang,
		Id: 10003,
		CreatorId: 10000,
		CreatorName: "管理员",
		CreatorHead: "",
		Money: 15,
		MaxUser: 99999,
		Users: []uint32{},
		CreateTime: time.Now(),
		RedpackList: []*Redpack{
			&Redpack{
				ObjId: bson.NewObjectId(),
				CreatorUser: 10005,
				CreatorName: "啦啦啦",
				CreatorHead: "http://wx.qlogo.cn/mmopen/ajNVdqHZLLDR9YkFYEz0XhumSbNtrpn98PlbDp7K87CxAGYMhkRwV6LEiaYPNRftBoktV2yXTQlodYEUA7SpZkg/0",
				Id: 100009,
				Money: 5,
				Lost: 5,
				Piece: 5,
				Type: RoomTypeErBaGang,
				TailNumber: 0,
				RoomId: 10003,
				OpenRecord: []*OpenRecordItem{},
				Time: time.Now(),
			},
		},
	})

}

//房间抽象
type Room struct {
	sync.Mutex
	ObjId       bson.ObjectId `bson:"_id"`
	Type        RoomType   //房间类型
	Id          int32      //房间密码
	CreatorId   uint32     //建房者id
	CreatorName string     //建房者昵称
	CreatorHead string     //建房者头像
	Money       float64    //多大金额的金币房
	MaxUser     int      //最大人数
	Users       []uint32   //在线的人
	CreateTime  time.Time  //创建时间
	RedpackList []*Redpack //红包列表
}

//房间列表
var Rooms = []*Room{}

//获取房间
func GetRoomById(roomId int32) *Room {
	for _,r := range Rooms {
		if r.Id == roomId {
			return r
		}
	}
	return nil
}

//通过房间类型查找
func GetRoomByType(roomType RoomType) *Room {
	for _,r := range Rooms {
		if r.Type == roomType {
			return r
		}
	}
	return nil
}

//新建房间
func CreateRoom(room_type RoomType, creator *userModel.User, max_user int) (*Room, error) {
	new_room_id, err := db.GetNextSeq("redpack_room_id")
	if err != nil {
		return nil, err
	}

	new_room := &Room{
		ObjId: bson.NewObjectId(),
		Type:  room_type,
		Id: new_room_id,
		CreatorId: creator.Id,
		CreatorName: creator.NickName,
		CreatorHead: creator.HeadUrl,
		MaxUser: max_user,
		Users: []uint32{},
		CreateTime: time.Now(),
		RedpackList: []*Redpack{},
	}

	//创建房间
	Rooms = append(Rooms, new_room)
	return new_room, nil
}

//发红包
func (room *Room) SendRedpack(creator *userModel.User, money float64, piece int, tail_number int) (*Redpack,error) {
	new_redpack_id, err := db.GetNextSeq("redpack_redpack_id")
	if err != nil {
		return nil, err
	}
	if creator.AccountNumber == "" {
		creator.AccountNumber = creator.NickName
	}
	new_redpack := &Redpack{
		ObjId: bson.NewObjectId(),
		CreatorUser: creator.Id,
		CreatorName: creator.AccountNumber,
		CreatorHead: creator.HeadUrl,
		Id: new_redpack_id,
		Money: money,
		Lost: money,
		Piece: piece,
		Type: room.Type,
		TailNumber: tail_number,
		RoomId: room.Id,
		OpenRecord: []*OpenRecordItem{},
		Time: time.Now(),
	}

	//广播消息
	BroadCast(bson.M{TagRoomId: room.Id}, func(conn WsConn) []byte {
		userId := conn.Get(TagUserId).(uint32)
		return GetClientRedpackListJson([]*Redpack{new_redpack}, userId)
	})

	//上锁
	room.Lock()
	defer room.Unlock()
	//更新到mongo
	defer new_redpack.Upsert()

	//加入列表,保存100条缓存
	room.RedpackList = append(room.RedpackList, new_redpack)
	if length := len(room.RedpackList); length > 100 {
		room.RedpackList = room.RedpackList[length-100:]
	}
	log.T("发红包后 内存中的红包个数%d",len(room.RedpackList))
	return new_redpack, nil
}

//查找红包
func (room *Room) GetRedpackById(red_id int32) *Redpack {
	//先从内存找
	log.T("内存中红包个数：%d",len(room.RedpackList))
	for _,r := range room.RedpackList {
		if r.Id == red_id {
			return r
		}
	}
	log.T("没有找到此红包 ID：%d",red_id)
	return nil
}
