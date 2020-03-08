package mylog

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

// 在终端写日志的相关内容
type ConsoleLogger struct {
	Level LogLevel
}

//NewConsoleLog 构造函数
func NewConsoleLog(levelStr string) *ConsoleLogger {
	lv, err := parseLoggerLevel(levelStr)
	if err != nil {
		panic(err)
	}
	return &ConsoleLogger{
		Level: lv,
	}
}

// 比较日志等级
func (c *ConsoleLogger) enable(level LogLevel) bool {
	return level >= c.Level

}

//处理相同数据进行合并到此函数，共同调用
func (c *ConsoleLogger) consolelog(lv LogLevel, format string, a ...interface{}) {
	if c.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now().Format("2006-01-02 15:04:05")
		levelStr := getLoggerStrbyLevel(lv)
		funcName, fileName, lineNo := getInfo(3)
		fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", now, levelStr, fileName, funcName, lineNo, msg)
	}

}

//获取执行时函数名，文件名及行号
func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("getinfo failed")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	funcName = strings.Split(funcName, ".")[1]
	fileName = path.Base(file)
	return

}

//日志等级DEBUG方法
func (c *ConsoleLogger) Debug(format string, a ...interface{}) {
	c.consolelog(DEBUG, format, a...)
}

//日志等级WARNING方法
func (c *ConsoleLogger) Warning(format string, a ...interface{}) {
	c.consolelog(WARNING, format, a...)
}

//日志等级INFO方法
func (c *ConsoleLogger) Info(format string, a ...interface{}) {
	c.consolelog(INFO, format, a...)
}

//日志等级FATAL方法
func (c *ConsoleLogger) Fatal(format string, a ...interface{}) {
	c.consolelog(FATAL, format, a...)
}

//日志等级ERROR方法
func (c *ConsoleLogger) Error(format string, a ...interface{}) {
	c.consolelog(ERROR, format, a...)
}

//日志等级TRACE方法
func (c *ConsoleLogger) Trace(format string, a ...interface{}) {
	c.consolelog(TRACE, format, a...)
}
