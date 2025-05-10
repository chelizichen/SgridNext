package entity

type Server struct {
	ID           int    `gorm:"column:id;primaryKey;autoIncrement"`
	ServerName   string `gorm:"column:server_name;type:varchar(255)"`
	ServerType   int    `gorm:"column:server_type"`
	Status       int    `gorm:"column:status"`
	ExecFilePath string `gorm:"column:exec_file_path;type:varchar(255)"`
	CreateTime   string `gorm:"column:create_time;type:varchar(64)"`
	GroupId      int    `gorm:"column:group_id"`
	Description  string `gorm:"column:description;type:varchar(255)"`
}

type ServerGroup struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string `gorm:"column:name;type:varchar(255)"`
	EngLishName string `gorm:"column:english_name;type:varchar(255)"`
	Status      int    `gorm:"column:status"`
	CreateTime  string `gorm:"column:create_time;type:varchar(64)"`
}

type ServerNode struct {
	ID               int    `gorm:"column:id;primaryKey;autoIncrement"`
	ServerId         int    `gorm:"column:server_id"`
	NodeId           int    `gorm:"column:node_id"`
	Port             int    `gorm:"column:port"`
	ServerNodeStatus int    `gorm:"column:server_node_status"`
	PatchId          int    `gorm:"column:patch_id"`
	CreateTime       string `gorm:"column:create_time;type:varchar(64)"`
}

type Node struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement"`
	Host       string `gorm:"column:host;type:varchar(255)"`
	NodeStatus int    `gorm:"column:node_status"`
	Cpus       int    `gorm:"column:cpus"`
	Memory     int    `gorm:"column:memory"`
	Os         string `gorm:"column:os;type:varchar(64)"`
	CreateTime string `gorm:"column:create_time;type:varchar(64)"`
}

type ServerPackage struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement"`
	ServerId   int    `gorm:"column:server_id"`
	Hash       string `gorm:"column:hash;type:varchar(255)"`
	Commit     string `gorm:"column:commit;type:varchar(255)"`
	FileName   string `gorm:"column:file_name;type:varchar(255)"`
	CreateTime string `gorm:"column:create_time;type:varchar(64)"`
}

const (
	TYPE_SUCCESS = 1
	TYPE_ERROR  = 2
	TYPE_INFO   = 3
	TYPE_WARN   = 4
	TYPE_PATCH  = 5
	TYPE_CHECK  = 6
)

type NodeStat struct {
	Id           int    `gorm:"column:id;primaryKey;autoIncrement"`
	NodeId       int    `gorm:"column:node_id"`
	ServerName   string `gorm:"column:server_name;type:varchar(255)"`
	ServerId     int    `gorm:"column:server_id"`
	ServerNodeId int    `gorm:"column:server_node_id"`
	TYPE         int    `gorm:"column:type"`
	Content      string `gorm:"column:content;type:text"`
	CreateTime   string `gorm:"column:create_time;type:varchar(64)"`
}