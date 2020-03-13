package logger

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"bufio"
	"fmt"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var once sync.Once
var gLogger *logger

var (
	PanicLevel = uint8(logrus.PanicLevel)
	FatalLevel = uint8(logrus.FatalLevel)
	ErrorLevel = uint8(logrus.ErrorLevel)
	WarnLevel  = uint8(logrus.WarnLevel)
	InfoLevel  = uint8(logrus.InfoLevel)
	DebugLevel = uint8(logrus.DebugLevel)
)

// 创建日志记录器
func CreateLoggerOnce(level, filelevel uint8) {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	base := filepath.Base(os.Args[0])
	name := strings.TrimSuffix(base, filepath.Ext(base))
	path := "logs/" + name + "/" + tm.Format("20060102150405") + "/"

	once.Do(func() {
		gLogger = &logger{
			file:      newLoggerOfLFShook(1048576000, 100, 15, path),
			console:   newLoggerOfConsole(),
			fileLevel: logrus.Level(filelevel),
		}
		gLogger.console.SetLevel(logrus.Level(level))
		gLogger.file.SetLevel(logrus.Level(filelevel))
	})
}

// 格式化输出Debug日志
func Debugf(format string, params ...interface{}) {
	if gLogger != nil {
		for _, log := range gLogger.available(logrus.DebugLevel) {
			log.Debugf(format, params...)
		}
	}
}

// 格式化输出Info日志
func Infof(format string, params ...interface{}) {
	if gLogger != nil {
		for _, log := range gLogger.available(logrus.InfoLevel) {
			log.Infof(format, params...)
		}
	}
}

// 格式化输出Warn日志
func Warnf(format string, params ...interface{}) {
	if gLogger != nil {
		for _, log := range gLogger.available(logrus.WarnLevel) {
			log.Warnf(format, params...)
		}
	}
}

// 格式化输出Error日志
func Errorf(format string, params ...interface{}) {
	if gLogger != nil {
		for _, log := range gLogger.available(logrus.ErrorLevel) {
			log.Errorf(format, params...)
		}
	}
}

// 格式化输出Fatal日志
func Fatalf(format string, params ...interface{}) {
	if gLogger != nil {
		for _, log := range gLogger.available(logrus.FatalLevel) {
			log.Fatalf(format, params...)
		}
	}
}

// 格式化输出Panic日志
func Panicf(format string, params ...interface{}) {
	if gLogger != nil {
		for _, log := range gLogger.available(logrus.PanicLevel) {
			log.Panicf(format, params...)
		}
	}
}

// 日志选项
type logger struct {
	file      *logrus.Logger // 文件日志
	console   *logrus.Logger // 控制台日志
	fileLevel logrus.Level   // 文件日志级别
}

// 获取可用的日志记录器
func (lg *logger) available(level logrus.Level) []*logrus.Logger {
	logs := make([]*logrus.Logger, 0)
	if level <= lg.console.Level {
		logs = append(logs, lg.console)
	}
	if level <= lg.fileLevel {
		logs = append(logs, lg.file)
	}
	return logs
}

// 创建终端记录器
func newLoggerOfConsole() *logrus.Logger {
	lg := logrus.New()
	for _, level := range logrus.AllLevels {
		lg.Level |= level
	}
	lg.Formatter = &logrus.JSONFormatter{}
	return lg
}

// 创建文件记录器
func newLoggerOfLFShook(maxsize int, maxbackup int, maxage int, path string) *logrus.Logger {
	lg := logrus.New()

	// 文件记录器不输出到控制台
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	lg.SetOutput(writer)

	writerMap := lfshook.WriterMap{}
	for _, level := range logrus.AllLevels {
		lg.Level |= level
		writer := &lumberjack.Logger{
			Filename: path + level.String() + ".log",
			//MaxSize:    maxsize,
			MaxBackups: maxbackup,
			MaxAge:     maxage,
		}
		if level == logrus.DebugLevel { // debug日志只保存3天
			writer.MaxAge = 3
		}
		writerMap[level] = writer
	}
	lg.Formatter = &logrus.JSONFormatter{}
	lg.Hooks.Add(lfshook.NewHook(writerMap, nil))
	return lg
}
