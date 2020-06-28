package library

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

type Slog struct{}

func Logf(msg string,format string,args interface{}) {
	logInfo(msg,format,args,"warning")
}

func LogI(msg string,format string,args interface{}) {
	logInfo(msg,format,args,"info")
}

// 日志记录到文件
func  logInfo(msg string,format string,args interface{},logtype string) {
	tmpConf ,_ := ReadWebConfig()
	logConf := *tmpConf
	logFilePath := logConf["servicelog.logfile"]
	logFileName := logConf["servicelog.logname"]
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	// 写入文件
	src, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err", err)
	}

	// 实例化
	logger := logrus.New()

	// 设置输出
	logger.Out = src

	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName + ".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat:"2006-01-02 15:04:05",
	})
	// 新增 Hook
	logger.AddHook(lfHook)
	format = format + msg

	if logtype == "warning" {
		logger.Warnf(format,args)
	}  else {
		logger.Infof(format,args)
	}
}