package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"sgridnext.com/src/constant"
)

func CreateLogger(logName string) *logrus.Logger {
	cwd, _ := os.Getwd()
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	})
	SGRID_LOG_DIR := filepath.Join(cwd, constant.TARGET_LOG_DIR, constant.MAIN_SERVER_NAME)
	if SGRID_LOG_DIR == "" {
		return logger
	}
	logPath := filepath.Join(SGRID_LOG_DIR, fmt.Sprintf("%s.log", logName))
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
		&logrus.JSONFormatter{
			TimestampFormat: time.DateTime,
		},
	))
	return logger
}
