package util

import (
	"os"
	"path/filepath"

	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
)

func defaultGetLogDir(serverInfo *entity.Server) string {
	if serverInfo.LogPath != "" {
		return serverInfo.LogPath
	}
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, constant.TARGET_LOG_DIR, serverInfo.ServerName)
}

func GetLogPath(serverInfo *entity.Server, logFile string) string {
	return filepath.Join(defaultGetLogDir(serverInfo), logFile)
}

func GetLogDir(serverInfo *entity.Server) string {
	return defaultGetLogDir(serverInfo)
}
