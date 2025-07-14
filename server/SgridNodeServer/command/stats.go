package command

import (
	"encoding/json"
	"os"
	"path/filepath"

	"sgridnext.com/src/constant"
	"sgridnext.com/src/logger"
)

type SvrNodeStat struct {
	NodeId     int    `json:"node_id,omitempty"`     // 服务节点ID
	ServerName string `json:"server_name,omitempty"` // 服务名称
	Pid        int    `json:"pid,omitempty"`         // Pid
	ServerHost string `json:"host,omitempty"`        // 主机地址
	ServerPort int    `json:"port,omitempty"`        // 主机端口
	MachineId  int    `json:"machine_id,omitempty"`  // 机器ID
	ServerId   int    `json:"server_id,omitempty"`   // 服务ID
	DockerName string `json:"docker_name,omitempty"` // 容器名称
}


type SvrNodeStatMap struct {
	UpdateTime string         `json:"update_time,omitempty"`
	StatList   []*SvrNodeStat `json:"stat_list,omitempty"`
}

func LoadStatList() *SvrNodeStatMap {
	cwd, _ := os.Getwd()
	stat_path := filepath.Join(cwd, "stat.json")
	jsonStr, err := os.ReadFile(stat_path)
	if err != nil {
		logger.App.Errorf("读取stat.json文件失败: %v", err)
		return nil
	}
	var statMap *SvrNodeStatMap
	err = json.Unmarshal(jsonStr, &statMap)
	if err != nil {
		logger.App.Errorf("解析stat.json文件失败: %v", err)
		return nil
	}
	return statMap
}

func InitStatFile(){
	cwd, _ := os.Getwd()
	stat_path := filepath.Join(cwd, "stat.json")
	file,err  := os.Create(stat_path)
	if err != nil {
		logger.App.Errorf("创建stat.json文件失败: %v", err)
		return
	}
	logger.App.Infof("创建stat.json文件成功: %v", stat_path)
	defer file.Close()
	jsonStr, _ := json.Marshal(&SvrNodeStatMap{
		UpdateTime: constant.GetCurrentTime(),
		StatList: []*SvrNodeStat{},
	})
	_, err = file.Write(jsonStr)
	if err != nil {
		logger.App.Errorf("创建stat.json文件失败: %v", err)
		return
	}
	logger.App.Infof("创建stat.json文件成功: %v", stat_path)
	return 
}


func InitCommands(snsList []*SvrNodeStat) {
	for _, svr := range snsList {
		cmd := NewPidCommand(
			svr.Pid,
			svr.ServerName,
			svr.NodeId,
		)
		cmd.SetHost(svr.ServerHost)
		cmd.SetPort(svr.ServerPort)
		cmd.SetLocalMachineId(svr.MachineId)
		cmd.SetServerId(svr.ServerId)
		CenterManager.AddCommand(svr.NodeId,cmd)
	}
}

func LoadSvrStat(snsList []*SvrNodeStat, nodeid int) *SvrNodeStat {
	if snsList == nil {
		logger.App.Errorf("snsList is nil")
		return nil
	}
	for _, sns := range snsList {
		if sns.NodeId == nodeid {
			return sns
		}
	}
	return nil
}

func (cm *centerManager) SyncStat(createTime string) {
	cwd, _ := os.Getwd()
	stat_path := filepath.Join(cwd, "stat.json")
	// 将这块信息同步到本地文件
	statMap := &SvrNodeStatMap{
		StatList: make([]*SvrNodeStat, 0),
		UpdateTime: createTime,
	}
	for _, cmd := range cm.GetCommands() {
		sns := &SvrNodeStat{
			NodeId:     cmd.GetNodeId(),
			ServerName: cmd.GetServerName(),
			Pid:        cmd.GetPid(),
			MachineId:  cmd.GetLocalMachineId(),
			ServerHost: cmd.GetHost(),
			ServerPort: cmd.GetPort(),
			ServerId: 	cmd.GetServerId(),
			DockerName: cmd.GetDockerName(),
		}
		statMap.StatList = append(statMap.StatList, sns)
	}
	jsonStr, _ := json.Marshal(statMap)
	logger.Alive.Infof("同步状态: SyncStat ｜%s", string(jsonStr))
	outFile, err := os.OpenFile(stat_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.App.Errorf("创建文件失败: SyncStat |%v", err)
		return
	}
	defer outFile.Close()
	if _, err := outFile.Write(jsonStr); err != nil {
		logger.App.Errorf("文件写入失败: SyncStat | %v", err)
		return
	}
}
