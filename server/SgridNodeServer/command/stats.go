package command

import (
	"encoding/json"
	"os"
	"path/filepath"

	"sgridnext.com/src/logger"
)

func LoadStatList() []*SvrNodeStat {
	cwd, _ := os.Getwd()
	stat_path := filepath.Join(cwd, "stat.json")
	jsonStr, err := os.ReadFile(stat_path)
	if err != nil {
		logger.App.Errorf("读取stat.json文件失败: %v", err)
		return nil
	}
	var snsList []*SvrNodeStat
	err = json.Unmarshal(jsonStr, &snsList)
	if err != nil {
		logger.App.Errorf("解析stat.json文件失败: %v", err)
		return nil
	}
	return snsList
}

func InitCommands(snsList []*SvrNodeStat) {
	for _, svr := range snsList {
		CenterManager.AddCommand(svr.NodeId,
			NewPidCommand(
				svr.Pid,
				svr.ServerName,
				svr.NodeId,
			),
		)
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

type SvrNodeStat struct {
	NodeId     int    `json:"node_id,omitempty"`
	ServerName string `json:"server_name,omitempty"`
	Pid        int    `json:"pid,omitempty"`
}

func (cm *centerManager) SyncStat() {
	cwd, _ := os.Getwd()
	stat_path := filepath.Join(cwd, "stat.json")
	// 将这块信息同步到本地文件
	var snsList []*SvrNodeStat
	for _, cmd := range cm.GetCommands() {
		sns := &SvrNodeStat{
			NodeId:     cmd.GetNodeId(),
			ServerName: cmd.GetServerName(),
			Pid:        cmd.GetPid(),
		}
		snsList = append(snsList, sns)
	}
	if snsList == nil {
		snsList = make([]*SvrNodeStat, 0)
	}
	jsonStr, _ := json.Marshal(snsList)
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
