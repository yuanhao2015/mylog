package mylog

import (
	"fmt"
	"os"
	"path"
	"time"
)

// 往文件里写相关的日志
type FileLogger struct {
	level         LogLevel
	filePath      string
	fileName      string
	fileObj       *os.File
	errFileObj    *os.File
	maxFileSize   int64
	prevHour      int
	isSplitByDate bool
	msgChan       chan *logMsg
}

type logMsg struct {
	level     LogLevel
	msg       string
	funcName  string
	fileName  string
	timestamp string
	line      int
}

//NewFileLog 构造函数
func NewFileLog(levelStr, fp, fn string, mxfsize int64, issplitbydate bool) *FileLogger {
	lv, err := parseLoggerLevel(levelStr)
	if err != nil {
		panic(err)
	}
	fl := FileLogger{
		level:         lv,
		fileName:      fn,
		filePath:      fp,
		maxFileSize:   mxfsize,
		isSplitByDate: issplitbydate,
		msgChan:       make(chan *logMsg, 50000),
	}
	err = fl.initFile()
	if err != nil {
		panic(err)
	}

	return &fl
}

//初始化文件操作
func (f *FileLogger) initFile() error {
	//打开文件
	fullFileName := path.Join(f.filePath, f.fileName)
	errfullFileName := path.Join(f.filePath, f.fileName+".err")
	f1, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open file is faied,err %v", err)
		return err
	}
	f2, err := os.OpenFile(errfullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open file is faied,err %v", err)
		return err
	}
	//赋值给全局fileObj
	f.fileObj = f1
	f.errFileObj = f2

	//开启1个groutine异步写日志
	//for i := 1; i <= 5; i++ {
	go f.writeLogBackgroud()
	//}
	return nil

}

//根据f.isSplitByDate值来判断是进行大小或者时间分隔
//1.检查当前是几点，然后和上次保存的时间比较，如果不同则进行分隔
//2.检查文件大小是否匹配设定值，否则就进行切割
func (f *FileLogger) checkFileSplit(file *os.File) (*os.File, error) {
	fullPath := path.Join(f.filePath, file.Name())
	fileinfo, err := os.Stat(fullPath)
	if err != nil {
		fmt.Printf("get file info failed,error %v", err)
		return nil, err
	}
	if f.isSplitByDate {
		fileHour := fileinfo.ModTime().Hour()
		nowtimeHour := time.Now().Hour()
		if f.prevHour == 0 {
			f.prevHour = nowtimeHour
		}
		if fileHour > f.prevHour {
			f.prevHour = fileHour
			return SplitFile(file, fullPath, true)
		}

	} else {
		if fileinfo.Size() >= f.maxFileSize {
			return SplitFile(file, fullPath, false)
		}
	}
	return file, nil
}

//文件切割
func SplitFile(file *os.File, fullPath string, isdate bool) (*os.File, error) {
	//关闭当前日志
	file.Close()
	//备份一下 rename xx.log --> xx.log.bak202003111709
	nowStr := time.Now().Format("20060102150405000")
	if isdate {
		h, _ := time.ParseDuration("-1h")
		nowStr = time.Now().Add(h).Format("2006010215")
	}
	newLogName := fmt.Sprintf("%s.bak%s", fullPath, nowStr)
	os.Rename(fullPath, newLogName)
	//重新打开一个新的日志文件
	fileObj, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("open new file failed,error: %v", err)
		return nil, err
	}
	return fileObj, nil
}

//比较日志等级
func (f *FileLogger) enable(level LogLevel) bool {
	return level >= f.level

}

//后台写日志
func (f *FileLogger) writeLogBackgroud() {
	// msg := fmt.Sprintf(format, a...)
	// now := time.Now().Format("2006-01-02 15:04:05")

	// funcName, fileName, lineNo := getInfo(3)
	for {
		select {
		case logTmp := <-f.msgChan:
			levelStr := getLoggerStrbyLevel(logTmp.level)
			var err error
			f.fileObj, err = f.checkFileSplit(f.fileObj)
			if err != nil {
				panic(err)
			}
			logInfo := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", logTmp.timestamp, levelStr, logTmp.fileName, logTmp.funcName, logTmp.line, logTmp.msg)
			fmt.Println(logInfo)
			fmt.Fprintf(f.fileObj, logInfo)
			if logTmp.level >= ERROR {
				f.errFileObj, err = f.checkFileSplit(f.errFileObj)
				if err != nil {
					panic(err)
				}
				fmt.Fprintf(f.errFileObj, logInfo)
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}

}

//将获取的等级日志写入文件
func (f *FileLogger) consolelog(lv LogLevel, format string, a ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName, fileName, lineNo := getInfo(3)
		logData := &logMsg{
			level:     lv,
			fileName:  fileName,
			funcName:  funcName,
			line:      lineNo,
			msg:       msg,
			timestamp: now,
		}
		//防止通道满了，阻塞
		select {
		case f.msgChan <- logData:
		default:
		}

	}

}

//日志等级DEBUG方法
func (f *FileLogger) Debug(format string, a ...interface{}) {
	f.consolelog(DEBUG, format, a...)
}

//日志等级WARNING方法
func (f *FileLogger) Warning(format string, a ...interface{}) {
	f.consolelog(WARNING, format, a...)
}

//日志等级INFO方法
func (f *FileLogger) Info(format string, a ...interface{}) {
	f.consolelog(INFO, format, a...)
}

//日志等级FATAL方法
func (f *FileLogger) Fatal(format string, a ...interface{}) {
	f.consolelog(FATAL, format, a...)
}

//日志等级ERROR方法
func (f *FileLogger) Error(format string, a ...interface{}) {
	f.consolelog(ERROR, format, a...)
}

//日志等级TRACE方法
func (f *FileLogger) Trace(format string, a ...interface{}) {
	f.consolelog(TRACE, format, a...)
}
