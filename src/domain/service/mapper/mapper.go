package mapper

import (
	"encoding/json"

	"gorm.io/gorm"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
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
	T_Mapper.db.AutoMigrate(&entity.NodeStat{})
	T_Mapper.db.AutoMigrate(&entity.ServerNodeLimit{})
	T_Mapper.db.AutoMigrate(&entity.Document{})
	T_Mapper.db.AutoMigrate(&entity.DocumentServerRelation{})
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

func (t *T_PatchServer_Mapper) GetNodeInfo(id int) (entity.Node, error) {
	var node entity.Node
	res := t.db.Debug().Where("id =?", id).First(&node)
	return node, res.Error
}

func (t *T_PatchServer_Mapper) UpdateMachineNodeStatus(id int, status int) error {
	err := t.db.Debug().
		Model(&entity.Node{}).
		Where("id = ?", id).
		Update("node_status", status).
		Error
	return err
}

func (t *T_PatchServer_Mapper) UpdateMachineNodeAlias(id int, alias string) error {
	err := t.db.Debug().
		Model(&entity.Node{}).
		Where("id = ?", id).
		Update("alias", alias).
		Error
	return err
}

func (t *T_PatchServer_Mapper) UpdateNodePatch(ids []int, patchId int) error {
	logger.Mapper.Info("更新服务节点：", ids, patchId)
	if len(ids) == 0 {
		return nil
	}
	err := t.db.Debug().
		Model(&entity.ServerNode{}).
		Where("id in ?", ids).
		Update("patch_id", patchId).
		Error
	return err
}

func (t *T_PatchServer_Mapper) UpdateNodeStatus(id int, status int) error {
	err := t.db.Debug().
		Model(&entity.ServerNode{}).
		Where("id = ?", id).
		Update("server_node_status", status).
		Error
	return err
}

func (t *T_PatchServer_Mapper) UpdateServerNode(node entity.ServerNode) error {
	err := t.db.Debug().
		Model(&entity.ServerNode{}).
		Where("id = ?", node.ID).
		Update("additional_args", node.AdditionalArgs).
		Update("server_run_type", node.ServerRunType).
		Update("view_page", node.ViewPage).
		Error
	return err
}

func (t *T_PatchServer_Mapper) DeleteServerNode(ids []int) error {
	err := t.db.Debug().
		Model(&entity.ServerNode{}).
		Where("id in (?)", ids).
		Update("server_node_status", constant.COMM_STATUS_DELETE).
		Error
	return err
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
	Id               int     `json:"id"`
	Host             string  `json:"host"`
	Port             int     `json:"port"`
	PatchId          int     `json:"patch_id"`
	NodeCreateTime   string  `json:"node_create_time"`
	ServerNodeStatus int     `json:"server_node_status"`
	CpuLimit         float64 `json:"cpu_limit"`
	MemoryLimit      int     `json:"memory_limit"`
	ServerRunType    int     `json:"server_run_type"`
	AdditionalArgs   string  `json:"additional_args"`
	ServerId         int     `json:"server_id"`
	ViewPage         string  `json:"view_page"`
	Alias            string  `json:"alias"`
}

func (t *T_PatchServer_Mapper) GetServerNodes(serverId int, nodeId int) ([]ServerNodesVo, error) {
	var servers []ServerNodesVo
	var params []interface{}
	where := " where 1 = 1"

	if serverId > 0 {
		where += " and server_nodes.server_id = ? "
		params = append(params, serverId)
	}

	if nodeId > 0 {
		where += " and server_nodes.node_id = ? "
		params = append(params, nodeId)
	}
	where += " and server_node_status in (1,2) "

	query := `
	SELECT 
		server_nodes.id as id,
		server_nodes.port as port,
		server_nodes.patch_id as patch_id,
		server_nodes.create_time as node_create_time,
		server_nodes.server_node_status as server_node_status,
		server_nodes.server_run_type as server_run_type,
		server_nodes.additional_args as additional_args,
		server_nodes.server_id as server_id,
		server_nodes.view_page as view_page,
		nodes.host as host,
		server_node_limits.cpu_limit as cpu_limit,
		server_node_limits.memory_limit as memory_limit,
		nodes.alias as alias
	FROM server_nodes
	LEFT JOIN 
		nodes ON server_nodes.node_id = nodes.id
	LEFT JOIN
		server_node_limits ON server_nodes.id = server_node_limits.server_node_id
	`
	query += where
	res := t.db.Debug().Raw(query, params...).Scan(&servers)
	return servers, res.Error
}

func (t *T_PatchServer_Mapper) GetNodeList() ([]entity.Node, error) {
	var nodes []entity.Node
	res := t.db.Debug().Find(&nodes)
	return nodes, res.Error
}

func (t *T_PatchServer_Mapper) UpdateNodeUpdateTime(id int) error {
	err := t.db.Debug().
		Model(&entity.Node{}).
		Where("id = ?", id).
		Update("update_time", constant.GetCurrentTime()).
		Error
	return err
}
func (t *T_PatchServer_Mapper) GetServerPackageList(id int, offset int, size int) ([]entity.ServerPackage, int64, error) {
	var packages []entity.ServerPackage
	var total int64
	// 根据id 倒叙
	res := t.db.Debug().Model(&entity.ServerPackage{}).
		Where("server_id = ?", id).
		Count(&total).
		Order("id desc").
		Offset(offset).
		Limit(size).
		Find(&packages)
	return packages, total, res.Error
}

func (t *T_PatchServer_Mapper) GetServerPackageInfo(id int) (entity.ServerPackage, error) {
	var serverPackage entity.ServerPackage
	res := t.db.Debug().Where("id =?", id).First(&serverPackage)
	return serverPackage, res.Error
}

// NODE STAT DAT
func (t *T_PatchServer_Mapper) SaveNodeStat(req *entity.NodeStat) (int, error) {
	// 从配置文件中获取nodeId
	req.NodeId = config.Conf.GetLocalNodeId()
	req.CreateTime = constant.GetCurrentTime()
	res := t.db.Debug().Create(req)
	if res.Error != nil {
		return 0, res.Error
	}
	return req.Id, nil
}

type PageGetNodeStatListRsp struct {
	Total int64             `json:"total"`
	List  []entity.NodeStat `json:"list"`
}

func (t *T_PatchServer_Mapper) GetNodeStatList(req *entity.NodeStat, offset int, size int) (PageGetNodeStatListRsp, error) {
	var rsp PageGetNodeStatListRsp = PageGetNodeStatListRsp{
		Total: 0,
		List:  []entity.NodeStat{},
	}
	var queryParams []interface{}
	where := "  1 = 1 "
	if req.NodeId > 0 {
		where += " and node_id = ?"
		queryParams = append(queryParams, req.NodeId)
	}
	if req.TYPE > 0 {
		where += " and type = ? "
	}
	if req.ServerId > 0 {
		where += " and server_id = ? "
		queryParams = append(queryParams, req.ServerId)
	}
	if req.ServerNodeId > 0 {
		where += " and server_node_id = ? "
		queryParams = append(queryParams, req.ServerNodeId)
	}
	res := t.db.Debug().
		Model(&entity.NodeStat{}).
		Where(where, queryParams...).
		Count(&rsp.Total).
		Limit(size).
		Offset(offset).
		Order("id desc").
		Find(&rsp.List)
	return rsp, res.Error
}

// cgroup limit
func (t *T_PatchServer_Mapper) GetServerNodeLimitList(serverNodeIds []int) ([]entity.ServerNodeLimit, error) {
	var limits []entity.ServerNodeLimit
	res := t.db.Debug().
		Model(&entity.ServerNodeLimit{}).
		Where("server_node_id in ?", serverNodeIds).
		Find(&limits)
	return limits, res.Error
}

func (t *T_PatchServer_Mapper) UpsertServerNodeLimit(req *entity.ServerNodeLimit) error {
	res := t.db.Debug().
		Where("server_node_id = ?", req.ServerNodeId).
		Assign(entity.ServerNodeLimit{
			CpuLimit:    req.CpuLimit,
			MemoryLimit: req.MemoryLimit,
			UpdateTime:  constant.GetCurrentTime(),
		}).
		FirstOrCreate(req)
	return res.Error
}

func (t *T_PatchServer_Mapper) GetHost(nodeId int) (string, error) {
	var Node entity.Node
	res := t.db.Debug().
		Model(&entity.Node{}).
		Where("id =?", nodeId).
		First(&Node)
	return Node.Host, res.Error
}

func (t *T_PatchServer_Mapper) GetNodeIdByHost(host string) (int, error) {
	var Node entity.Node
	res := t.db.Debug().
		Model(&entity.Node{}).
		Where("host =?", host).
		First(&Node)
	return Node.ID, res.Error
}

func (t *T_PatchServer_Mapper) UpdateServer(req *entity.Server) error {
	res := t.db.Debug().
		Model(&entity.Server{}).
		Where("id = ?", req.ID).
		Update("docker_name", req.DockerName).
		Update("log_path", req.LogPath).
		Update("exec_file_path", req.ExecFilePath).
		Update("description", req.Description).
		Update("server_type", req.ServerType).
		Update("config_path", req.ConfigPath)
	return res.Error
}

// ========== 文档管理相关方法 ==========

// CreateDocument 创建文档
func (t *T_PatchServer_Mapper) CreateDocument(req *entity.Document) (int, error) {
	res := t.db.Debug().Create(req)
	if res.Error != nil {
		return 0, res.Error
	}
	return req.ID, nil
}

// GetDocumentList 获取文档列表
func (t *T_PatchServer_Mapper) GetDocumentList() ([]entity.Document, error) {
	var documents []entity.Document
	res := t.db.Debug().Find(&documents)
	return documents, res.Error
}

// GetDocumentById 根据ID获取文档
func (t *T_PatchServer_Mapper) GetDocumentById(id int) (entity.Document, error) {
	var document entity.Document
	res := t.db.Debug().Where("id = ?", id).First(&document)
	return document, res.Error
}

// UpdateDocument 更新文档
func (t *T_PatchServer_Mapper) UpdateDocument(req *entity.Document) error {
	res := t.db.Debug().
		Model(&entity.Document{}).
		Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"title":       req.Title,
			"content":     req.Content,
			"update_time": req.UpdateTime,
			"description": req.Description,
		})
	return res.Error
}

