package constant

import "time"

const (
	SGRID_LOG_DIR     = "SGRID_LOG_DIR"
	SGRID_CONF_DIR    = "SGRID_CONF_DIR"
	SGRID_PACKAGE_DIR = "SGRID_PACKAGE_DIR"
	SGRID_SERVANT_DIR = "SGRID_SERVANT_DIR"

	TARGET_LOG_DIR     = "server/SgridPatchServer/log"
	TAGET_CONF_DIR     = "server/SgridPatchServer/conf"
	TARGET_PACKAGE_DIR = "server/SgridPatchServer/package"
	TARGET_SERVANT_DIR = "server/SgridPatchServer/servant"
)

const (
	MAIN_SERVER_NAME = "SgridNext"
	PATCH_SERVER     = "SgridPatchServer"
)

const (
	SERVER_TYPE_NODE = 1
	SERVER_TYPE_JAVA = 2
	SERVER_TYPE_BINARY = 3
)

const (
	SGRID_TARGET_PORT = "SGRID_TARGET_PORT"
	SGRID_TARGET_HOST = "SGRID_TARGET_HOST"
)


const (
	COMM_STATUS_ONLINE = 1
	COMM_STATUS_OFFLINE = 2
	COMM_STATUS_DELETE = 3
)


// 拿到当前时间 YYYY-MM-DD HH:MM:SS
func GetCurrentTime() string {
    return time.Now().Format("2006-01-02 15:04:05")
}