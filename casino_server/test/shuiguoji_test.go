package mongodb

import (
	"testing"
	"net"
	"casino_server/msg/bbproto"
	"casino_server/utils/test"
	"fmt"
)

func TestShuiGuoJi(t *testing.T){

	conn, err := net.Dial(TCP, url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var id uint16 = 7
	var data bbproto.Shuiguoji

	var nApple int32 = 1						//苹果的数量
	data.ScoresApple = &nApple
	data.ScoresOrange = &nApple
	data.Scores77 = &nApple
	data.ScoresBar = &nApple
	data.ScoresMango = &nApple
	data.ScoresStar = &nApple
	data.ScoresWatermelon = &nApple
	data.ScoresBell = &nApple

	m := utils.AssembleData(id,&data)
	conn.Write(m)
	result := utils.Read(conn).(*bbproto.ShuiguojiRes)
	fmt.Println("读取的结果:",result.Result)
}