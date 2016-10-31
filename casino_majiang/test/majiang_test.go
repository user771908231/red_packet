package test

import (
	"testing"
	"casino_majiang/service/majiang"
	"fmt"
	"casino_server/utils/numUtils"
)

func Test(t *testing.T) {
	testCanGangPai()
}


//测试ganGang

/**

玩家的手牌[15-4筒	 16-5筒	 17-5筒	 18-5筒	 19-5筒	 20-6筒	 21-6筒	 22-6筒	 23-6筒	 24-7筒	 12-4筒	 ]
玩家的碰牌[12-4筒	 13-4筒	 14-4筒	 ],玩家的杠牌[],玩家的胡牌[],玩家的inpai[6条]
bill[玩家[3302]总共输赢[0],下边是细节:

 */

func testCanGangPai() {

	//pai := majiang.InitMjPaiByIndex()
	var handPai *majiang.MJHandPai = &majiang.MJHandPai{}
	//初始化手牌
	handPai.Pais = append(handPai.Pais, majiang.InitMjPaiByIndex(15))
	handPai.Pais = append(handPai.Pais, majiang.InitMjPaiByIndex(16))
	handPai.Pais = append(handPai.Pais, majiang.InitMjPaiByIndex(17))
	handPai.Pais = append(handPai.Pais, majiang.InitMjPaiByIndex(18))
	handPai.Pais = append(handPai.Pais, majiang.InitMjPaiByIndex(19))
	handPai.Pais = append(handPai.Pais, majiang.InitMjPaiByIndex(20))
	handPai.Pais = append(handPai.Pais, majiang.InitMjPaiByIndex(21))
	handPai.Pais = append(handPai.Pais, majiang.InitMjPaiByIndex(22))
	handPai.Pais = append(handPai.Pais, majiang.InitMjPaiByIndex(23))
	handPai.Pais = append(handPai.Pais, majiang.InitMjPaiByIndex(24))

	handPai.InPai = majiang.InitMjPaiByIndex(25)

	handPai.PengPais = append(handPai.PengPais, majiang.InitMjPaiByIndex(12))
	handPai.PengPais = append(handPai.PengPais, majiang.InitMjPaiByIndex(13))
	handPai.PengPais = append(handPai.PengPais, majiang.InitMjPaiByIndex(14))
	fmt.Println("")
	fmt.Println("处理之前的手牌:", getPaiInfo(handPai.Pais))
	fmt.Println("处理之前的碰牌:", getPaiInfo(handPai.PengPais))
	fmt.Println("处理之前的杠牌:", getPaiInfo(handPai.GangPais))

	can, result := handPai.GetCanGang(nil)

	fmt.Println(can)
	fmt.Println(result)
	fmt.Println("result:", getPaiInfo(result))

	fmt.Println("处理之后的手牌:", getPaiInfo(handPai.Pais))
	fmt.Println("处理之后的碰牌:", getPaiInfo(handPai.PengPais))
	fmt.Println("处理之后的杠牌:", getPaiInfo(handPai.GangPais))

}

func getPaiInfo(pais []*majiang.MJPai) string {
	s := ""
	for _, p := range pais {
		ii, _ := numUtils.Int2String(int32(p.GetIndex()))
		s = s + ii + "-" + p.LogDes() + "\t "
	}

	return s

}


