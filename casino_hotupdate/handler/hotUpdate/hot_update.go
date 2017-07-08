package hotUpdate

import (
	"mime/multipart"
	"casino_hotupdate/utils"
	"casino_hotupdate/modules"
	"casino_hotupdate/model/updateModel"
)

//列表
func ListHandler(ctx *modules.Context)  {
	ctx.HTML(200, "hotupdate/list")
}

//上传表单
type UploadForm struct {
	OverWrite bool
	File *multipart.FileHeader
}

//上传
func UploadHandler(ctx *modules.Context, form UploadForm) {
	if form.File == nil {
		ctx.Ajax(-1, "上传参数错误！", nil)
		return
	}

	file_name := form.File.Filename
	file_path := updateModel.GetPkgPath(file_name)

	if form.OverWrite == false && updateModel.IsPkgExist(file_name) == true {
		ctx.Ajax(-2, "文件已存在，是否覆盖上传？", nil)
		return
	}

	if updateModel.CheckPkgName(file_name) == false {
		ctx.Ajax(-3, "文件名称非法！", nil)
		return
	}

	//保存包裹
	utils.SaveFileTo(form.File, file_path)

	//解压

}