// DeleteDocument 删除文档
func (t *T_PatchServer_Mapper) DeleteDocument(id int) error {
	// 先删除关联关系
	t.db.Debug().Where("document_id = ?", id).Delete(&entity.DocumentServerRelation{})
	// 再删除文档
	res := t.db.Debug().Where("id = ?", id).Delete(&entity.Document{})
	return res.Error
}

// CreateDocumentServerRelation 创建文档服务关联
func (t *T_PatchServer_Mapper) CreateDocumentServerRelation(documentId int, serverId int) error {
	relation := &entity.DocumentServerRelation{
		DocumentId: documentId,
		ServerId:   serverId,
		CreateTime: constant.GetCurrentTime(),
	}
	res := t.db.Debug().Create(relation)
	return res.Error
}

// DeleteDocumentServerRelation 删除文档服务关联
func (t *T_PatchServer_Mapper) DeleteDocumentServerRelation(documentId int, serverId int) error {
	res := t.db.Debug().
		Where("document_id = ? AND server_id = ?", documentId, serverId).
		Delete(&entity.DocumentServerRelation{})
	return res.Error
}

// GetDocumentServerRelations 获取文档关联的服务列表
func (t *T_PatchServer_Mapper) GetDocumentServerRelations(documentId int) ([]int, error) {
	var relations []entity.DocumentServerRelation
	res := t.db.Debug().
		Where("document_id = ?", documentId).
		Find(&relations)
	
	serverIds := make([]int, 0, len(relations))
	for _, rel := range relations {
		serverIds = append(serverIds, rel.ServerId)
	}
	return serverIds, res.Error
}

// GetServerDocumentRelations 获取服务关联的文档列表
func (t *T_PatchServer_Mapper) GetServerDocumentRelations(serverId int) ([]int, error) {
	var relations []entity.DocumentServerRelation
	res := t.db.Debug().
		Where("server_id = ?", serverId).
		Find(&relations)
	
	documentIds := make([]int, 0, len(relations))
	for _, rel := range relations {
		documentIds = append(documentIds, rel.DocumentId)
	}
	return documentIds, res.Error
}
