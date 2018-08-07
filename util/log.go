package util

import (
	"github.com/astaxie/beego/logs"
	"os"
	"time"
)

var (
	ConsoleLog *logs.BeeLogger
	FileLog    *logs.BeeLogger
)

func init() {
	InitLog()
}

// 初始化日志
func InitLog() {
	ConsoleLog = logs.NewLogger(10000)
	ConsoleLog.SetLogger("console", ``)
	ConsoleLog.EnableFuncCallDepth(true)

	logFile := time.Now().Format("2006-01-02")
	FileLog = logs.NewLogger(10000)

	if _, err := os.Stat("logs"); err != nil {
		os.MkdirAll("logs", 0755)
	}

	FileLog.SetLogger("file", `{"filename":"logs/`+logFile+`.log"}`)
	FileLog.EnableFuncCallDepth(true)
}
