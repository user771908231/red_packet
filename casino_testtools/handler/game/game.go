package game

import (
	"bufio"
	"casino_common/utils/numUtils"
	"casino_common/utils/testUtils"
	"casino_testtools/modules"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func GameTest(ctx *modules.Context) {
	ctx.HTML(200, "game/game")
}

func GameEdit(ctx *modules.Context) {
	gameId := ctx.Query("gameid")
	fileName := "./" + gameId + "/xipai.json"
	fmt.Printf("开始编辑文件:%v\n", fileName)
	err := os.Remove(fileName)
	if err != nil {
		fmt.Printf("删除文件的时候，失败:%v\n", fileName)
	}

	fmt.Printf("开始打开文件:%v\n", fileName)

	outputFile, outputError := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666) //0666是标准的权限掩码,关于打开标识看下面
	if outputError != nil {
		fmt.Printf("An error occurred with file creation:%v\n", outputError)
		return
	}
	defer outputFile.Close()

	outputString := ctx.Query("game")
	s := strings.Split(outputString, ",")

	xipai := &testUtils.XiPai{}
	for _, s2 := range s {
		xipai.Ids = append(xipai.Ids, numUtils.String2Int(s2))
	}

	if len(xipai.Ids) < 10 {
		xipai.Ids = nil
	}

	b, e := json.Marshal(xipai)
	if e != nil {
		return
	}
	outputWriter := bufio.NewWriter(outputFile)
	outputWriter.Write(b)
	outputWriter.Flush()
	ctx.Success("提交成功！", "/game?source=1", 1)
}
