package modules

import (
	"github.com/davecgh/go-spew/spew"
)

func init() {
	//spew.Config.MaxDepth = 3
	spew.Config.Indent = "    "
	spew.Config.DisableMethods = true
	spew.Config.DisablePointerMethods = true
}

//打印出对象的结构
func (this *Context)Dump(obj ...interface{}) {
	obj_str := spew.Sdump(obj)

	this.Write([]byte(obj_str))
}
