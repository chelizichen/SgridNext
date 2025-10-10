package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/entity"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
)

type FileReq struct {
	ServerId int    `json:"serverId"`
	FileName string `json:"fileName"`
	Type     int    `json:"type"`
}

func (r FileReq) ToJSON() []byte {
	return []byte(fmt.Sprintf(`{"serverId":%d,"fileName":"%s","type":%d}`,
		r.ServerId, r.FileName, r.Type))
}

func GetFile(req FileReq) error {
	client := &http.Client{}
	mainNodePath := config.Conf.Get("mainNode")
	// mainNodePath := "http://124.220.19.199:15872"
	apiPath := fmt.Sprintf("%s/api/server/getFile", mainNodePath)
	fmt.Printf("apiPath: %s \n", apiPath)
	resp, err := client.Post(apiPath, "application/json", bytes.NewBuffer(req.ToJSON()))
	if err != nil {
		err := fmt.Errorf("HTTP请求失败: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("接口返回错误状态码: %d", resp.StatusCode)
		return err
	}

	var filePath string
	cwd, _ := os.Getwd()
	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err != nil {
		err := fmt.Errorf("获取服务器信息失败: %v", err)
		return err
	}
	serverName := serverInfo.ServerName
	switch req.Type {
	case constant.FILE_TYPE_PACKAGE:
		filePath = filepath.Join(cwd, constant.TARGET_PACKAGE_DIR, serverName, req.FileName)
	case constant.FILE_TYPE_CONFIG:
		filePath = filepath.Join(cwd, constant.TARGET_CONF_DIR, serverName, req.FileName)
	default:
		err := fmt.Errorf("未知的文件类型: %d", req.Type)
		return err
	}
	logger.Package.Info("创建目录 | filePath", filePath)
	if err = os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		err := fmt.Errorf("创建目录失败: %v", err)
		return err
	}
	logger.Package.Info("写入文件 | filePath", filePath)
	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		err := fmt.Errorf("创建文件失败: %v", err)
		return err
	}
	defer outFile.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("读取响应体失败: %v", err)
		return err
	}
	if _, err := outFile.Write(bodyBytes); err != nil {
		err := fmt.Errorf("文件写入失败: %v", err)
		return err
	}
	// 如果配置文件路径不为空，则代表服务启动时配置文件不在默认路径下，需要在外部进行同时写入，不备份
	if serverInfo.ConfigPath != "" && req.Type == constant.FILE_TYPE_CONFIG {
		// 覆盖写入配置文件
		logger.Config.Infof("写入外部配置文件 %s", serverInfo.ConfigPath)
		filePath := filepath.Join(serverInfo.ConfigPath, req.FileName)

		// 确保目录存在
		if err = os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			err = fmt.Errorf("创建外部配置目录失败 %s :%v", filepath.Dir(filePath), err)
			return err
		}

		err = os.WriteFile(filePath, bodyBytes, 0755)
		if err != nil {
			err = fmt.Errorf("写入外部配置文件失败 %s :%v", serverInfo.ConfigPath, err)
			return err
		}
	}

	mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
		ServerId: req.ServerId,
		TYPE:     entity.TYPE_SUCCESS,
		Content: fmt.Sprintf("node %v | serverName: %s | type %v |download file: %s success",
			config.Conf.GetLocalNodeId(), serverName, req.Type, req.FileName),
		ServerName: serverName,
	})
	return nil
}

type GetConfigListReq struct {
	ServerId int `json:"serverId"`
}

type GetConfigListResp struct {
	Success bool     `json:"success"`
	Msg     string   `json:"msg"`
	Data    []string `json:"data"`
}

func (r GetConfigListReq) ToJSON() []byte {
	return []byte(fmt.Sprintf(`{"serverId":%d}`,
		r.ServerId))
}

// 服务启动前从主控拉取配置文件列表
func GetConfigList(req GetConfigListReq) error {
	client := &http.Client{}
	mainNodePath := config.Conf.Get("mainNode")
	// mainNodePath := "http://124.220.19.199:15872"
	apiPath := fmt.Sprintf("%s/api/server/getConfigList", mainNodePath)
	fmt.Printf("apiPath: %s \n", apiPath)
	resp, err := client.Post(apiPath, "application/json", bytes.NewBuffer(req.ToJSON()))
	if err != nil {
		err := fmt.Errorf("HTTP请求失败: %v", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("接口返回错误状态码: %d", resp.StatusCode)
		return err
	}
	var respData GetConfigListResp
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		err := fmt.Errorf("JSON解码失败: %v", err)
		return err
	}
	if !respData.Success {
		err := fmt.Errorf("接口返回错误: %s", respData.Msg)
		return err
	}
	configList := respData.Data
	if configList == nil {
		logger.Config.Infof("接口返回空配置列表")
		return nil
	}
	cwd, _ := os.Getwd()
	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	if err != nil {
		err := fmt.Errorf("获取服务器信息失败: %v", err)
		return err
	}
	serverName := serverInfo.ServerName
	configDir := filepath.Join(cwd, constant.TARGET_CONF_DIR, serverName)
	if err = os.MkdirAll(configDir, 0755); err != nil {
		err := fmt.Errorf("创建目录失败: %v", err)
		return err
	}
	for _, configName := range configList {
		if filepath.Ext(configName) == "_" {
			continue
		}
		fmt.Println("configName", configName)
		fileReq := FileReq{
			ServerId: req.ServerId,
			FileName: configName,
			Type:     constant.FILE_TYPE_CONFIG,
		}
		err = GetFile(fileReq)
		if err != nil {
			err := fmt.Errorf("下载文件失败 ｜ req %v｜ error: %v", fileReq, err)
			mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
				ServerId: req.ServerId,
				TYPE:     entity.TYPE_ERROR,
				Content: fmt.Sprintf("node %v | serverName: %s | type %v | download file %s error : %s",
					config.Conf.GetLocalNodeId(), serverName, constant.FILE_TYPE_CONFIG, configName, err.Error()),
				ServerName: serverName,
			})
			return err
		}
	}
	return nil
}
