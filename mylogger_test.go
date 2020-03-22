package mylog

import (
	"testing"
	"time"
)

func TestLoger(t *testing.T) {
	log := ChoiceLoggerMode("c", "debug", false)
	for {
		t.Log("------控制台输出---------")
		log.Debug("这是一条debug日志")
		log.Trace("这是一条trace日志")
		log.Info("这是一条info日志")
		log.Warning("这是一条warning日志")
		log.Error("这是一条error日志带参数，字符串参数：%s 整形参数：%d", "mmmmm", 11111)
		log.Fatal("这是一条fatal日志")
		time.Sleep(1 * time.Second)
		t.Log("------写文件输出---------")
		log = ChoiceLoggerMode("f", "debug", false)
		log.Debug("这是一条debug日志")
		log.Trace("这是一条trace日志")
		log.Info("这是一条info日志")
		log.Warning("这是一条warning日志")
		log.Error("这是一条error日志带参数，字符串参数：%s 整形参数：%d", "mmmmm", 11111)
		log.Fatal("这是一条fatal日志")
		time.Sleep(1 * time.Second)
	}
}
