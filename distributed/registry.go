package distributed

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"sgridnext.com/server/SgridNodeServer/command"
	"sgridnext.com/src/constant"
)

type SgridDistributedRegistry interface {
	// 找注册表
	FindRegistry() ([]*command.SvrNodeStat, error)
	FindRegistryByServerName(serverName string) ([]*command.SvrNodeStat, error)
}

type DefaultRegistry struct{}

func (r *DefaultRegistry) FindRegistry() ([]*command.SvrNodeStat, error) {
	// cwd, _ := os.Getwd()
	// stat_remote_path := filepath.Join(cwd, "stat-remote.json")
	sgrid_node_dir := os.Getenv(constant.SGRID_NODE_DIR)
	stat_remote_path := filepath.Join(sgrid_node_dir, "stat-remote.json")
	fmt.Println("FindRegistry >> stat_remote_path: ", stat_remote_path)
	file, err := os.ReadFile(stat_remote_path)
	if err != nil {
		return nil, err
	}
	var nodeStatMap *command.SvrNodeStatMap
	err = json.Unmarshal(file, &nodeStatMap)
	if err != nil {
		fmt.Println("FindRegistry >> json.Unmarshal error: ", err.Error())
		return nil, err
	}
	fmt.Printf("FindRegistry | Content | %s \n", string(file))
	return nodeStatMap.StatList, nil
}

func (r *DefaultRegistry) FindRegistryWithPath(p string) (*command.SvrNodeStatMap, error) { 
	stat_remote_path := p
	fmt.Println("FindRegistry >> stat_remote_path: ", stat_remote_path)
	file, err := os.ReadFile(stat_remote_path)
	if err != nil {
		return nil, err
	}
	var nodeStatMap *command.SvrNodeStatMap
	err = json.Unmarshal(file, &nodeStatMap)
	if err != nil {
		fmt.Println("FindRegistry >> json.Unmarshal error: ", err.Error())
		return nil, err
	}
	fmt.Printf("FindRegistry | Content | %s \n", string(file))
	return nodeStatMap, nil
}

func (r *DefaultRegistry) FindRegistryByServerName(serverName string) ([]*command.SvrNodeStat, error) {
	statList, err := r.FindRegistry()
	if err != nil {
		fmt.Println("FindRegistryByServerName >> FindRegistry error: ", err.Error())
		return nil, err
	}
	var result []*command.SvrNodeStat
	for _, stat := range statList {
		if stat.ServerName == serverName {
			result = append(result, stat)
		}
	}
	fmt.Println("FindRegistryByServerName >> result: ", result)
	return result, nil
}
