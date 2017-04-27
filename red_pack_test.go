package main

import (
	"testing"
	"casino_admin/model/weixinModel"
	"gopkg.in/mgo.v2/bson"
)

func TestSendRedPack(t *testing.T) {
	err := weixinModel.SendRedPack("ot9ZM00vyjYsDK6oGFKwqn6nsQgs", 1, bson.NewObjectId().Hex())
	t.Log(err)
}
