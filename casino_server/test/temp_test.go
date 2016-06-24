package mongodb

import (
	"testing"
	"casino_server/service/fruitService"
)

func TestTemp(t *testing.T){

	p := &fruitService.SGJLuckP
	p.GetMax()


}

func uPro(p fruitService.ShuiGuoPro){
	p.ORANGE_1 = 88888
}

func uProa(p *fruitService.ShuiGuoPro){
	p.ORANGE_1 = 777777
}
