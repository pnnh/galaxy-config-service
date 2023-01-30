package main

import (
	"path"
	"runtime"
	"time"

	"github.com/pnnh/multiverse-server/config"
	"github.com/sirupsen/logrus"
)

func init() {
	if config.Debug() {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			TimestampFormat: time.RFC3339,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				//处理文件名
				fileName := path.Base(frame.File)
				return "", fileName
			},
		})
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.SetReportCaller(false)
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:     false,
			TimestampFormat: time.RFC3339,
		})
	}
}

func main() {
	app := NewApplication()

	if err := app.Init(); err != nil {
		logrus.Fatalln("应用程序初始化出错: %w", err)
	}
	if err := app.Start(); err != nil {
		logrus.Fatalln("应用程序执行出错: %w", err)
	}
}
