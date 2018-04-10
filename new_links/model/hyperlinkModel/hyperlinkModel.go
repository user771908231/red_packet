package hyperlinkModel

import (
	"new_links/model/groupingModel"
	"new_links/model/linksModel"
)

func GetGroup(string string) string {
	G := groupingModel.GetGroupHost(string)
	val := linksModel.RandLink(G.ObjId)
	defer func() {
		val.Visssts()
	}()
	if val == nil {
		return ""
	}
	return val.LinkName
}

func RandMap()  {

}