package modules

import "fmt"

func (this *Context)Dump(obj interface{}) {
	this.Write([]byte(fmt.Sprintln(obj)))
}
