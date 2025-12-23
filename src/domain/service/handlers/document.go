package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
)

// UploadDocument 上传文档
func UploadDocument(ctx *gin.Context) {
	cwd, _ := os.Getwd()
	docDir := filepath.Join(cwd, "doc")

	// 确保目录存在
	if err := os.MkdirAll(docDir, 0755); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建文档目录失败", "error": err.Error()})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文件上传失败", "error": err.Error()})
		return
	}

	// 检查文件扩展名
	fileName := file.Filename
	if !strings.HasSuffix(strings.ToLower(fileName), ".md") {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "只支持上传 .md 文件"})
		return
	}

	// 保存文件
	filePath := filepath.Join(docDir, fileName)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "保存文件失败", "error": err.Error()})
		return
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "读取文件内容失败", "error": err.Error()})
		return
	}

	// 创建文档记录
	title := strings.TrimSuffix(fileName, ".md")
	document := &entity.Document{
		ID:          0,
		Title:       title,
		FileName:    fileName,
		Content:     string(content),
		CreateTime:  constant.GetCurrentTime(),
		UpdateTime:  constant.GetCurrentTime(),
		Description: "",
	}

	docId, err := mapper.T_Mapper.CreateDocument(document)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建文档记录失败", "error": err.Error()})
		return
	}

	logger.App.Infof("上传文档成功: %s, ID: %d", fileName, docId)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "上传文档成功", "data": docId})
}

// CreateDocument 创建/编写文档
func CreateDocument(ctx *gin.Context) {
	var req struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		Description string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	if req.Title == "" {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文档标题不能为空"})
		return
	}

	cwd, _ := os.Getwd()
	docDir := filepath.Join(cwd, "doc")
	if err := os.MkdirAll(docDir, 0755); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建文档目录失败", "error": err.Error()})
		return
	}

	// 生成文件名
	fileName := req.Title
	if !strings.HasSuffix(strings.ToLower(fileName), ".md") {
		fileName += ".md"
	}
	filePath := filepath.Join(docDir, fileName)

	// 保存文件
	if err := os.WriteFile(filePath, []byte(req.Content), 0644); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "保存文件失败", "error": err.Error()})
		return
	}

	// 创建文档记录
	document := &entity.Document{
		ID:          0,
		Title:       req.Title,
		FileName:    fileName,
		Content:     req.Content,
		CreateTime:  constant.GetCurrentTime(),
		UpdateTime:  constant.GetCurrentTime(),
		Description: req.Description,
	}

	docId, err := mapper.T_Mapper.CreateDocument(document)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "创建文档记录失败", "error": err.Error()})
		return
	}

	logger.App.Infof("创建文档成功: %s, ID: %d", req.Title, docId)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "创建文档成功", "data": docId})
}

// GetDocumentList 获取文档列表
func GetDocumentList(ctx *gin.Context) {
	documents, err := mapper.T_Mapper.GetDocumentList()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取文档列表失败", "error": err.Error()})
		return
	}

	// 为每个文档获取关联的服务列表
	result := make([]map[string]interface{}, 0, len(documents))
	for _, doc := range documents {
		serverIds, _ := mapper.T_Mapper.GetDocumentServerRelations(doc.ID)
		result = append(result, map[string]interface{}{
			"id":          doc.ID,
			"title":       doc.Title,
			"fileName":    doc.FileName,
			"description": doc.Description,
			"createTime":  doc.CreateTime,
			"updateTime":  doc.UpdateTime,
			"serverIds":   serverIds,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// GetDocument 查看文档
func GetDocument(ctx *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 兼容 GET 请求
		idStr := ctx.Query("id")
		if idStr == "" {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文档ID不能为空"})
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文档ID格式错误"})
			return
		}
		req.ID = id
	}

	id := req.ID

	document, err := mapper.T_Mapper.GetDocumentById(id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文档不存在"})
		return
	}

	// 获取关联的服务列表
	serverIds, _ := mapper.T_Mapper.GetDocumentServerRelations(id)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"id":          document.ID,
			"title":       document.Title,
			"fileName":    document.FileName,
			"content":     document.Content,
			"description": document.Description,
			"createTime":  document.CreateTime,
			"updateTime":  document.UpdateTime,
			"serverIds":   serverIds,
		},
	})
}

