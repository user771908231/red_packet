package majiang

import (
	"testing"
	"fmt"
	"casino_majiang/conf/config"
)

func Test(t *testing.T) {
	config.InitConfig(false)

	pais := []*MJPai{}

	//index:13 flower:3 value:4 des:"T_4"
	//index:14 flower:3 value:4 des:"T_4"
	//index:15 flower:3 value:4 des:"T_4"
	// index:16 flower:3 value:5 des:"T_5"
	// index:17 flower:3 value:5 des:"T_5"
	// index:18 flower:3 value:5 des:"T_5"
	// index:19 flower:3 value:5 des:"T_5"
	// index:20 flower:3 value:6 des:"T_6"
	// index:21 flower:3 value:6 des:"T_6"
	// index:22 flower:3 value:6 des:"T_6"
	// index:23 flower:3 value:6 des:"T_6"
	// index:24 flower:3 value:7 des:"T_7"
	// index:25 flower:3 value:7 des:"T_7"
	// index:12 flower:3 value:4 des:"T_4" ]

	p1 := InitMjPaiByIndex(13)
	p2 := InitMjPaiByIndex(14)
	p3 := InitMjPaiByIndex(15)
	p4 := InitMjPaiByIndex(16)
	p5 := InitMjPaiByIndex(17)
	p6 := InitMjPaiByIndex(18)
	p7 := InitMjPaiByIndex(19)
	p8 := InitMjPaiByIndex(20)
	p9 := InitMjPaiByIndex(21)
	p10 := InitMjPaiByIndex(22)
	p11 := InitMjPaiByIndex(23)
	p12 := InitMjPaiByIndex(24)
	p13 := InitMjPaiByIndex(25)
	p14 := InitMjPaiByIndex(12)

	pais = append(pais, p1)
	pais = append(pais, p2)
	pais = append(pais, p3)
	pais = append(pais, p4)
	pais = append(pais, p5)
	pais = append(pais, p6)
	pais = append(pais, p7)
	pais = append(pais, p8)
	pais = append(pais, p9)
	pais = append(pais, p10)
	pais = append(pais, p11)
	pais = append(pais, p12)
	pais = append(pais, p13)
	pais = append(pais, p14)

	counts := GettPaiStats(pais)
	canHu, isAll19 := tryHU(counts, len(pais))
	fmt.Println("canHu", canHu)
	fmt.Println("isAll19", isAll19)
	fmt.Println("\n")

}
