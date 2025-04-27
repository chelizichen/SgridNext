package db

type Server struct {
	ID           int
	ServerName   string
	ServerType   int
	Status       int
	ExecFilePath string
	CreateTime   string
	GroupId      int
}

type ServerGroup struct {
	ID          int
	Name        string
	EngLishName string
	Status      int
	CreateTime  string
}

type ServerNode struct {
	ID               int
	ServerId         int
	NodeId           int
	Port             int
	ServerNodeStatus int
	PatchId          int
	CreateTime       string
}

type Node struct {
	ID         int
	Host       string
	NodeStatus int
	Cpus       int
	Memory     int
	Os         string
	CreateTime string
}

type ServerPackage struct {
	ID         int
	ServerId   int
	Hash       string
	CreateTime string
}
