package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/choirulanwar/textify/backend/pkg/helper"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logger      *logrus.Logger
	once        sync.Once
	runtimePath string
)

func NewLogger(logPath, module string, debug bool) *logrus.Logger {
	once.Do(func() {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic("创建日志文件失败: " + err.Error())
		}
		runtimeDir := helper.GetRuntimeUserHomeDir()
		path, _ := os.Executable()
		_, pathAppName := filepath.Split(path)

		runtimePath = runtimeDir + "/" + pathAppName + "/" + logPath + "/"
		latestLogFile := runtimePath + module + ".log"
		logger = logrus.New()

		logger.Out = src

		if debug == true {
			logger.SetLevel(logrus.DebugLevel)
		}

		logWriter, err := rotatelogs.New(
			runtimePath+module+"-%Y-%m-%d.log",
			rotatelogs.WithLinkName(latestLogFile),
			rotatelogs.WithMaxAge(30*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			panic("创建日志文件失败: " + err.Error())
		}
		logger.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.DebugLevel: logWriter,
				logrus.InfoLevel:  logWriter,
				logrus.FatalLevel: logWriter,
				logrus.PanicLevel: logWriter,
				logrus.WarnLevel:  logWriter,
				logrus.ErrorLevel: logWriter,
			},
			&logrus.JSONFormatter{
				TimestampFormat: time.DateTime,
			},
		))
	})
	return logger
}

func GetRuntimePath() string {
	return runtimePath
}

func PrintInfo(args ...interface{}) {
	fmt.Println(fmt.Sprintf("\u001B[34m[%s] %v\u001B[0m", time.Now().Format(time.DateTime), args))
}

func PrintError(args ...interface{}) {
	fmt.Println(fmt.Sprintf("\x1b[31m[%s] %v\x1b[0m", time.Now().Format(time.DateTime), args))
}
