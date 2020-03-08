package mylog

import (
	"errors"
	"strings"
)

//定义logger相关常量及函数
type LogLevel uint16

//定义mlogger接口
type Mlogger interface {
	Debug(format string, a ...interface{})
	Trace(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warning(format string, a ...interface{})
	Error(format string, a ...interface{})
	Fatal(format string, a ...interface{})
}

const (
	UNKOWN LogLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

func parseLoggerLevel(s string) (LogLevel, error) {
	s = strings.ToUpper(s)
	switch s {
	case "DEBUG":
		return DEBUG, nil
	case "TRACE":
		return TRACE, nil
	case "INFO":
		return INFO, nil
	case "WARNING":
		return WARNING, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	default:
		return UNKOWN, errors.New("未匹配到指定的类型")
	}
}

func getLoggerStrbyLevel(lv LogLevel) string {
	switch lv {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "DEBUG"
}

// logger接口初始化选择 c-->窗口打印 f-->文件写入
func ChoiceLoggerMode(s, l string, check bool) Mlogger {
	switch s {
	case "console", "c":
		return NewConsoleLog(l)

	case "file", "f":
		return NewFileLog(l, "./", "xx.log", 100*1024*1024, check)
	}
	return nil
}
