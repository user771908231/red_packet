package redModel

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"casino_common/utils/db"
	"casino_redpack/model/userModel"
	"sync"
)

type RoomType string

const (
	//五人房红包
	RoomTypeWurenDZ RoomType = "wurendz"
	//炸弹接龙
	RoomTypeZhandanJL RoomType = "zhadanjl"
)

func init() {
	Rooms = append(Rooms, &Room{
		ObjId: bson.NewObjectId(),
		Type: RoomTypeZhandanJL,
		Id: 10000,
		CreatorId: 10000,
		CreatorName: "管理员",
		CreatorHead: "",
		Money: 15,
		MaxUser: 999,
		Users: []uint32{},
		CreateTime: time.Now(),
		RedpackList: []*Redpack{},
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

//单个红包
type Redpack struct {
	ObjId       bson.ObjectId `bson:"_id"`
	CreatorUser uint32            //发包的人id
	CreatorName string            //发包人昵称
	CreatorHead string            //发包人头像
	Id          int32             //红包id
	Money       float64           //红包大小
	Lost        float64           //剩余大小
	Piece       int             //分成几份
	Type        RoomType          //红包类型
	TailNumber  int             //尾数 0-9, 默认-1
	RoomId      int32             //房间id
	OpenRecord  []*OpenRecordItem //开红包记录
	Time        time.Time         //发红包时间
}

//开红包记录
type OpenRecordItem struct {
	UserId uint32  //领红包的人id
	NickName string  //领红包的人的昵称
	Head string  //领红包的人的头像
	Money float64  //领了多少钱
	Time time.Time  //时间
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
	new_redpack := &Redpack{
		ObjId: bson.NewObjectId(),
		CreatorUser: creator.Id,
		CreatorName: creator.NickName,
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

	//上锁
	room.Lock()
	defer room.Unlock()

	//加入列表
	room.RedpackList = append(room.RedpackList, new_redpack)
	if length := len(room.RedpackList); length > 10 {
		room.RedpackList = room.RedpackList[length-10:]
	}

	//广播消息
	BroadCast(GetClientRedpackListJson([]*Redpack{new_redpack}), bson.M{TagRoomId: room.Id})
	return new_redpack, nil
}