// UpdateDocument 更新文档
func UpdateDocument(ctx *gin.Context) {
	var req struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Content     string `json:"content"`
		Description string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	// 获取原文档
	oldDoc, err := mapper.T_Mapper.GetDocumentById(req.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文档不存在"})
		return
	}

	cwd, _ := os.Getwd()
	docDir := filepath.Join(cwd, "doc")

	// 如果标题改变，需要重命名文件
	fileName := oldDoc.FileName
	if req.Title != oldDoc.Title {
		// 生成新文件名
		fileName = req.Title
		if !strings.HasSuffix(strings.ToLower(fileName), ".md") {
			fileName += ".md"
		}
		// 重命名文件
		oldPath := filepath.Join(docDir, oldDoc.FileName)
		newPath := filepath.Join(docDir, fileName)
		if err := os.Rename(oldPath, newPath); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "重命名文件失败", "error": err.Error()})
			return
		}
	}

	// 保存文件内容
	filePath := filepath.Join(docDir, fileName)
	if err := os.WriteFile(filePath, []byte(req.Content), 0644); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "保存文件失败", "error": err.Error()})
		return
	}

	// 更新数据库记录
	document := &entity.Document{
		ID:          req.ID,
		Title:       req.Title,
		FileName:    fileName,
		Content:     req.Content,
		UpdateTime:  constant.GetCurrentTime(),
		Description: req.Description,
	}

	if err := mapper.T_Mapper.UpdateDocument(document); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "更新文档记录失败", "error": err.Error()})
		return
	}

	logger.App.Infof("更新文档成功: %s, ID: %d", req.Title, req.ID)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "更新文档成功"})
}

// DownloadDocument 下载文档
func DownloadDocument(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文档ID不能为空"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文档ID格式错误"})
		return
	}

	document, err := mapper.T_Mapper.GetDocumentById(id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文档不存在"})
		return
	}

	cwd, _ := os.Getwd()
	filePath := filepath.Join(cwd, "doc", document.FileName)

	if _, err := os.Stat(filePath); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文件不存在"})
		return
	}

	ctx.File(filePath)
}

// DeleteDocument 删除文档
func DeleteDocument(ctx *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	// 获取文档信息
	document, err := mapper.T_Mapper.GetDocumentById(req.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文档不存在"})
		return
	}

	// 删除文件
	cwd, _ := os.Getwd()
	filePath := filepath.Join(cwd, "doc", document.FileName)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		logger.App.Warnf("删除文件失败: %v", err)
	}

	// 删除数据库记录（会自动删除关联关系）
	if err := mapper.T_Mapper.DeleteDocument(req.ID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "删除文档失败", "error": err.Error()})
		return
	}

	logger.App.Infof("删除文档成功: ID: %d", req.ID)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "删除文档成功"})
}

// LinkDocumentToServer 关联文档到服务
func LinkDocumentToServer(ctx *gin.Context) {
	var req struct {
		DocumentId int   `json:"documentId"`
		ServerIds  []int `json:"serverIds"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "参数错误"})
		return
	}

	// 先删除该文档的所有关联
	oldServerIds, _ := mapper.T_Mapper.GetDocumentServerRelations(req.DocumentId)
	for _, serverId := range oldServerIds {
		mapper.T_Mapper.DeleteDocumentServerRelation(req.DocumentId, serverId)
	}

	// 创建新的关联
	for _, serverId := range req.ServerIds {
		if err := mapper.T_Mapper.CreateDocumentServerRelation(req.DocumentId, serverId); err != nil {
			logger.App.Errorf("创建文档服务关联失败: %v", err)
		}
	}

	logger.App.Infof("关联文档到服务成功: DocumentId: %d, ServerIds: %v", req.DocumentId, req.ServerIds)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "msg": "关联成功"})
}

// GetDocumentServerRelations 获取文档关联的服务列表
func GetDocumentServerRelations(ctx *gin.Context) {
	idStr := ctx.Query("documentId")
	if idStr == "" {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文档ID不能为空"})
		return
	}

	documentId, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "文档ID格式错误"})
		return
	}

	serverIds, err := mapper.T_Mapper.GetDocumentServerRelations(documentId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "msg": "获取关联服务失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": serverIds})
}
