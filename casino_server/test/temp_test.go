package mongodb

import (
	"testing"
	"casino_server/service/fruitService"
	"fmt"
	"time"
)

func TestTemp(t *testing.T) {

	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(1 * time.Second)
		}
	}

}

func uPro(p fruitService.ShuiGuoPro) {
	p.ORANGE_1 = 88888
}

func uProa(p *fruitService.ShuiGuoPro) {
	p.ORANGE_1 = 777777
}
