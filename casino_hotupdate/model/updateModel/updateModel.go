package updateModel

import (
	"casino_hotupdate/conf"
	"casino_hotupdate/utils"
)

//包裹路径
func GetPkgPath(pkg_name string) string {
	return conf.Server.UploadPath + "/" + pkg_name
}

//包裹是否存在
func IsPkgExist(pkg_name string) bool {
	return utils.FileIsExist(GetPkgPath(pkg_name))
}

//检查包裹名称是否正确
func CheckPkgName(pkg_name string) bool {

	return true
}
