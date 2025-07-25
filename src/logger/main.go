package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func CreateLogger(logName string) *logrus.Logger {
	cwd, _ := os.Getwd()
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.DateTime,
	})

	logDir := cwd
	envDir := os.Getenv("SGRID_LOG_DIR")
	if envDir != "" {
		logDir = envDir
	}else {
		// return logger
	}

	logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", logName))
	fmt.Println(logPath)
	// 配置日志轮转
	writer, _ := rotatelogs.New(
		logPath+".%Y%m%d",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(time.Duration(14*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)

	// 设置日志级别映射
	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}
	// 设置 Hook
	logger.AddHook(lfshook.NewHook(
		writerMap,
		&logrus.TextFormatter{
			TimestampFormat: time.DateTime,
		},
	))
	return logger
}
