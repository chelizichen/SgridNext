package constant

const (
	SGRID_LOG_DIR     = "SGRID_LOG_DIR"
	SGRID_CONF_DIR    = "SGRID_CONF_DIR"
	SGRID_PACKAGE_DIR = "SGRID_PACKAGE_DIR"
	SGRID_SERVANT_DIR = "SGRID_SERVANT_DIR"
	SGRID_DATA_DIR    = "SGRID_DATA_DIR"
	TARGET_LOG_DIR     = "server/SgridPatchServer/log"
	TARGET_CONF_DIR    = "server/SgridPatchServer/conf"
	TARGET_PACKAGE_DIR = "server/SgridPatchServer/package"
	TARGET_SERVANT_DIR = "server/SgridPatchServer/servant"
	TARGET_DATA_DIR    = "server/SgridPatchServer/data"
)

const (
	MAIN_SERVER_NAME = "SgridNext"
	PATCH_SERVER     = "SgridPatchServer"
)

const (
	SERVER_TYPE_NODE   = 1
	SERVER_TYPE_JAVA   = 2
	SERVER_TYPE_BINARY = 3
)

const (
	SGRID_TARGET_PORT = "SGRID_TARGET_PORT"
	SGRID_TARGET_HOST = "SGRID_TARGET_HOST"
)

const (
	COMM_STATUS_ONLINE  = 1
	COMM_STATUS_OFFLINE = 2
	COMM_STATUS_DELETE  = 3
)

const (
	ACTIVATE_DEPLOY  = 1
	ACTIVATE_RESTART = 2
	ACTIVATE_STOP    = 3
)

// NodeServer.proto DownloadFileRequest 的 type 字段
const (
	FILE_TYPE_CONFIG  = 1
	FILE_TYPE_PACKAGE = 2
	FILE_TYPE_LOG     = 3
)

const (
	CGROUP_TYPE_CPU    = 1
	CGROUP_TYPE_MEMORY = 2
	CGROUP_TYPE_DELETE = -1
)

const (
	NODE_PORT      = "25528"
	SGRID_NODE_DIR = "SGRID_NODE_DIR"
)

const (
	SERVER_RUN_TYPE_RESTART_ALWAYS = 12
)

const (
	DB_TYPE_MYSQL    = "mysql"
	DB_TYPE_POSTGRES = "postgres"
)


const (
	MSG_CALL_SIZE_MAX = 20 * 1024 * 1024
	MSG_RECV_SIZE_MAX = 20 * 1024 * 1024
)


// 日志类型定义
const (
	LOG_TYPE_BUSINESS = 1 // 业务日志
	LOG_TYPE_MASTER   = 2 // 主控日志
	LOG_TYPE_NODE     = 3 // 节点日志
)


type ConfObj struct{
	Host string `json:"host"`
	Db string `json:"db"`
	DbType string `json:"dbtype"`
	NodeIndex string `json:"nodeIndex"`
	MainNode string `json:"mainNode"`
}