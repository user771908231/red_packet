package majiang

import (
	"testing"
)

func TestGetJiaoInfos(t *testing.T) {
	gameData := getPinghu()
	t.Logf("初始化的数据:[%v]", ServerPais2string(gameData.Pais))

	//获取用户手牌 包括inPai
	userHandPai := NewMJHandPai()
	*userHandPai = *gameData       //手牌
	t.Logf("gameData.p[%p]", gameData.Pais)
	t.Logf("userHandPai.p[%p]", userHandPai.Pais)

	//userHandPai.Pais = make([]*MJPai, len(gameData.Pais))
	//copy(userHandPai.Pais, gameData.Pais)        //复制一份，避免圆数据背修改

	//
	userPais := make([]*MJPai, len(userHandPai.Pais))
	copy(userPais, userHandPai.Pais)
	if userHandPai.InPai != nil {
		//碰牌 无inPai的情况
		userPais = append(userPais, userHandPai.InPai)
	}

	lenth := len(userPais)
	for i := 0; i < lenth; i++ {
		t.Logf("----------------第[%v]次循环----------------", i)

		t.Logf("删除牌之前 userpais[%v]", ServerPais2string(userPais))
		//从用户手牌中移除当前遍历的元素
		removedPai := userPais[i]
		t.Logf("要删除的牌removedPai[%v]", removedPai.LogDes())
		userPais = removeFromPais(userPais, i)
		t.Logf("删除牌之后 userPais[%v]", ServerPais2string(userPais))
		userHandPai.Pais = userPais
		t.Logf("userHandPai.Pais[%v]", ServerPais2string(userHandPai.Pais))
		t.Logf("gameData.Pais[%v]", ServerPais2string(gameData.Pais))

		//回复手牌
		userPais = addPaiIntoPais(removedPai, userPais, i) //将移除的牌添加回原位置继续遍历
		t.Logf("userPais[%v]", ServerPais2string(userPais))

		t.Logf("----------------第[%v]次循环结束----------------", i)

	}
}

//平胡 1番
func getPinghu() *MJHandPai {
	inPaiDes := "W_9"
	paisDes := []string{"W_7", "W_7", "W_8", "W_8", "W_9", "T_6", "T_6"}
	pengPaisDes := []string{"T_1", "T_1", "T_1"}
	gangPaisDes := []string{"T_3", "T_3", "T_3", "T_3"}

	return getMjHandPai(inPaiDes, pengPaisDes, gangPaisDes, paisDes)

}

func getMjHandPai(inPaiDes string, pengPaisDes []string, gangPaisDes []string, paisDes []string) *MJHandPai {
	hand := NewMJHandPai()

	hand.InPai = InitMjPaiByDes(inPaiDes, hand)

	for i := 0; i < len(paisDes); i++ {
		hand.Pais = append(hand.Pais, InitMjPaiByDes(paisDes[i], hand))
	}
	for i := 0; i < len(pengPaisDes); i++ {
		hand.PengPais = append(hand.PengPais, InitMjPaiByDes(pengPaisDes[i], hand))
	}
	for i := 0; i < len(gangPaisDes); i++ {
		hand.GangPais = append(hand.GangPais, InitMjPaiByDes(gangPaisDes[i], hand))
	}

	//ignore params
	hand.HuPais = nil
	hand.OutPais = nil
	*hand.QueFlower = W        //定缺的花色
	return hand
}