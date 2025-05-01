package mapper

import (
	"encoding/json"

	"gorm.io/gorm"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/logger"
)

var T_Mapper *T_PatchServer_Mapper

func LoadMapper(db *gorm.DB) {
	T_Mapper = &T_PatchServer_Mapper{
		db: *db,
	}

	// 自动迁移模式
	T_Mapper.db.AutoMigrate(&entity.Server{})
	T_Mapper.db.AutoMigrate(&entity.Node{})
	T_Mapper.db.AutoMigrate(&entity.ServerGroup{})
	T_Mapper.db.AutoMigrate(&entity.ServerPackage{})
	T_Mapper.db.AutoMigrate(&entity.ServerNode{})
}

type T_PatchServer_Mapper struct {
	db gorm.DB
}

func (t *T_PatchServer_Mapper) CreateServer(req *entity.Server) (int, error) {
	jsonStr, _ := json.Marshal(req)
	logger.Mapper.Info("创建服务：", string(jsonStr))
	res := t.db.Debug().Create(req)
	if res.Error != nil {
		return 0, res.Error
	}
	return req.ID, nil
}

func (t *T_PatchServer_Mapper) CreateNode(req *entity.Node) (int, error) {
	res := t.db.Debug().Create(req)
	if res.Error != nil {
		return 0, res.Error
	}
	return req.ID, nil
}

// 如果有同名的组名，返回错误
func (t *T_PatchServer_Mapper) CreateGroup(req *entity.ServerGroup) (int, error) {
	logger.Mapper.Info("创建服务组：", req)
	res := t.db.Debug().Create(req)
	if res.Error != nil {
		return 0, res.Error
	}
	return req.ID, nil
}

func (t *T_PatchServer_Mapper) CreateServerNode(req []*entity.ServerNode) error {
	for _, v := range req {
		res := t.db.Debug().Create(v)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func (t *T_PatchServer_Mapper) CreateServerPackage(req *entity.ServerPackage) (int, error) {
	res := t.db.Debug().Create(req)
	if res.Error != nil {
		return 0, res.Error
	}
	return req.ID, nil
}

func (t *T_PatchServer_Mapper) GetGroupList() ([]entity.ServerGroup, error) {
	var groups []entity.ServerGroup
	res := t.db.Debug().Find(&groups)
	return groups, res.Error
}

type ServerWithGroupVO struct {
	ServerID   int    `json:"server_id"`
	ServerName string `json:"server_name"`
	GroupName  string `json:"group_name"`
	GroupId    int    `json:"group_id"`
}

func (t *T_PatchServer_Mapper) GetServerListWithGroup() ([]ServerWithGroupVO, error) {
	var servers []ServerWithGroupVO
	query := `
	SELECT 
		servers.id as server_id, 
		servers.server_name as server_name, 
		server_groups.name as group_name,
		server_groups.id as group_id
	FROM servers 
	LEFT JOIN 
		server_groups ON servers.group_id = server_groups.id`
	res := t.db.Raw(query).Scan(&servers)
	return servers, res.Error
}

func (t *T_PatchServer_Mapper) GetServerInfo(id int) (entity.Server, error) {
	var server entity.Server
	res := t.db.Debug().Where("id = ?", id).First(&server)
	return server, res.Error
}

type ServerNodesVo struct {
	Id               int    `json:"id"`
	Host             string `json:"host"`
	Port             int    `json:"port"`
	PatchId          int    `json:"patch_id"`
	NodeCreateTime   string `json:"node_create_time"`
	ServerNodeStatus int    `json:"server_node_status"`
}

func (t *T_PatchServer_Mapper) GetServerNodes(serverId int,nodeId int) ([]ServerNodesVo, error) {
	var servers []ServerNodesVo
	var params []interface{}
	where := " where 1 = 1"

	if serverId > 0 {
		where += " and server_nodes.server_id = ?"
		params = append(params, serverId)
	}

	if nodeId > 0 {
		where += " and server_nodes.node_id =?"
		params = append(params, nodeId)
	}

	query := `
	SELECT 
		server_nodes.id as id,
		server_nodes.port as port,
		server_nodes.patch_id as patch_id,
		server_nodes.create_time as node_create_time,
		server_nodes.server_node_status as server_node_status,
		nodes.host as host
	FROM server_nodes
	LEFT JOIN 
		nodes ON server_nodes.node_id = nodes.id
	`
	query += where
	res := t.db.Raw(query, params...).Scan(&servers)
	return servers, res.Error
}

func (t *T_PatchServer_Mapper) GetNodeList() ([]entity.Node, error) {
	var nodes []entity.Node
	res := t.db.Debug().Find(&nodes)
	return nodes, res.Error
}


func (t *T_PatchServer_Mapper) GetServerPackageList(id int) ([]entity.ServerPackage, error) {
	var packages []entity.ServerPackage
	// 根据id 倒叙
	res := t.db.Debug().Order("id desc").Where("server_id =?", id).Find(&packages)
	return packages, res.Error
}

func (t *T_PatchServer_Mapper) GetServerPackageInfo(id int) (entity.ServerPackage, error) {
	var serverPackage entity.ServerPackage
	res := t.db.Debug().Where("id =?", id).First(&serverPackage)
	return serverPackage, res.Error
}