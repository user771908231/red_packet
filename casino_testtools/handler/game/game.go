package game

import (
	"bufio"
	"casino_testtools/modules"
	"fmt"
	"os"
)

func GameTest(ctx *modules.Context) {
	ctx.HTML(200, "game/game")
}

func GameEdit(ctx *modules.Context) {
	outputFile, outputError := os.OpenFile("/usr/local/gametest/gameid/test.json",
		os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation\n")
		return
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	outputString := ctx.Query("game")
	outputWriter.WriteString(outputString)
	outputWriter.Flush()
	ctx.Success("提交成功！", "/game", 1)
}
