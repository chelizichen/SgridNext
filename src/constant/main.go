package constant

const (
	SGRID_LOG_DIR     = "SGRID_LOG_DIR"
	SGRID_CONF_DIR    = "SGRID_CONF_DIR"
	SGRID_PACKAGE_DIR = "SGRID_PACKAGE_DIR"
	SGRID_SERVANT_DIR = "SGRID_SERVANT_DIR"

	TARGET_LOG_DIR     = "server/SgridPatchServer/log"
	TARGET_CONF_DIR    = "server/SgridPatchServer/conf"
	TARGET_PACKAGE_DIR = "server/SgridPatchServer/package"
	TARGET_SERVANT_DIR = "server/SgridPatchServer/servant"
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
	ACITVATE_DEPLOY  = 1
	ACTIVATE_RESTART = 2
	ACTIVATE_STOP    = 3
)

const (
	FILE_TYPE_CONFIG  = 1
	FILE_TYPE_PACKAGE = 2
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
