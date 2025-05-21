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
		logger.App.Errorf("HTTP请求失败: %v", err)
		return fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.App.Errorf("接口返回错误状态码: %d", err)
		return fmt.Errorf("接口返回错误状态码: %d", resp.StatusCode)
	}

	var filePath string
	cwd, _ := os.Getwd()
	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	serverName := serverInfo.ServerName
	// serverName := "SgridTestJavaServer"
	switch req.Type {
	case constant.FILE_TYPE_PACKAGE:
		filePath = filepath.Join(cwd, constant.TARGET_PACKAGE_DIR, serverName, req.FileName)
	case constant.FILE_TYPE_CONFIG:
		filePath = filepath.Join(cwd, constant.TAGET_CONF_DIR, serverName, req.FileName)
	default:
		logger.App.Errorf("未知的文件类型: %d", req.Type)
		return fmt.Errorf("未知的文件类型: %d", req.Type)
	}
	fmt.Println("创建目录 | filePath", filePath)
	if err = os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		logger.App.Errorf("创建目录失败: %v", err)
		return fmt.Errorf("创建目录失败: %v", err)
	}
	fmt.Println("写入文件 | filePath", filePath)
	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.App.Errorf("创建文件失败: %v", err)
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		logger.App.Errorf("文件写入失败: %v", err)
		return fmt.Errorf("文件写入失败: %v", err)
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

func GetConfigList(req GetConfigListReq) error {
	client := &http.Client{}
	mainNodePath := config.Conf.Get("mainNode")
	// mainNodePath := "http://124.220.19.199:15872"
	apiPath := fmt.Sprintf("%s/api/server/getConfigList", mainNodePath)
	fmt.Printf("apiPath: %s \n", apiPath)
	resp, err := client.Post(apiPath, "application/json", bytes.NewBuffer(req.ToJSON()))
	if err != nil {
		logger.App.Errorf("HTTP请求失败: %v", err)
		return fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.App.Errorf("接口返回错误状态码: %d", err)
		return fmt.Errorf("接口返回错误状态码: %d", resp.StatusCode)
	}
	var respData GetConfigListResp
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		logger.App.Errorf("JSON解码失败: %v", err)
		return fmt.Errorf("JSON解码失败: %v", err)
	}
	if !respData.Success {
		logger.App.Errorf("接口返回错误: %s", respData.Msg)
		return fmt.Errorf("接口返回错误: %s", respData.Msg)
	}
	configList := respData.Data
	if configList == nil {
		logger.App.Errorf("接口返回空配置列表")
		return nil
	}
	cwd, _ := os.Getwd()
	serverInfo, err := mapper.T_Mapper.GetServerInfo(req.ServerId)
	serverName := serverInfo.ServerName
	// serverName := "SgridTestJavaServer"
	configDir := filepath.Join(cwd, constant.TAGET_CONF_DIR, serverName)
	if err = os.MkdirAll(configDir, 0755); err != nil {
		logger.App.Errorf("创建目录失败: %v", err)
		return fmt.Errorf("创建目录失败: %v", err)
	}
	for _, configName := range configList {
		if filepath.Ext(configName) == "_" {
			continue
		}
		fmt.Println("configName", configName)
		err = GetFile(FileReq{
			ServerId: req.ServerId,
			FileName: configName,
			Type:     constant.FILE_TYPE_CONFIG,
		})
		if err != nil {
			fmt.Println("download file error", err)
			mapper.T_Mapper.SaveNodeStat(&entity.NodeStat{
				ServerId: req.ServerId,
				TYPE:     entity.TYPE_SUCCESS,
				Content: fmt.Sprintf("node %v | serverName: %s | type %v |download file error : %s success",
					config.Conf.GetLocalNodeId(), serverName, constant.FILE_TYPE_CONFIG, configName),
				ServerName: serverName,
			})
			return err
		}
	}
	return nil
}
